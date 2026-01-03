package main

import (
	"fmt"
	"os"
	"smartcommunity/internal/config" // 引入新写的 config 包
	"smartcommunity/internal/global"
	"smartcommunity/internal/middleware"
	"smartcommunity/internal/model"
	"smartcommunity/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// cmd/main.go
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	config.Init(env)

	// 2. 初始化数据库和Redis
	// 直接从 config.Conf 中获取 yaml 里的值
	global.InitDB(config.Conf.DB.DSN)
	global.InitRedis(config.Conf.Redis.Addr, "") // 假设yaml里没配密码，暂时传空
	global.InitMinio(config.Conf.MinIO)
	// 自动迁移所有表
	global.DB.AutoMigrate(
		&model.SysUser{},
		&model.SysRole{},
		&model.SysMenu{},
		&model.SysUserRole{},
		&model.SysRoleMenu{},

		&model.Product{},
		&model.ProductCategory{},
		&model.HotProduct{},
		&model.Promotion{},
		&model.Store{},
		&model.StoreProduct{},
		&model.Cart{},
		&model.Order{},
		&model.OrderItem{},
		&model.Notice{},
		&model.NoticeRead{},
		&model.Repair{},
		&model.Visitor{},
		&model.Parking{},
		&model.PropertyFee{},
		&model.Favorite{},
		&model.SysTransaction{},
	)

	// 3. 启动 Gin
	r := gin.Default()

	// 添加 CORS 中间件（重要！）
	r.Use(middleware.CORS())

	router.InitRouter(r)

	port := config.Conf.Server.Port
	if port == "" {
		port = "8080"
	}
	fmt.Printf("服务启动在 :%s\n", port)
	r.Run(":" + port)
}
