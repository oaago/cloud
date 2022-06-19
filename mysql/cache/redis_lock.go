package cache

import (
	"github.com/oaago/component/logx"
	"github.com/oaago/component/redis"
	"sync"
	"time"
)

type Lock struct {
	LockName          string                                                            //锁名
	LockValue         string                                                            //锁值 改为uuid生成
	LockExpireSeconds uint                                                              //锁的过期时间
	LockSleepTime     uint                                                              //未获取到锁的休眠时间
	SearchKey         string                                                            //查询键值
	Result            interface{}                                                       //查询结果
	SetData           func(db *DBType, expire int, result interface{}, redisKey string) //设置缓存方法
	FindData          func(key string, result interface{}) (interface{}, error)         //查询缓存方法
}

var rwlock sync.RWMutex

func (lock *Lock) GetLock(db *DBType, expire int) {
	//查询缓存数据存不存在
	rwlock.RLock()
	result, err := lock.FindData(lock.SearchKey, lock.Result)
	rwlock.RUnlock()
	if result != nil {
		logx.Logger.Info("缓存中存在，从缓存内读取，直接返回结果，lockKey：", lock.LockName, ",lockValue：", lock.LockValue, ";id:", GetGid)
		return
	} else if err != nil {
		logx.Logger.Error("error:" + err.Error())
		lock.GetLock(db, expire)
	} else {
		logx.Logger.Info("缓存中不存在，开始占用锁，lockKey：", lock.LockName, ",lockValue：", lock.LockValue, ";id:", GetGid)
	}

	//查询分布锁锁是否被占用
	rwlock.Lock()
	logx.Logger.Info("开始尝试占用锁; lockKey：", lock.LockName, ",lockValue：", lock.LockValue, ";id:", GetGid)
	exists := redis.Client.SetNX(lock.LockName, lock.LockValue, time.Second*time.Duration(lock.LockExpireSeconds)).Val()
	logx.Logger.Info("开始尝试占用锁结束; lockKey：", lock.LockName, ",lockValue：", lock.LockValue, ";id:", GetGid)
	rwlock.Unlock()
	if !exists {
		//锁存在
		logx.Logger.Info("锁被占用，进行重复占用锁，开始休眠; lockKey：", lock.LockName+",lockValue：", lock.LockValue, ";id:", GetGid)
		time.Sleep(time.Millisecond * time.Duration(lock.LockSleepTime))
		//重复占有锁
		logx.Logger.Info("锁被占用，进行重复占用锁，休眠结束; lockKey：", lock.LockName+",lockValue：", lock.LockValue, ";id:", GetGid)
		lock.GetLock(db, expire)
	} else {
		//锁不存在，进行加锁
		logx.Logger.Info("锁不存在，进行占用锁成功，lockKey：", lock.LockName, ",lockValue：", lock.LockValue, ";id:", GetGid)
		lock.setData(db, expire)
	}
}

func (lock *Lock) setData(db *DBType, expire int) {
	lock.SetData(db, expire, lock.Result, lock.SearchKey)
	//拿到锁后方法执行完成，解锁
	defer func(key string, value string) {
		lua := "if redis.call(\"get\",KEYS[1]) == ARGV[1] then\n    return redis.call(\"del\",KEYS[1])\nelse\n    return 0\nend"
		varKey := []string{key}
		varValue := []string{value}
		redis.Client.Eval(lua, varKey, varValue)
		logx.Logger.Info("锁不存在，占用锁后开始解锁，lockKey：", lock.LockName, ",lockValue：", lock.LockValue, ";id:", GetGid)
	}(lock.LockName, lock.LockValue)
}
