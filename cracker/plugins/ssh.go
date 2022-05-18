package plugins

import (
	"fmt"

	"github.com/toki-plus/tkscan/cracker/config"
	"github.com/toki-plus/tkscan/cracker/log"
	"github.com/toki-plus/tkscan/cracker/models"
	"golang.org/x/crypto/ssh"
)

func TScanSsh(target models.TargetsModel) (result models.TScanResult, err error) {
	// 先将传入的参数设入返回的结果
	result.TargetsModel = target
	// 配置扫描选项
	config := &ssh.ClientConfig{
		User: target.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(target.Password),
		},
		Timeout:         config.TimeOut,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 测试连接
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", target.Ip, target.Port), config)
	// 连接错误
	if err != nil {
		return result, err
	}
	log.Log.Infof("%v:%v 连接成功", target.Username, target.Password)

	// 连接成功，创建会话，执行命令
	session, err := client.NewSession()
	// 创建会话错误
	if err != nil {
		return result, err
	}
	log.Log.Infof("%v:%v 创建会话成功", target.Username, target.Password)


	// 延迟关闭连接，闭包匿名函数
	defer func() {
		// 关闭会话
		if session != nil {
			_ = session.Close()
		}
		if client != nil {
			_ = client.Close()
		}
	}()


	// 执行命令，这里的错误与上面的类型不同，需要重型定义
	err = session.Run("whoami")
	// 执行命令失败
	if err != nil {
		return result, err
	}
	log.Log.Infof("%v:%v 执行命令成功", target.Username, target.Password)


	// 成功，设置结果为 true
	result.Result = true


	// 返回结果
	return result, err
}
