package config

import (
	"strings"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"gopkg.in/cheggaaa/pb.v2"
)

// 定义文件名
var (
	// 指定字典名称
	TargetsList  = "targets_list.txt"
	UsernameDict = "username_dict.txt"
	PasswordDict = "password_dict.txt"
	// 结果保存文件
	ResultFile = "crack_result.txt"
)

// 定义进度条
var (
	// 检测 task 是否存活的进度条
	ProcessBarActive *pb.ProgressBar
	// 爆破进度条
	ProcessBarScan *pb.ProgressBar
)

// 定义超时时间，经测试，如果为 3 会超时
var TimeOut = 5 * time.Second

// 定义默认并发数，经测试，并发量不宜过大，否则会降低命中率
var Concurrency = 50

// Debug 模式
var DebugModel bool = false

// 标记特定的任务是否爆破成功，成功的话不再尝试破解该用户
var SuccessHash sync.Map

// 将爆破的结果存入 cache 中，该 cache 库支持内存数据落盘
var CacheTarget *cache.Cache

// 启动时间
var StartTime time.Time

// 定义端口和协议
var (
	// 定义用来判断协议是否支持的 map
	SupportProtocols map[string]bool

	// 指定当前支持的端口和协议
	PortProtocol = map[int]string{
		22: "SSH",
	}
)

// 初始化函数
func init() {
	// 初始化 map
	SupportProtocols = make(map[string]bool)
	for _, protocol := range PortProtocol {
		// 给 PortProtocol 中的每一个端口协议打上一个 true 的标签
		SupportProtocols[strings.ToUpper(protocol)] = true
	}

	// 初始化 CacheTarget
	CacheTarget = cache.New(cache.NoExpiration, cache.DefaultExpiration)
}
