package global

import (
	"context"
	"log"
	"smartcommunity/internal/config" // 引入 config 包
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// DB 全局 MySQL 数据库连接对象
	DB *gorm.DB
	// RDB 全局 Redis 连接对象
	RDB *redis.Client

	MinioClient *minio.Client // 新增全局客户端
)

// InitDB 初始化 MySQL 连接
// dsn: 数据库连接字符串 (如 "root:root@tcp(127.0.0.1:3306)/dbname?...")
func InitDB(dsn string) {
	var err error

	// 配置 GORM 的日志模式，打印 SQL 语句方便调试
	dbConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	DB, err = gorm.Open(mysql.Open(dsn), dbConfig)
	if err != nil {
		log.Fatalf("MySQL 连接失败: %v", err)
	}

	// 获取底层的 sqlDB 对象，用于设置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("获取底层 SQL DB 失败: %v", err)
	}

	// 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(100)
	// 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ MySQL 连接成功")
}

// InitRedis 初始化 Redis 连接
// addr: Redis 地址 (如 "localhost:6379")
// password: Redis 密码 (如果没有则传空字符串)
func InitRedis(addr string, password string) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // 如果没有密码，传 ""
		DB:       0,        // 使用默认 DB 0
	})

	// 测试连接 (Ping)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("❌ Redis 连接失败: %v", err)
	}

	log.Println("✅ Redis 连接成功")
}

// --- 新增：初始化 MinIO ---
func InitMinio(conf config.MinIOConfig) {
	var err error
	// 初始化 MinIO 客户端
	MinioClient, err = minio.New(conf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKey, conf.SecretKey, ""),
		Secure: conf.UseSSL,
	})
	if err != nil {
		log.Fatalf("❌ MinIO 客户端创建失败: %v", err)
	}

	// 检查桶是否存在，不存在则创建 (可选，为了严谨建议加上)
	ctx := context.Background()
	exists, err := MinioClient.BucketExists(ctx, conf.Bucket)
	if err != nil {
		log.Printf("⚠️ 检查 Bucket 失败: %v", err)
	} else if !exists {
		err = MinioClient.MakeBucket(ctx, conf.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("⚠️ 创建 Bucket 失败: %v", err)
		} else {
			log.Printf("✅ Bucket '%s' 创建成功", conf.Bucket)
			// 注意：实际生产中还需要设置 Policy 为 Public 才能直接访问图片
			// 这里假设你已经在 MinIO Console 设置好了 Public 权限
		}
	}

	log.Println("✅ MinIO 连接成功")
}
