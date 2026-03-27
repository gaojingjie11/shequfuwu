package service

import (
	"fmt"

	"smartcommunity/internal/model"
)

func buildFallbackCommunityReport(report *model.AIReport) string {
	return fmt.Sprintf(`## 社区运营简报（近7日）

### 一、核心数据概览
- 报修新增：%d 条
- 未处理报修：%d 条
- 访客新增：%d 条
- 物业费缴费笔数：%d 笔
- 物业费收缴金额：%.2f 元

### 二、管理风险提示
- 若未处理报修持续高位，可能影响住户满意度并增加投诉风险。
- 若物业费缴费笔数或金额偏低，需关注催缴机制和账单触达效果。

### 三、管理建议
1. 对未处理报修执行分级处置，优先处理超时工单。
2. 对未缴费住户开展分层提醒（短信、电话、上门）。
3. 每周复盘报修闭环时效和缴费转化率，持续优化流程。

> 说明：当前为系统本地模板报告（AI 服务暂不可用时自动降级）。`, report.RepairNewCount, report.RepairPendingCount, report.VisitorNewCount, report.PropertyPaidCount, report.PropertyPaidAmount)
}
