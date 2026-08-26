package domain

import "strings"

func NormalizeFeeder(v string) string { return strings.ToUpper(strings.TrimSpace(v)) }
func NormalizeDevice(v string) string { return strings.ToUpper(strings.TrimSpace(v)) }
func SectionLabel(v string) string    { return strings.ToUpper(strings.TrimSpace(v)) }
func StatusLabel(v StepStatus) string {
	switch v {
	case Pending:
		return "待执行"
	case Dispatched:
		return "已派发"
	case Acknowledged:
		return "已确认"
	case Completed:
		return "已完成"
	case Failed:
		return "失败复核"
	}
	return "未知"
}
