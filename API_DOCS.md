
# 智能社区 API 文档

**Base URL**: `http://43.138.85.114:81/api`

> 认证说明：
> 大部分接口需要在 Headers 中添加 `Authorization: Bearer <token>`
> Token 可通过登录接口获取

---

## 1. 用户认证 (Auth)

### 注册
- **URL**: `/v1/register`
- **Method**: `POST`
- **Body** (JSON):
  ```json
  {
    "username": "testuser",
    "password": "password123",
    "mobile": "13800138000"
  }
  ```

### 登录
- **URL**: `/v1/login`
- **Method**: `POST`
- **Body** (JSON):
  ```json
  {
    "username": "testuser",
    "password": "password123"
  }
  ```

### 获取当前用户信息
- **URL**: `/user/info`
- **Method**: `GET`
- **Headers**:
  - `Authorization`: `Bearer <token>`

---

## 2. 商品模块 (Mall)

### 获取商品列表
- **URL**: `/v1/products`
- **Method**: `GET`
- **Params**:
  - `page`: 1
  - `size`: 10
  - `category_id`: (可选) 1
  - `name`: (可选) 搜索关键词

### 获取商品分类
- **URL**: `/v1/categories`
- **Method**: `GET`

### 获取商品详情
- **URL**: `/v1/product/{id}`
- **Method**: `GET`
- **Example**: `/v1/product/1`

---

## 3. 评论模块 (Comment)

### 获取商品评论
- **URL**: `/v1/comments`
- **Method**: `GET`
- **Params**:
  - `product_id`: 1
  - `page`: 1
  - `size`: 10

### 发表评论 (需登录)
- **URL**: `/user/comment/create`
- **Method**: `POST`
- **Headers**:
  - `Authorization`: `Bearer <token>`
- **Body** (JSON):
  ```json
  {
    "product_id": 1,
    "content": "东西不错，物流很快！",
    "rating": 5
  }
  ```

---

## 4. 订单模块 (Order)

### 创建订单 (加入购物车/直接购买)
- **URL**: `/user/cart`
- **Method**: `POST`
- **Headers**:
  - `Authorization`: `Bearer <token>`
- **Body** (JSON):
  ```json
  {
    "product_id": 1,
    "quantity": 1
  }
  ```

### 获取购物车/订单列表
- **URL**: `/user/cart`
- **Method**: `GET`
- **Headers**:
  - `Authorization`: `Bearer <token>`

### 支付订单
- **URL**: `/user/cart/pay`
- **Method**: `POST`
- **Headers**:
  - `Authorization`: `Bearer <token>`
- **Body** (JSON):
  ```json
  {
    "order_id": 123
  }
  ```

---

## 5. 管理员模块 (Admin)

### 上传图片
- **URL**: `/user/upload`
- **Method**: `POST`
- **Headers**:
  - `Authorization`: `Bearer <token>`
  - `Content-Type`: `multipart/form-data`
- **Body**:
  - `file`: (选择文件)

### 获取数据大屏统计
- **URL**: `/v1/dashboard/stats`
- **Method**: `GET`
- **Headers**:
  - `Authorization`: `Bearer <token>`

### 搜索订单 (按用户ID)
- **URL**: `/user/orders/all`
- **Method**: `GET`
- **Headers**:
  - `Authorization`: `Bearer <token>`
- **Params**:
  - `page`: 1
  - `size`: 10
  - `user_id`: (可选) 123

### 新增车位
- **URL**: `/parking/admin/create`
- **Method**: `POST`
- **Headers**:
  - `Authorization`: `Bearer <token>`
- **Body** (JSON):
  ```json
  {
    "number": "A-888",
    "status": 0
  }
  ```
