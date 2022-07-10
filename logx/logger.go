package logx

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/oaago/cloud/config"
	"log"
	"os"
	"strings"

	"github.com/golang-module/carbon"
	"github.com/natefinch/lumberjack"
	"github.com/oaago/common/ipstring"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

type LoggerType struct {
	Path        string `json:"path,omitempty"`
	Name        string `json:"name,omitempty"`
	EnableKafka bool   `json:"enablekafka,omitempty"`
}

var Logx *zap.Logger
var LoggerOptions = &LoggerType{}

type kafkaProducerLog struct {
	topic        string
	syncProducer sarama.SyncProducer
}

func getEncoder() zapcore.Encoder {
	// 以下两种都是EncoderConfig类型 可以使用源码中封装的 也可以自定义
	// zap.NewProductionEncoderConfig()
	// zap.NewDevelopmentEncoderConfig()
	// return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	//自定义 我们可以修改里面的key和value值实现日志的定制
	encodingConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder, //时间格式更改
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return zapcore.NewJSONEncoder(encodingConfig)

}
func init() {
	LoggerOptions.Path = config.Op.GetString("logger.path")
	LoggerOptions.Name = config.Op.GetString("server.name")
	baseLogPath := LoggerOptions.Path
	time := carbon.Now().ToDateString()
	ipstr, err := ipstring.GetIp()
	if err != nil {
		ipstr = ""
	}
	if len(LoggerOptions.Name) == 0 {
		baseLogPath = "./logs/_" + "-level-" + time + ipstr
	} else {
		if LoggerOptions.Path == "" {
			LoggerOptions.Path = "./logs/"
		}
		baseLogPath = LoggerOptions.Path + "/" + LoggerOptions.Name + "-level-" + time + "-" + ipstr
	}
	baseLogPath = baseLogPath + ".log"

	var coreArr []zapcore.Core

	//日志级别
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //error级别
		return lev >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //info和debug级别,debug级别是最低的
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})
	kafkaPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.DebugLevel
	})
	//info文件writeSyncer
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   strings.Replace(baseLogPath, "level", "info", 1), //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    300,                                              //文件大小限制,单位MB
		MaxBackups: 50,                                               //最大保留日志文件数量
		MaxAge:     30,                                               //日志文件保留天数
		Compress:   true,                                             //是否压缩处理
	})
	infoFileCore := zapcore.NewCore(getEncoder(), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(infoFileWriteSyncer)), lowPriority) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	//error文件writeSyncer
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   strings.Replace(baseLogPath, "level", "error", 1), //日志文件存放目录
		MaxSize:    300,                                               //文件大小限制,单位MB
		MaxBackups: 50,                                                //最大保留日志文件数量
		MaxAge:     30,                                                //日志文件保留天数
		Compress:   true,                                              //是否压缩处理
	})
	errorFileCore := zapcore.NewCore(getEncoder(), zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), highPriority) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	nodes := config.Op.GetStringSlice("kafka.producer.nodes")
	if LoggerOptions.EnableKafka && len(nodes) > 0 {
		config := sarama.NewConfig()
		config.Producer.RequiredAcks = sarama.WaitForLocal
		config.Producer.Return.Successes = true
		kl, err := sarama.NewSyncProducer(nodes, config)
		if err != nil {
			fmt.Errorf(err.Error())
		}
		var kpl = kafkaProducerLog{
			topic:        "logs-" + LoggerOptions.Name + "-" + time + "-" + ipstr,
			syncProducer: kl,
		}
		prodEncoder := zap.NewProductionEncoderConfig()
		prodEncoder.EncodeTime = zapcore.ISO8601TimeEncoder
		kws := zapcore.AddSync(kpl)
		kafkaCore := zapcore.NewCore(zapcore.NewJSONEncoder(prodEncoder), kws, kafkaPriority)
		coreArr = append(coreArr, infoFileCore, errorFileCore, kafkaCore)
	} else {
		coreArr = append(coreArr, infoFileCore, errorFileCore)
	}
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	filed := zap.Fields(zap.String("serviceName", LoggerOptions.Name),
		zap.String("ip", ipstr),
		zap.String("ts", time))
	Logx = zap.New(zapcore.NewTee(coreArr...), development, zap.AddCaller(), filed) //zap.AddCaller()为显示文件名和行号，可省略
	Logger = Logx.Sugar()
}

func (kl kafkaProducerLog) Write(p []byte) (n int, err error) {
	_, _, err = kl.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: kl.topic,
		Value: sarama.ByteEncoder(p),
	})
	if err != nil {
		log.Println(err)
		return 0, err
	}
	//fmt.Println("send msgs:",string(p))
	return len(p), nil
}

func (kl *kafkaProducerLog) Close() {
	kl.syncProducer.Close()
}
