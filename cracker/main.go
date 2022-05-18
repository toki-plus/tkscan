package main

import (
	"os"
	"runtime"

	"github.com/toki-plus/tkscan/cracker/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "tkscan-cracker"
	app.Usage = "A password cracker"
	app.Author = "Toki"
	app.Commands = []cli.Command{cmd.Scan}
	app.Flags = append(app.Flags, cmd.Scan.Flags...)

	// 入口点
	err := app.Run(os.Args)
	_ = err
}

func init() {
	// 设置当前可用的 CPU 数为同时可执行的最大 CPU 数
	runtime.GOMAXPROCS(runtime.NumCPU())
}