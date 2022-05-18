package utils

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/toki-plus/tkscan/cracker/config"
	"github.com/toki-plus/tkscan/cracker/log"
	"github.com/toki-plus/tkscan/cracker/models"
	"github.com/toki-plus/tkscan/cracker/plugins"
	"github.com/toki-plus/tkscan/cracker/utils/hash"
	"github.com/urfave/cli"
	"gopkg.in/cheggaaa/pb.v2"
)

// 生成扫描任务，用 tasks、用户名、密码 初始化一个 models.TargetsModel 结构
func GenerateTask(tasks []models.TScanTasks, usernames []string, passwords []string) (targets []models.TargetsModel, targetNum int) {
	targets = make([]models.TargetsModel, 0)

	for _, username := range usernames {
		for _, password := range passwords {
			for _, task := range tasks {
				target := models.TargetsModel{Ip: task.Ip, Port: task.Port, Protocol: task.Protocol, Username: username, Password: password}
				targets = append(targets, target)
			}
		}
	}

	return targets, len(targets)
}

// 开始扫描任务
func RunTask(tasks []models.TargetsModel) {
	// 进度条
	config.ProcessBarScan = pb.StartNew(len(tasks))
	config.ProcessBarScan.SetTemplate(`{{ rndcolor "Scanning progress: " }} {{  percent . "[%.02f%%]" "[?]"| rndcolor}} {{ counters . "[%s/%s]" "[%s/?]" | rndcolor}} {{ bar . "「" "-" (rnd "ᗧ" "◔" "◕" "◷" ) "•" "」" | rndcolor }} {{rtime . | rndcolor}}`)

	// 开启协程
	wg := &sync.WaitGroup{}
	// 创建一个大小为 config.Concurrency 的 channel
	taskChan := make(chan models.TargetsModel, config.Concurrency)
	// 创建 config.Concurrency 个协程
	for i := 0; i < config.Concurrency; i++ {
		go crackPassword(taskChan, wg)
	}
	// 生产者，不断地往 taskChan 发送任务，直到 taskChan 阻塞
	for _, task := range tasks {
		wg.Add(1)
		taskChan <- task
	}
	// 关闭 taskChan
	close(taskChan)
	// 等待最大超时时间，如果超时返回 true
	_ = waitTimeout(wg, config.TimeOut)


	// 将内存中的爆破结果落盘
	{
		// 将爆破结果保存到一个 DB 文件中，DB文件的格式为 go-cache 库定义的格式
		_ = models.SaveResultToFile()
		// 打印爆破的结果的状态信息，如爆破的用时、爆破得到的有效弱口令的总数
		models.ResultTotal()
		// 将结果导出到一个 config 指定的 txt 文件中
		_ = models.DumpToFile(config.ResultFile)
	}
}

// 每个协程都从 channel 中读取数据后开始爆破并保存
func crackPassword(targetChan chan models.TargetsModel, wg *sync.WaitGroup) {
	for target := range targetChan {
		// 进度条
		config.ProcessBarScan.Increment()

		if config.DebugModel {
			log.Log.Debugf("checking: IP: %v, Port: %v, [%v], UserName: %v, Password: %v, goroutineNum: %v",
				target.Ip,
				target.Port,
				target.Protocol,
				target.Username,
				target.Password,
				runtime.NumGoroutine(),
			)
		}

		// 格式化 target
		k := fmt.Sprintf("%v-%v-%v", target.Ip, target.Port, target.Username)
		// 将 target 进行 hash
		h := hash.MakeTaskHash(k)
		// 这个 hash 已经爆破成功
		if hash.CheckTaskHash(h) {
			// 等待组的数量减 1
			wg.Done()
			continue
		}

		// 协议名称转大写
		protocol := strings.ToUpper(target.Protocol)

		// 获取协议对应的爆破函数
		function := plugins.TScanFuncMap[protocol]

		// 保存结果
		models.SaveResult(function(target))

		wg.Done()
	}
}

// 等待最大超时时间，如果超时返回 true
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	// 正常返回 false
	case <-c:
		return false
	// 超时返回 true
	case <-time.After(timeout):
		return true
	}
}

// util.Scan 是 Scan 对象的 Action，对命令行传入的参数进行处理
// 如果传入了具体的参数就把 config 包中定义的默认值替换掉
func Scan(ctx *cli.Context) (err error) {

	if ctx.IsSet("debug_model") {
		config.DebugModel = ctx.Bool("debug_model")
	}

	if config.DebugModel {
		log.Log.Level = logrus.DebugLevel
	}

	if ctx.IsSet("timeout") {
		config.TimeOut = time.Duration(ctx.Int("timeout")) * time.Second
	}

	if ctx.IsSet("concurrency") {
		config.Concurrency = ctx.Int("concurrency")
	}

	if ctx.IsSet("username_dict") {
		config.UsernameDict = ctx.String("username_dict")
	}

	if ctx.IsSet("password_dict") {
		config.PasswordDict = ctx.String("password_dict")
	}

	if ctx.IsSet("result_file") {
		config.ResultFile = ctx.String("result_file")
	}

	config.StartTime = time.Now()

	usernameDict, uErr := ReadUsernameDict(config.UsernameDict)
	passwordDict, pErr := ReadPasswordDict(config.PasswordDict)
	targetsList := ReadTargetsList(config.TargetsList)
	aliveTargets := CheckAlive(targetsList)

	if uErr == nil && pErr == nil {
		tasks, _ := GenerateTask(aliveTargets, usernameDict, passwordDict)
		RunTask(tasks)
	}

	return err
}
