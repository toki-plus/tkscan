package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/toki-plus/tkscan/cracker/config"
	"github.com/toki-plus/tkscan/cracker/log"
	"github.com/toki-plus/tkscan/cracker/models"
)

// 将 targets_list.txt 中的目标加入任务组
func ReadTargetsList(fileName string) (tasks []models.TScanTasks) {
	// 读取文件
	targetsListFile, err := os.Open(fileName)
	// 排除
	if err != nil {
		fmt.Printf("Open file failed, err: %v", err)
	}
	// 延迟关闭文件
	defer func() {
		if targetsListFile != nil {
			_ = targetsListFile.Close()
		}
	}()

	// 通过标准库的 bufio 包逐行读取，然后用 strings 包的 Split 分割
	scanner := bufio.NewScanner(targetsListFile)
	// 按行分割，扫描开始后调用 Split 会报错
	scanner.Split(bufio.ScanLines)

	// 开始扫描，注：这里怎么判断扫描到末尾了？
	for scanner.Scan() {
		line := scanner.Text()
		// 如果扫到空行，将进继续扫描下一行，注：应该不是这个意思，扫描完退出不应该用 break 吗？
		if line == "" {
			continue
		}
		// 删除字符串的首尾空白符号
		task := strings.TrimSpace(line)
		// 以冒号分割字符串成 ip 和 端口号|协议，则通过协议来判断服务
		ipPort := strings.Split(task, ":")
		ip := ipPort[0]
		portProtocol := ipPort[1]
		// 以竖杠分割端口号和协议
		tmpPortProtocol := strings.Split(portProtocol, "|")
		// 如果 targets_list 写入的格式为 IP:PORT|PROTOCOL
		if len(tmpPortProtocol) == 2 {
			// 将端口号转成整型
			port, err := strconv.Atoi(tmpPortProtocol[0])
			if err != nil {
				fmt.Printf("str to int failed, err: %v", err)
			}
			// 将协议转成大写
			protocal := strings.ToUpper(tmpPortProtocol[1])
			// 判断协议是否支持
			if config.SupportProtocols[protocal] {
				// 协议支持，将上面三个参数传入 TScanTasks
				task := models.TScanTasks{Ip: ip, Port: port, Protocol: protocal}
				// 将扫描任务加入切片
				tasks = append(tasks, task)
			}
		} else {
			// 如果 targets_list 写入的格式为 IP:PORT，则通过端口来判断服务
			port, err := strconv.Atoi(tmpPortProtocol[0])
			if err != nil {
				fmt.Printf("str to int failed, err: %v", err)
			}
			// 判断是否能从 PortProtocol 中取出当前端口对应的协议
			protocol, ok := config.PortProtocol[port]
			// 如果有对应的协议，且协议有对应的插件功能，则证明支持该端口
			if ok && config.SupportProtocols[protocol] {
				// 协议支持，将上面三个参数传入 TScanTasks
				task := models.TScanTasks{Ip: ip, Port: port, Protocol: protocol}
				// 将扫描任务加入切片
				tasks = append(tasks, task)
			}
		}
	}

	// 返回 tasks 切片
	return tasks
}

// 传入用户名字典的名字，返回装着所有用户的切片
func ReadUsernameDict(usernameDict string) (usernames []string, err error) {
	file, err := os.Open(usernameDict)
	if err != nil {
		log.Log.Fatalf("Open username dict file err: %v", err)
	}

	// 延迟关闭文件，为啥要这样写？
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		// 去除空格
		username := strings.TrimSpace(scanner.Text())
		if username != "" {
			usernames = append(usernames, username)
		}
	}

	return usernames, err
}

// 传入密码字典的名字，返回装着所有用户的切片
func ReadPasswordDict(passwordDict string) (passwords []string, err error) {
	file, err := os.Open(passwordDict)
	if err != nil {
		log.Log.Fatalf("Open password dict filed err: %v", err)
	}

	// 延迟关闭文件
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		// 去除空格
		password := strings.TrimSpace(scanner.Text())
		if password != "" {
			passwords = append(passwords, password)
		}
	}

	return passwords, err
}
