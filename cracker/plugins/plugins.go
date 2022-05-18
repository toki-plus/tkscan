package plugins

import (
	"github.com/toki-plus/tkscan/cracker/models"
)

// 定义一个用来实现插件功能的函数，函数应接收一个目标模型，返回一个 TScanResult
type TScanFunc func(targets models.TargetsModel) (result models.TScanResult, err error)

// 声明一个用来接收 TScanFunc 的 map，用来存储对应的插件函数名
var (
	TScanFuncMap map[string]TScanFunc
)

// 自动执行
func init() {
	// 初始化 TScanFuncMap
	TScanFuncMap = make(map[string]TScanFunc)
	// 实现一个 SSH 扫描函数
	TScanFuncMap["SSH"] = TScanSsh
}
