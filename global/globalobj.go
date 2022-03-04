package global

// Package utils 提供zinx相关工具类函数
// 包括:
//		全局配置
//		配置文件加载
//
// 当前文件描述:
// @Title  globalobj.go
// @Description  相关配置文件定义及加载方式
// @Author  Aceld - Thu Mar 11 10:32:29 CST 2019

import (
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/sun-fight/zinx-websocket/ziface"
)

// Object 定义一个全局的对象
var Object *obj

var Redis *redis.Client
var Mysql *gorm.DB

type mysqlConfig struct {
	Path         string //服务器地址:端口
	Config       string // 高级配置
	Dbname       string // 数据库名
	Username     string // 数据库用户名
	Password     string // 数据库密码
	MaxIdleConns int    // 空闲中的最大连接数
	MaxOpenConns int    // 打开到数据库的最大连接数
	LogMode      string // 日志模式
}

type redisConfig struct {
	DB       int    // redis的哪个数据库
	Addr     string // 服务器地址:端口
	Password string // 密码
}

type zapConfig struct {
	Level         string // 级别
	Format        string // 输出
	Prefix        string // 日志前缀
	Director      string // 日志文件夹
	LinkName      string // 软链接名称
	ShowLine      bool   // 显示行
	EncodeLevel   string // 编码级
	StacktraceKey string // 栈名
	LogInConsole  bool   // 输出控制台
}

/*
	存储一切有关Zinx框架的全局参数，供其他模块使用
	一些参数也可以通过 用户根据 zinx.yaml 来配置
*/
type obj struct {
	/*
		Server
	*/
	TCPServer ziface.IServer //当前Zinx的全局Server对象
	Host      string         //当前服务器主机IP
	TCPPort   int            //当前服务器主机监听端口号
	Name      string         //当前服务器名称
	// 详见[doublemsgid](https://github.com/sun-fight/zinx-websocket/tree/master/examples/doublemsgid)案例
	DoubleMsgID uint16 //(主子)双命令号模式(默认1单命令号模式)
	Env         string // develop production

	/*
		Zinx
	*/
	Version          string        //当前Zinx版本号
	MaxPacketSize    uint16        //读取数据包的最大值
	MaxConn          int           //当前服务器主机允许的最大链接个数
	WorkerPoolSize   uint32        //业务工作Worker池的数量
	MaxWorkerTaskLen uint32        //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen    uint32        //SendBuffMsg发送消息的缓冲最大长度
	HeartbeatTime    time.Duration //心跳间隔默认60秒,0=永不超时

	/*
		config file path
	*/
	ConfFilePath string

	/*
		zap
	*/
	ZapConfig zapConfig

	//mysql
	MysqlConfig mysqlConfig
	//redis
	RedisConfig redisConfig
	// 额外的配置 .ExtraConfig["web-host"]
	ExtraConfig map[string]interface{}
}

//PathExists 判断一个文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//Reload 读取用户的配置文件
func (g *obj) Reload() {
	if confFileExists, _ := PathExists(g.ConfFilePath); !confFileExists {
		fmt.Println("Config File " + g.ConfFilePath + " is not exist!!")
		return
	}

	v := viper.New()
	v.SetConfigFile(g.ConfFilePath)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		Glog.Info("config file changed:" + e.Name)
		if err := v.Unmarshal(&g); err != nil {
			Glog.Error("配置文件更新失败", zap.Error(err))
		}
	})
	if err := v.Unmarshal(&g); err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}

}

// InitObject 初始化全局配置
func InitObject() {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = "."
	}
	//初始化Object变量，设置一些默认值
	Object = &obj{
		Name:             "ZinxServerApp",
		Version:          "V0.11",
		TCPPort:          8999,
		Host:             "0.0.0.0",
		Env:              "production",
		DoubleMsgID:      1,
		MaxConn:          12000,
		MaxPacketSize:    4096,
		ConfFilePath:     pwd + "/conf/zinx.yaml",
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
		HeartbeatTime:    60,
		ZapConfig: zapConfig{
			Level:         "info",
			Format:        "console",
			Prefix:        "[zinx-websocket]",
			Director:      "log",
			LinkName:      "latest_log",
			ShowLine:      true,
			EncodeLevel:   "LowercaseColorLevelEncoder",
			StacktraceKey: "stacktrace",
			LogInConsole:  true,
		},
	}

	//NOTE: 从配置文件中加载一些用户配置的参数
	Object.Reload()
}
