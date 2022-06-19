package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/oaago/component/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"strings"
	"time"
)

var DB *gorm.DB

type DBL struct {
	DB   *gorm.DB
	Name string
}
type DBSType map[string]*DBL

var DBS = DBSType{}

var Err error

func init() {
	mapKeys := config.Op.GetStringMapString("mysql")
	if len(mapKeys) == 0 {
		return
	}
	for i, _ := range mapKeys {
		mysql, _ := NewConnect(i)
		// 每次挂载最后一个 到上面
		DB = mysql
	}
}

func NotConfigConnect(url string) *gorm.DB {
	if DB != nil {
		return DB
	}
	dsn := fmt.Sprintf(url + "?charset=utf8mb4&parseTime=true&loc=Local&readTimeout=7200s&writeTimeout=300s")
	//连接MYSQL
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,       // 禁用彩色打印
		},
	)
	DB, Err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableAutomaticPing: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "",
		},
		Logger: newLogger,
	})
	sqlDB, err := DB.Debug().DB()
	sqlDB.SetMaxIdleConns(20)   //空闲连接数
	sqlDB.SetMaxOpenConns(1000) //最大连接数
	sqlDB.SetConnMaxIdleTime(time.Hour)
	sqlDB.SetConnMaxLifetime(10 * time.Hour)
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	return DB.Debug()
}

func NewConnect(mode string) (*gorm.DB, string) {
	modeUrl := config.Op.GetString("mysql." + mode)
	if mode == "" {
		mode = "default"
	}
	if DBS[mode] != nil {
		return DBS[mode].DB.Debug(), DBS[mode].Name
	}
	dsn := fmt.Sprintf(modeUrl + "?charset=utf8mb4&parseTime=true&loc=Local&readTimeout=7200s&writeTimeout=300s")
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,       // 禁用彩色打印
		},
	)
	//连接MYSQL
	dbl := &DBL{}
	dns := strings.Split(modeUrl, "/")
	if len(dns) != 2 {
		panic("连接数据库地址不正确")
	}
	dbl.Name = dns[1]
	dbl.DB, Err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableAutomaticPing: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	sqlDB, err := dbl.DB.DB()
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	sqlDB.SetMaxIdleConns(20)   //空闲连接数
	sqlDB.SetMaxOpenConns(1000) //最大连接数
	sqlDB.SetConnMaxIdleTime(time.Hour)
	sqlDB.SetConnMaxLifetime(10 * time.Hour)
	DBS[mode] = &DBL{}
	DBS[mode] = dbl
	return DBS[mode].DB.Debug(), DBS[mode].Name
}

func GetDBByName(mode string) *gorm.DB {
	if mode == "" {
		mode = "default"
	}
	return DBS[mode].DB.Debug()
}

func Close(mode string) {
	if mode == "" {
		mode = "default"
	}
	DB, _ := DBS[mode].DB.DB()
	err := DB.Close()
	if err != nil {
		return
	}
}

func CloseAll() {
	for _, db := range DBS {
		DB, _ := db.DB.DB()
		err := DB.Close()
		if err != nil {
			return
		}
	}
}
