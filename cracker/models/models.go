package models

// TargetsModel 用来传递从命令行接收到的目标参数
type TargetsModel struct {
	Ip       string
	Port     int
	Protocol string
	Username string
	Password string
}

// TScanResult 表示扫描结果，返回 bool
type TScanResult struct {
	TargetsModel TargetsModel
	Result       bool
}

// TScanTasks 表示扫描的服务类型
type TScanTasks struct {
	Ip       string
	Port     int
	Protocol string
}
