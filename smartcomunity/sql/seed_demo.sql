-- smartcommunity demo seed data
-- Usage:
-- 1) Ensure tables are created by running backend once (AutoMigrate).
-- 2) Import this SQL into smart_community database.
-- 3) Login demo accounts (all passwords are: 123456).

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ------------------------------------------------------------
-- 基础权限数据
-- ------------------------------------------------------------
INSERT INTO `sys_role` (`id`, `name`, `code`, `remark`, `created_at`) VALUES
(1, '系统管理员', 'admin', '全局管理权限', '2026-03-27 20:00:00'),
(2, '物业管理员', 'property', '物业与报修管理', '2026-03-27 20:00:00'),
(3, '普通用户', 'user', '居民用户', '2026-03-27 20:00:00')
ON DUPLICATE KEY UPDATE
`name`=VALUES(`name`), `code`=VALUES(`code`), `remark`=VALUES(`remark`);

INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `path`, `component`, `sort`, `type`, `created_at`) VALUES
(1, 0, '仪表盘', '/admin/dashboard', 'views/admin/Dashboard.vue', 1, 1, '2026-03-27 20:00:00'),
(2, 0, '订单管理', '/admin/order', 'views/admin/order/OrderList.vue', 2, 1, '2026-03-27 20:00:00'),
(3, 0, '物业管理', '/admin/property', 'views/admin/property/PropertyFeeList.vue', 3, 1, '2026-03-27 20:00:00'),
(4, 0, 'AI报表', '/admin/ai-report', 'views/admin/AIReport.vue', 4, 1, '2026-03-27 20:00:00')
ON DUPLICATE KEY UPDATE
`name`=VALUES(`name`), `path`=VALUES(`path`), `component`=VALUES(`component`), `sort`=VALUES(`sort`), `type`=VALUES(`type`);

INSERT INTO `sys_role_menu` (`id`, `role_id`, `menu_id`) VALUES
(1, 1, 1), (2, 1, 2), (3, 1, 3), (4, 1, 4),
(5, 2, 1), (6, 2, 3), (7, 2, 4)
ON DUPLICATE KEY UPDATE
`role_id`=VALUES(`role_id`), `menu_id`=VALUES(`menu_id`);

-- 密码均为 123456 (bcrypt, cost=14)
INSERT INTO `sys_user` (`id`, `username`, `password`, `real_name`, `mobile`, `age`, `gender`, `email`, `avatar`, `green_points`, `balance`, `role`, `status`, `created_at`, `updated_at`) VALUES
(1, 'admin', '$2a$14$5Kc3aGL3vFnv3LAUDOCEZOTKQWfqag5edtPgvsKuS0C2qePaD4i46', '系统管理员', '13800000001', 30, 1, 'admin@community.com', 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png', 5000, 8000.00, 'admin', 1, '2026-03-27 20:00:00', '2026-03-27 20:00:00'),
(2, 'property', '$2a$14$5Kc3aGL3vFnv3LAUDOCEZOTKQWfqag5edtPgvsKuS0C2qePaD4i46', '物业经理', '13800000002', 35, 1, 'property@community.com', 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png', 1200, 3000.00, 'property', 1, '2026-03-27 20:00:00', '2026-03-27 20:00:00'),
(3, 'user_demo', '$2a$14$5Kc3aGL3vFnv3LAUDOCEZOTKQWfqag5edtPgvsKuS0C2qePaD4i46', '测试住户', '13800000003', 28, 2, 'user@community.com', 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png', 860, 520.00, 'user', 1, '2026-03-27 20:00:00', '2026-03-27 20:00:00')
ON DUPLICATE KEY UPDATE
`username`=VALUES(`username`), `password`=VALUES(`password`), `real_name`=VALUES(`real_name`), `mobile`=VALUES(`mobile`),
`age`=VALUES(`age`), `gender`=VALUES(`gender`), `email`=VALUES(`email`), `avatar`=VALUES(`avatar`),
`green_points`=VALUES(`green_points`), `balance`=VALUES(`balance`), `role`=VALUES(`role`), `status`=VALUES(`status`), `updated_at`=VALUES(`updated_at`);

INSERT INTO `sys_user_role` (`id`, `user_id`, `role_id`) VALUES
(1, 1, 1), (2, 2, 2), (3, 3, 3)
ON DUPLICATE KEY UPDATE
`user_id`=VALUES(`user_id`), `role_id`=VALUES(`role_id`);

-- ------------------------------------------------------------
-- 商城与店铺
-- ------------------------------------------------------------
INSERT INTO `pms_store` (`id`, `name`, `address`, `phone`, `region`, `business_hours`) VALUES
(1, '智慧社区生活馆', '幸福路 100 号 1 栋底商', '0755-88886666', 'A区', '09:00-22:00')
ON DUPLICATE KEY UPDATE
`name`=VALUES(`name`), `address`=VALUES(`address`), `phone`=VALUES(`phone`), `region`=VALUES(`region`), `business_hours`=VALUES(`business_hours`);

INSERT INTO `pms_product_category` (`id`, `name`, `icon`, `sort`) VALUES
(1, '家清洗护', 'Brush', 1),
(2, '居家日用', 'HomeFilled', 2),
(3, '粮油食品', 'Goods', 3)
ON DUPLICATE KEY UPDATE
`name`=VALUES(`name`), `icon`=VALUES(`icon`), `sort`=VALUES(`sort`);

INSERT INTO `pms_product` (`id`, `category_name`, `name`, `description`, `price`, `original_price`, `stock`, `image_url`, `is_promotion`, `sales`, `status`, `created_at`, `category_id`) VALUES
(101, '家清洗护', '洗衣液 3kg', '低泡易漂，家庭常用款', 42.00, 50.00, 260, 'https://picsum.photos/seed/laundry/300/300', 1, 182, 1, '2026-03-27 20:00:00', 1),
(102, '居家日用', '抽纸 24包', '原木本色，4层加厚', 29.90, 39.90, 300, 'https://picsum.photos/seed/tissue/300/300', 1, 96, 1, '2026-03-27 20:00:00', 2),
(103, '粮油食品', '东北大米 10kg', '当季新米，真空包装', 79.90, 89.90, 180, 'https://picsum.photos/seed/rice/300/300', 0, 56, 1, '2026-03-27 20:00:00', 3)
ON DUPLICATE KEY UPDATE
`category_name`=VALUES(`category_name`), `name`=VALUES(`name`), `description`=VALUES(`description`),
`price`=VALUES(`price`), `original_price`=VALUES(`original_price`), `stock`=VALUES(`stock`), `image_url`=VALUES(`image_url`),
`is_promotion`=VALUES(`is_promotion`), `sales`=VALUES(`sales`), `status`=VALUES(`status`), `category_id`=VALUES(`category_id`);

INSERT INTO `pms_store_product` (`id`, `store_id`, `product_id`, `stock`) VALUES
(1, 1, 101, 260),
(2, 1, 102, 300),
(3, 1, 103, 180)
ON DUPLICATE KEY UPDATE
`store_id`=VALUES(`store_id`), `product_id`=VALUES(`product_id`), `stock`=VALUES(`stock`);

INSERT INTO `pms_hot_product` (`id`, `product_id`, `created_at`) VALUES
(1, 101, '2026-03-27 20:00:00'),
(2, 102, '2026-03-27 20:00:00')
ON DUPLICATE KEY UPDATE
`product_id`=VALUES(`product_id`);

-- ------------------------------------------------------------
-- 公告 / 报修 / 访客 / 车位 / 物业费
-- ------------------------------------------------------------
INSERT INTO `cms_notice` (`id`, `title`, `content`, `publisher`, `view_count`, `created_at`) VALUES
(1, '关于电梯保养的通知', '3月30日 09:00-12:00 对 1-3 栋电梯进行保养，请提前安排出行。', '物业中心', 35, '2026-03-25 09:00:00'),
(2, '垃圾分类积分活动', '上传垃圾分类照片可获得绿色积分，积分可用于商城和物业费抵扣。', '物业中心', 52, '2026-03-26 10:30:00'),
(3, '清明假期服务安排', '客服与安保正常值班，报修工单 24 小时响应。', '物业中心', 18, '2026-03-27 08:20:00')
ON DUPLICATE KEY UPDATE
`title`=VALUES(`title`), `content`=VALUES(`content`), `publisher`=VALUES(`publisher`), `view_count`=VALUES(`view_count`);

INSERT INTO `cms_repair` (`id`, `user_id`, `type`, `category`, `content`, `status`, `result`, `created_at`) VALUES
(1, 3, 1, '水电', '厨房水龙头漏水，夜间持续滴水。', 0, '', '2026-03-26 20:12:00'),
(2, 3, 3, '门禁', '单元门禁卡识别失败，无法刷卡进门。', 1, '已安排工程师明日上午上门处理', '2026-03-25 18:42:00'),
(3, 3, 4, '电梯', '2栋电梯有异响，晚高峰明显。', 2, '已更换导轨部件，运行正常', '2026-03-23 14:15:00')
ON DUPLICATE KEY UPDATE
`user_id`=VALUES(`user_id`), `type`=VALUES(`type`), `category`=VALUES(`category`), `content`=VALUES(`content`),
`status`=VALUES(`status`), `result`=VALUES(`result`);

INSERT INTO `cms_visitor` (`id`, `user_id`, `visitor_name`, `visitor_phone`, `reason`, `visit_time`, `status`, `audit_remark`, `created_at`) VALUES
(1, 3, '张三', '13900001111', '亲友来访', '2026-03-27 19:30:00', 1, '已放行', '2026-03-27 16:20:00'),
(2, 3, '李四', '13900002222', '送货上门', '2026-03-28 10:00:00', 0, '', '2026-03-27 18:05:00')
ON DUPLICATE KEY UPDATE
`user_id`=VALUES(`user_id`), `visitor_name`=VALUES(`visitor_name`), `visitor_phone`=VALUES(`visitor_phone`),
`reason`=VALUES(`reason`), `visit_time`=VALUES(`visit_time`), `status`=VALUES(`status`), `audit_remark`=VALUES(`audit_remark`);

INSERT INTO `cms_parking` (`id`, `parking_no`, `status`, `user_id`, `car_plate`) VALUES
(1, 'A-001', 1, 3, '粤B12345'),
(2, 'A-002', 0, 0, ''),
(3, 'B-008', 1, 2, '粤B88888')
ON DUPLICATE KEY UPDATE
`parking_no`=VALUES(`parking_no`), `status`=VALUES(`status`), `user_id`=VALUES(`user_id`), `car_plate`=VALUES(`car_plate`);

INSERT INTO `cms_property_fee` (`id`, `user_id`, `month`, `amount`, `used_points`, `used_balance`, `status`, `pay_time`) VALUES
(1, 3, '2026-03', 120.00, 800, 40.00, 1, '2026-03-27 12:05:00'),
(2, 3, '2026-04', 120.00, 0, 0.00, 0, NULL)
ON DUPLICATE KEY UPDATE
`user_id`=VALUES(`user_id`), `month`=VALUES(`month`), `amount`=VALUES(`amount`),
`used_points`=VALUES(`used_points`), `used_balance`=VALUES(`used_balance`), `status`=VALUES(`status`), `pay_time`=VALUES(`pay_time`);

-- ------------------------------------------------------------
-- 订单、收藏、流水、积分流水
-- ------------------------------------------------------------
INSERT INTO `oms_order` (`id`, `order_no`, `user_id`, `store_id`, `total_amount`, `used_points`, `used_balance`, `status`, `paid_at`, `created_at`) VALUES
(1, '17700000000000000001', 3, 1, 42.00, 0, 0.00, 0, NULL, '2026-03-27 21:52:00'),
(2, '17700000000000000002', 3, 1, 29.90, 200, 9.90, 1, '2026-03-27 20:10:00', '2026-03-27 20:00:00'),
(3, '17700000000000000003', 3, 1, 79.90, 0, 79.90, 2, '2026-03-26 19:00:00', '2026-03-26 18:30:00')
ON DUPLICATE KEY UPDATE
`order_no`=VALUES(`order_no`), `user_id`=VALUES(`user_id`), `store_id`=VALUES(`store_id`), `total_amount`=VALUES(`total_amount`),
`used_points`=VALUES(`used_points`), `used_balance`=VALUES(`used_balance`), `status`=VALUES(`status`), `paid_at`=VALUES(`paid_at`);

INSERT INTO `oms_order_item` (`id`, `order_id`, `product_id`, `price`, `quantity`) VALUES
(1, 1, 101, 42.00, 1),
(2, 2, 102, 29.90, 1),
(3, 3, 103, 79.90, 1)
ON DUPLICATE KEY UPDATE
`order_id`=VALUES(`order_id`), `product_id`=VALUES(`product_id`), `price`=VALUES(`price`), `quantity`=VALUES(`quantity`);

INSERT INTO `pms_favorite` (`id`, `user_id`, `product_id`, `created_at`) VALUES
(1, 3, 101, '2026-03-27 20:30:00'),
(2, 3, 102, '2026-03-27 20:35:00')
ON DUPLICATE KEY UPDATE
`user_id`=VALUES(`user_id`), `product_id`=VALUES(`product_id`);

INSERT INTO `sys_transaction` (`id`, `user_id`, `type`, `amount`, `related_id`, `remark`, `created_at`) VALUES
(1, 3, 3, 500.00, 0, '初始充值', '2026-03-27 19:00:00'),
(2, 3, 1, -9.90, 2, '支付订单 17700000000000000002', '2026-03-27 20:10:00'),
(3, 3, 2, -40.00, 1, '缴纳物业费 2026-03', '2026-03-27 12:05:00')
ON DUPLICATE KEY UPDATE
`user_id`=VALUES(`user_id`), `type`=VALUES(`type`), `amount`=VALUES(`amount`), `related_id`=VALUES(`related_id`), `remark`=VALUES(`remark`);

INSERT INTO `green_point_record` (`id`, `user_id`, `action`, `points`, `created_at`) VALUES
(1, 3, 'garbage_classification', 120, '2026-03-24 09:30:00'),
(2, 3, 'garbage_classification', 80, '2026-03-25 18:20:00'),
(3, 3, 'property_fee', -800, '2026-03-27 12:05:00'),
(4, 3, 'mall_consume', -200, '2026-03-27 20:10:00'),
(5, 1, 'garbage_classification', 500, '2026-03-26 08:30:00'),
(6, 2, 'garbage_classification', 300, '2026-03-26 09:10:00')
ON DUPLICATE KEY UPDATE
`user_id`=VALUES(`user_id`), `action`=VALUES(`action`), `points`=VALUES(`points`), `created_at`=VALUES(`created_at`);

-- ------------------------------------------------------------
-- AI 报表与聊天记录
-- ------------------------------------------------------------
INSERT INTO `cms_ai_report` (`id`, `repair_new_count`, `repair_pending_count`, `visitor_new_count`, `property_paid_count`, `property_paid_amount`, `report_summary`, `report_markdown`, `generated_by`, `created_at`, `updated_at`) VALUES
(1, 3, 2, 2, 1, 120.00,
'社区近7日运营数据分析报告',
'# 社区近7日运营数据分析报告\n\n## 1. 核心数据概览\n- 报修新增：3条（未处理2条）\n- 新增访客：2条\n- 物业费收缴：1笔，共120.00元\n\n## 2. 管理风险\n- 报修存在积压，存在投诉风险。\n- 物业费收缴率有提升空间。\n\n## 3. 建议\n- 48小时内清理未处理工单。\n- 对未缴费住户建立分层催收机制。\n- 对访客高峰时段增加门岗支持。',
1, '2026-03-27 21:00:00', '2026-03-27 21:00:00')
ON DUPLICATE KEY UPDATE
`repair_new_count`=VALUES(`repair_new_count`), `repair_pending_count`=VALUES(`repair_pending_count`),
`visitor_new_count`=VALUES(`visitor_new_count`), `property_paid_count`=VALUES(`property_paid_count`),
`property_paid_amount`=VALUES(`property_paid_amount`), `report_summary`=VALUES(`report_summary`),
`report_markdown`=VALUES(`report_markdown`), `generated_by`=VALUES(`generated_by`), `updated_at`=VALUES(`updated_at`);

INSERT INTO `cms_chat_messages` (`id`, `user_id`, `role`, `content`, `created_at`) VALUES
(1, 3, 'user', '帮我买洗衣液', '2026-03-27 21:30:00'),
(2, 3, 'assistant', '已为您找到在售洗衣液：洗衣液 3kg，价格42元。需要我为您下单吗？', '2026-03-27 21:30:03'),
(3, 3, 'user', '下单一瓶', '2026-03-27 21:30:15'),
(4, 3, 'assistant', '订单已创建，请确认支付并输入登录密码。', '2026-03-27 21:30:18')
ON DUPLICATE KEY UPDATE
`user_id`=VALUES(`user_id`), `role`=VALUES(`role`), `content`=VALUES(`content`), `created_at`=VALUES(`created_at`);

SET FOREIGN_KEY_CHECKS = 1;

