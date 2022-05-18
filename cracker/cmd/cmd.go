package cmd

import (
	"github.com/toki-plus/tkscan/cracker/utils"
	"github.com/urfave/cli"
)

// 根据预设的命令行参数，定义相应的 cli.Command 对象
var Scan = cli.Command{
	Name:        "crack",
	Usage:       "./tkscan-cracker crack -i targets_list.txt -u username_dict.txt -p password_list.txt -t 5 -c 50 -d true -o crack_result.txt",
	Action:      utils.Scan,
	// 命令行参数
	Flags: []cli.Flag{
		boolFlag("debug_model, d", "-d true"),
		intFlag("timeout, t", 5, "-t 5"),
		intFlag("concurrency, c", 50, "-c 50"),
		stringFlag("targets_list, i", "targets_list.txt", "-i targets_list.txt"),
		stringFlag("username_dict, u", "username_dict.txt", "-u username_dict.txt"),
		stringFlag("password_list, p", "password_dict.txt", "-p password_list.txt"),
		stringFlag("crack_result, o", "crack_result.txt", "-o crack_result.txt"),
	},
}

func boolFlag(name, usage string) cli.BoolFlag {
	return cli.BoolFlag{
		Name:  name,
		Usage: usage,
	}
}

func intFlag(name string, value int, usage string) cli.IntFlag {
	return cli.IntFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func stringFlag(name, value, usage string) cli.StringFlag {
	return cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}
