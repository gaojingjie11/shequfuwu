package config

import (
	"log"

	"github.com/spf13/viper"
)

// 定义全局配置变量，外部直接调用 config.Conf.DB.DSN 即可
var Conf *Config

// Config 结构体映射 dev.yaml 的结构
type Config struct {
	Server ServerConfig `mapstructure:"server"` // 新增 Server 配置
	DB     DBConfig     `mapstructure:"db"`
	Redis  RedisConfig  `mapstructure:"redis"`
	MinIO  MinIOConfig  `mapstructure:"minio"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DBConfig struct {
	DSN string `mapstructure:"dsn"`
}

type RedisConfig struct {
	Addr string `mapstructure:"addr"`
}

// 新增 MinIO 配置结构
type MinIOConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	UseSSL    bool   `mapstructure:"use_ssl"`
}

// Init 初始化配置
func Init(env string) {
	// Init 支持传入环境参数，比如 "dev" 或 "prod"

	fileName := "dev"
	if env != "" {
		fileName = env
	}
	viper.SetConfigName(fileName)
	// ... 其他代码不变
	viper.SetConfigType("yaml")     // 配置文件类型
	viper.AddConfigPath("./config") // 查找配置文件的路径 (相对于项目根目录)
	// 如果你在IDE里运行可能需要加这一行，防止路径不对：
	// viper.AddConfigPath("../config")

	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 将读取的配置反序列化到 Conf 变量中
	if err := viper.Unmarshal(&Conf); err != nil {
		log.Fatalf("配置解析失败: %v", err)
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		log.Fatalf("配置解析失败: %v", err)
	}
	log.Println("配置加载成功")
	log.Println("配置加载成功")
}
