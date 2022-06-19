package cache

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/oaago/component/logx"
	"github.com/oaago/component/redis"
	"gorm.io/gorm"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

func GetGid() (gid uint64) {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		panic(err)
	}
	return n
}

type DBType struct {
	*gorm.DB
}

type DaoDataSource struct {
	DB DBType
	db *gorm.DB
	Lock
	RedisClient *redis.Cli
}

//初始化
func NewDaoDataSource(db DBType) *DaoDataSource {
	return &DaoDataSource{
		DB: db,
		db: db.DB,
	}
}

//获取原生的db对象
func (dataSource *DaoDataSource) GetNativeDb() *gorm.DB {
	return dataSource.db
}

//代理方法
func (db *DBType) Scan(expire int, searchObj interface{}, result interface{}) (interface{}, error) {
	if expire == 0 || expire < -1 {
		err := errors.New("expire expect -1 or gt; 0")
		return nil, err
	}
	val := reflect.ValueOf(searchObj)
	kd := val.Kind()
	if kd != reflect.Struct {
		err := errors.New("searchObj expect struct")
		return nil, err
	}

	types := reflect.TypeOf(searchObj)
	filedArray := make([]string, 1+types.NumField())
	//获取结构体所有字段
	for i := 0; i < types.NumField(); i++ {
		filedName := types.Field(i).Name
		filedArray = append(filedArray, filedName)
	}

	//将结构体属性名根据首字母排序
	bubbleSort(&filedArray)

	//遍历searchObj所有字段 获取字段值后进行json序列化 并且根据字段名首字母进行排序
	paramMap := make([]map[string]interface{}, 1+len(filedArray))
	for i := 0; i < len(filedArray); i++ {
		for j := 0; j < types.NumField(); j++ {
			if filedArray[i] == types.Field(j).Name {
				param := make(map[string]interface{}, 1)
				//获取属性的字段
				marshal, _ := json.Marshal(val.Field(j).Interface())
				param[types.Field(j).Name] = string(marshal)
				paramMap[i] = param
			}
		}
	}
	marshal, _ := json.Marshal(paramMap)
	data := []byte(types.String() + string(marshal))
	m5 := md5.New()
	m5.Write(data)
	md5String := hex.EncodeToString(m5.Sum(nil))
	redisKey := "dbCache:" + types.String() + ":" + md5String
	//构建锁
	newUUID, _ := uuid.NewUUID()
	lock := Lock{
		LockName:          "lock" + ":" + redisKey,
		LockValue:         newUUID.String(),
		LockExpireSeconds: 30,
		LockSleepTime:     500,
		SearchKey:         redisKey,
		Result:            result,
		SetData:           setToRedis,
		FindData:          getFromRedis,
	}
	lock.GetLock(db, expire)
	return result, nil
}

//设置缓存
func setToRedis(db *DBType, expire int, result interface{}, redisKey string) {
	db.DB.Scan(result)
	resultJson, _ := json.Marshal(result)
	storeToCache(redisKey, expire, string(resultJson))
}

func storeToCache(key string, expire int, value string) {
	val := redis.Client.Set(key, value, time.Minute*time.Duration(expire)).Val()
	logx.Logger.Info("缓存时redis的结果：" + val)
}

//结构体字段排序
func bubbleSort(arr *[]string) {
	for i := 0; i < len(*arr)-1; i++ {
		for j := 0; j < len(*arr)-1-i; j++ {
			temp := ""
			if (*arr)[j] > (*arr)[j+1] {
				temp = (*arr)[j]
				(*arr)[j] = (*arr)[j+1]
				(*arr)[j+1] = temp
			}
		}
	}
}

//从缓存读取
func getFromRedis(md5String string, result interface{}) (interface{}, error) {
	client := redis.Client
	cmdResult := client.Get(md5String)
	val := cmdResult.Val()
	if val != "" {
		logx.Logger.Info("缓存中存在，从缓存内读取，直接返回结果，key：", md5String)

		err := json.Unmarshal([]byte(val), &result)
		if err != nil {
			return result, err
		} else {
			return result, nil
		}
	}
	logx.Logger.Info("缓存中并不存在，继续查询，key：", md5String)
	if cmdResult.Err() != nil && cmdResult.Err().Error() != "redis: nil" {
		return nil, cmdResult.Err()
	} else {
		return nil, nil
	}
}
