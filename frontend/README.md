# 智慧社区前端项目

基于 Vue 3 + Vite 构建的现代化智慧社区管理系统前端应用，采用清新简洁的设计风格。

## ✨ 特性

- 🎨 **清新设计** - 使用现代化UI设计，配色清新，界面简洁
- 📱 **响应式布局** - 完美适配PC和移动端
- 🚀 **快速开发** - 基于 Vite 构建，开发体验极佳
- 🔐 **权限管理** - 完整的用户认证和权限控制
- 🛒 **商城功能** - 商品浏览、购物车、订单管理
- 🏘️ **社区服务** - 公告通知、报修投诉、访客登记等

## 🛠️ 技术栈

- **框架**: Vue 3.5 (Composition API)
- **构建工具**: Vite 6.0
- **路由**: Vue Router 4.5
- **状态管理**: Pinia 2.3
- **HTTP客户端**: Axios 1.7
- **日期处理**: Day.js 1.11
- **样式**: 原生 CSS + CSS Variables

## 📦 项目结构

```
frontend/
├── src/
│   ├── api/              # API接口
│   │   ├── auth.js       # 认证相关
│   │   ├── product.js    # 商品相关
│   │   ├── order.js      # 订单相关
│   │   └── service.js    # 社区服务
│   ├── assets/
│   │   └── styles/       # 全局样式
│   ├── components/       # 组件
│   │   └── layout/       # 布局组件
│   ├── router/           # 路由配置
│   ├── stores/           # 状态管理
│   │   ├── user.js       # 用户状态
│   │   └── cart.js       # 购物车状态
│   ├── utils/            # 工具函数
│   │   └── request.js    # HTTP请求封装
│   ├── views/            # 页面
│   │   ├── auth/         # 登录注册
│   │   ├── home/         # 首页
│   │   ├── mall/         # 商城
│   │   ├── order/        # 订单
│   │   ├── service/      # 社区服务
│   │   ├── profile/      # 个人中心
│   │   └── admin/        # 管理后台
│   ├── App.vue
│   └── main.js
├── index.html
├── package.json
└── vite.config.js
```

## 🚀 快速开始

### 1. 安装依赖

```bash
cd frontend
npm install
```

### 2. 启动开发服务器

```bash
npm run dev
```

项目将在 http://localhost:5173 启动

### 3. 构建生产版本

```bash
npm run build
```

## 🎯 功能模块

### 用户端功能

#### 1. 用户认证
- ✅ 用户注册
- ✅ 用户登录
- ✅ 忘记密码

#### 2. 首页
- ✅ 横幅展示
- ✅ 快捷入口
- ✅ 最新公告列表

#### 3. 商城模块
- ✅ 商品列表（分页、搜索）
- ✅ 商品详情
- ✅ 加入购物车
- ✅ 购物车管理
- ✅ 商品收藏

#### 4. 订单模块
- ✅ 订单创建
- ✅ 订单列表
- ✅ 订单支付
- ✅ 订单取消

#### 5. 社区服务
- ✅ 公告通知浏览
- ✅ 报修投诉提交
- ✅ 访客登记申请
- ✅ 历史记录查看

#### 6. 个人中心
- ✅ 个人信息查看
- ✅ 余额显示
- ✅ 快捷菜单

### 管理端功能

#### 7. 管理后台
- ✅ 管理入口界面
- 📝 用户管理（规划中）
- 📝 商品管理（规划中）
- 📝 订单管理（规划中）
- 📝 门店管理（规划中）
- 📝 数据统计（规划中）

## 🎨 设计系统

### 配色方案

```css
--primary-color: #00b894;      /* 清新绿 */
--secondary-color: #0984e3;    /* 天蓝色 */
--success-color: #00b894;      /* 成功-绿色 */
--warning-color: #fdcb6e;      /* 警告-黄色 */
--danger-color: #ff7675;       /* 危险-红色 */
--text-primary: #2d3436;       /* 主要文字 */
--text-secondary: #636e72;     /* 次要文字 */
--bg-color: #f8f9fa;           /* 背景色 */
```

### 设计原则

- **卡片式布局**: 所有内容使用圆角卡片呈现
- **充足留白**: 保持舒适的视觉间距
- **微动画**: hover、点击等交互动画
- **响应式**: 完美适配各种屏幕尺寸

## 🔌 API对接

### 后端API地址

开发环境: `http://localhost:8080/api/v1`

### 主要接口

- 认证: `/login`, `/register`
- 商品: `/products`, `/product/:id`
- 订单: `/order/create`, `/order/list`, `/order/pay`
- 购物车: `/cart/add`, `/cart/list`
- 公告: `/notices`, `/notice/:id`
- 报修: `/repair/create`, `/repair/list`

## 📝 测试账号

```
管理员:
用户名: admin / admin@example.com
密码: 123456

普通用户:
手机号: 13800000002
密码: 123456
```

## 🌐 浏览器支持

- Chrome >= 90
- Firefox >= 88
- Safari >= 14
- Edge >= 90

## 📄 开发说明

### 添加新页面

1. 在 `src/views` 下创建页面组件
2. 在 `src/router/index.js` 添加路由配置
3. 如需API，在 `src/api` 下添加接口函数

### 状态管理

使用 Pinia 进行状态管理，已创建：
- `user` store - 用户状态
- `cart` store - 购物车状态

### HTTP请求

所有请求通过 `src/utils/request.js` 统一处理，自动添加Token

## 🤝 贡献

欢迎提交 Issue 和 Pull Request

## 📜 许可证

MIT License
