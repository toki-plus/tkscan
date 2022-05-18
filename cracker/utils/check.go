package utils

// 检测从 file 中获得的 tasks 里的任务是否是活跃的

import (
	"fmt"
	"net"
	"sync"

	"github.com/toki-plus/tkscan/cracker/config"
	"github.com/toki-plus/tkscan/cracker/log"
	"github.com/toki-plus/tkscan/cracker/models"
	"gopkg.in/cheggaaa/pb.v2"
)

var AliveTask []models.TScanTasks

// 初始化 tasks 切片
func init() {
	AliveTask = make([]models.TScanTasks, 0)
}

func CheckAlive(tasks []models.TScanTasks) []models.TScanTasks {
	log.Log.Infoln("checking task active")
	// 进度条
	config.ProcessBarActive = pb.StartNew(len(tasks))
	config.ProcessBarActive.SetTemplate(`{{ rndcolor "Checking progress: " }} {{  percent . "[%.02f%%]" "[?]"| rndcolor}} {{ counters . "[%s/%s]" "[%s/?]" | rndcolor}} {{ bar . "「" "-" (rnd "ᗧ" "◔" "◕" "◷" ) "•" "」" | rndcolor}}  {{rtime . | rndcolor }}`)

	// 开启协程
	var wg sync.WaitGroup
	wg.Add(len(tasks))
	// 将每个 task 传入 check 函数进行检测
	for _, task := range tasks {
		// 闭包，自动加入协程
		go func(task models.TScanTasks) {
			// 告知当前goroutine完成
			defer wg.Done()
			SafeTask(check(task))
		}(task)
	}
	// 阻塞等待登记的goroutine完成
	wg.Wait()

	// 进度条
	config.ProcessBarActive.Finish()

	return AliveTask
}

// 使用 tcp 测试 task 连接的时延
func check(task models.TScanTasks) (bool, models.TScanTasks) {
	alive := false
	_, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", task.Ip, task.Port), config.TimeOut)
	// 如果指定时间内没有报错就是活跃的 task
	if err == nil {
		alive = true
	}

	// 进度条往前走
	config.ProcessBarActive.Increment()

	return alive, task
}

// 设置并发安全的互斥锁，保证同一时间只有一个 goroutine 可以访问共享资源
func SafeTask(alive bool, task models.TScanTasks) {
	// 定义互斥锁
	var mutex sync.Mutex
	if alive {
		// 获取互斥锁
		mutex.Lock()
		AliveTask = append(AliveTask, task)
		// 释放互斥锁
		mutex.Unlock()
	}
}
