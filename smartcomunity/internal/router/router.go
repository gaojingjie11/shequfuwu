package router

import (
	"smartcommunity/internal/controller"
	"smartcommunity/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	// 初始化 Handler
	uploadHandler := controller.UploadHandler{}
	userHandler := controller.UserHandler{}
	productHandler := controller.ProductHandler{}
	cartHandler := controller.CartHandler{}
	orderHandler := controller.OrderHandler{} // 新增
	noticeHandler := controller.NoticeHandler{}
	repairHandler := controller.RepairHandler{}
	financeHandler := controller.FinanceHandler{}   // 新增
	securityHandler := controller.SecurityHandler{} // 新增
	favoriteHandler := controller.FavoriteHandler{} // 新增
	// financeHandler, securityHandler, favoriteHandler already initialized

	storeHandler := controller.StoreHandler{}         // 初始化
	adminHandler := controller.AdminHandler{}         // 新增 后台管理
	marketingHandler := controller.MarketingHandler{} // 新增 营销管理

	// 开放接口 (不需要登录)
	apiGroup := r.Group("/api/v1")
	{ //用户
		apiGroup.POST("/register", userHandler.Register)
		apiGroup.POST("/login", userHandler.Login)
		apiGroup.POST("/send_code", userHandler.SendCode)             // 发送验证码
		apiGroup.POST("/login_code", userHandler.LoginCode)           // 验证码登录
		apiGroup.POST("/forget_password", userHandler.ForgetPassword) // 新增

		// 商品 (通常游客也可以看商品，所以放在开放接口里)
		apiGroup.GET("/products", productHandler.List)      // 列表
		apiGroup.GET("/product/:id", productHandler.Detail) //
		apiGroup.POST("/cart", productHandler.Create)

		// --- 【新增】门店列表 ---
		apiGroup.GET("/stores", storeHandler.List)

		// --- 【新增】商品分类 ---
		apiGroup.GET("/categories", productHandler.GetCategories)

		// --- 【新增】数据大屏 ---
		apiGroup.GET("/dashboard/stats", adminHandler.GetDashboardStats)

		// --- 【新增】商品评论 ---
		commentHandler := controller.CommentHandler{}
		apiGroup.GET("/comments", commentHandler.List) // 公开查看评论
	}

	// 2. 私有接口 (Private - 需要登录)
	// 使用 middleware.JWTAuth() 保护这个组
	private := r.Group("/api/v1")
	private.Use(middleware.JWTAuth())
	{
		// 购物车相关
		private.POST("/cart/add", cartHandler.Add)      // 添加
		private.GET("/cart/list", cartHandler.List)     // 列表
		private.DELETE("/cart/:id", cartHandler.Delete) // 删除
		private.POST("/cart/:id", cartHandler.Update)   // 修改数量
		// 订单相关
		private.POST("/order/create", orderHandler.Create)     // 下单
		private.GET("/order/list", orderHandler.List)          // 订单列表
		private.POST("/order/pay", orderHandler.Pay)           // 支付订单
		private.GET("/order/admin/list", orderHandler.ListAll) // 管理员看所有订单
		private.POST("/order/ship", orderHandler.Ship)         // 管理员发货
		private.POST("/order/receive", orderHandler.Receive)   // 用户确认收货
		private.POST("/order/cancel", orderHandler.Cancel)     // 用户取消

		// 报事报修
		private.POST("/repair/create", repairHandler.Create) // 提交
		private.GET("/repair/list", repairHandler.List)      // 查看历史..

		// --- 支付中心 ---
		private.POST("/finance/pay", financeHandler.Pay)                      // 统一支付
		private.GET("/property/list", financeHandler.ListPropertyFee)         // 物业费列表
		private.POST("/finance/recharge", financeHandler.Recharge)            // 充值
		private.POST("/finance/transfer", financeHandler.Transfer)            // 转账
		private.GET("/finance/transactions", financeHandler.ListTransactions) //获取交易流水
		// --- 营销管理 ---
		private.POST("/marketing/promotion/create", marketingHandler.Create)
		private.GET("/marketing/promotion/list", marketingHandler.List)
		private.DELETE("/marketing/promotion/:id", marketingHandler.Delete)

		// --- 门店管理 ---
		private.POST("/store/create", storeHandler.Create)
		private.POST("/store/update", storeHandler.Update)
		private.DELETE("/store/:id", storeHandler.Delete)
		private.POST("/store/bind_product", storeHandler.BindProduct)

		// --- 管理员接口 ---
		private.POST("/product/create", productHandler.Create) // 发布商品
		private.POST("/product/update", productHandler.Update) // 更新商品
		private.DELETE("/product/:id", productHandler.Delete)  // 删除商品
		private.GET("/product/rank", productHandler.GetRank)   // 销量排

		// --- 新增：安保管理 (访客 & 车位) ---
		private.POST("/visitor/create", securityHandler.CreateVisitor) // 访客登记
		private.GET("/visitor/list", securityHandler.ListVisitor)      // 访客记录
		private.GET("/parking/my", securityHandler.MyParking)          // 我的车位
		private.POST("/parking/bind", securityHandler.BindCar)
		private.GET("/visitor/admin/list", securityHandler.ListAllVisitor) // 管理员看列表
		private.POST("/visitor/audit", securityHandler.AuditVisitor)       // 审核通过/拒绝

		// --- 【新增】通用上传 ---
		private.POST("/upload", uploadHandler.UploadFile)

		// --- 【新增】管理员-公告管理 ---
		private.POST("/notice/create", noticeHandler.Create) // 发布
		private.DELETE("/notice/:id", noticeHandler.Delete)  // 删除
		private.POST("/notice/read/:id", noticeHandler.Read) // 标记已读 (用户)

		// --- 【新增】管理员-报修管理 ---
		private.GET("/repair/admin/list", repairHandler.ListAll) // 查看所有
		private.POST("/repair/process", repairHandler.Process)   // 处理/反馈

		// --- 收藏夹 ---
		private.POST("/favorite/add", favoriteHandler.Add)       // 收藏
		private.POST("/favorite/delete", favoriteHandler.Delete) // 取消
		private.GET("/favorites", favoriteHandler.List)          // 列表
		private.GET("/favorite/check", favoriteHandler.Check)    // 检查是否收藏
		// --- 用户个人中心 ---
		private.POST("/logout", userHandler.Logout)                       // 退出登录 (新增)
		private.POST("/user/update", userHandler.Update)                  // 修改头像/昵称
		private.POST("/user/change_password", userHandler.ChangePassword) // 修改密码
		private.GET("/user/info", userHandler.Info)                       // 获取最新信息

		// --- 车位管理 (Admin) 新增 ---
		private.GET("/parking/admin/list", securityHandler.GetAllParking)
		private.GET("/parking/admin/stats", securityHandler.GetParkingStats)
		private.POST("/parking/admin/assign", securityHandler.AssignParking)
		private.POST("/parking/admin/create", securityHandler.CreateParking) // 新增车位

		// --- 物业费管理 (Admin) 新增 ---
		private.POST("/property/admin/create", financeHandler.CreatePropertyFee)
		private.GET("/property/admin/list", financeHandler.ListAllPropertyFees)

		// --- 系统管理 (RBAC & Logs) ---
		private.POST("/admin/role/create", adminHandler.CreateRole)
		private.GET("/admin/role/list", adminHandler.ListRoles)
		private.POST("/admin/menu/create", adminHandler.CreateMenu)
		private.GET("/admin/menu/list", adminHandler.ListMenus)
		private.POST("/admin/role/bind_menu", adminHandler.BindRoleMenu)

		private.GET("/admin/user/list", adminHandler.ListUsers)     // 用户列表
		private.POST("/admin/user/freeze", adminHandler.FreezeUser) // 冻结用户
		// 新增用户管理接口
		private.POST("/admin/user/assign_role", adminHandler.AssignRole)
		private.POST("/admin/user/update_balance", adminHandler.UpdateUserBalance)

		// --- 【新增】商品评论 (Auth) ---
		commentHandler := controller.CommentHandler{}
		private.POST("/comment/create", commentHandler.Create)
	}

	public := r.Group("/api/v1")
	{
		// ... 其他公开接口 ...
		public.GET("/notices", noticeHandler.List)      // 列表
		public.GET("/notice/:id", noticeHandler.Detail) // 详情

	}
}
