package main

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

func TkScanSsh(ip string, port int, username, password string) (result string) {
	// ssh 包的配置选项
	// ClientConfig 结构用于配置客户端。在传递给 ssh 函数后，不能对其进行修改。
	config := &ssh.ClientConfig{
		User: username,
		// 设置 ssh 认证的方式为密码认证
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout: 3 * time.Second,
		// 在加密握手期间调用 HostKeyCallback 以验证服务器的主机密钥。客户端配置必须提供此回调才能成功连接。
		// 函数 UnsecureIgnoreHostKey 或 FixedHostKey 可用于简化的主机密钥检查。
		// UnsecureIgnoreHostKey 返回一个可用于 ClientConfig 的函数。HostKeyCallback 接受任何主机密钥。它不应用于生产代码。
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 测试连接
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", ip, port), config)
	// 连接失败
	if err != nil {
		fmt.Printf("connect failed, err: %v", err)
		return
	}
	// 没有错误，开启一个会话执行命令
	// 延迟关闭客户端
	defer client.Close()
	// 为本客户端开启一个会话，这个会话可以远程执行命令，底层会开启一个通道进行数据传输
	session, err := client.NewSession()
	// 执行命令，Run 在远程主机开启一个 cmd，可以执行一条命令，失败返回结果的错误信息
	errRet := session.Run("whoami")
	// 命令执行失败
	if errRet != nil {
		fmt.Printf("run command failed, err: %v", errRet)
		return
	}
	// 没有错误，返回结果给调用者，并关闭会话
	defer session.Close()
	return "open"
}

func main() {
	// 扫描 ssh 端口开放情况需要的参数
	ip := "192.168.80.134"
	// 端口用整数
	port := 22
	username := "root"
	password := "toki"

	// 传入参数到函数中执行
	result := TkScanSsh(ip, port, username, password)
	fmt.Printf("ssh => %v", result)
}
