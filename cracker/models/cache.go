package models

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/toki-plus/tkscan/cracker/config"
	"github.com/toki-plus/tkscan/cracker/log"
	"github.com/toki-plus/tkscan/cracker/utils/hash"
)

// 初始化，将扫描获取到的结果序列化成字符串进行保存
func init() {
	gob.Register(TargetsModel{})
	gob.Register(TScanResult{})
}

// 将插件运行出的结果保存起来
func SaveResult(result TScanResult, err error) {
	// 结果没有报错，且结果为 true
	if err == nil && result.Result {

		// 格式化 target
		k := fmt.Sprintf("%v-%v-%v", result.TargetsModel.Ip, result.TargetsModel.Port, result.TargetsModel.Username)
		// 将 task 转成哈希
		h := hash.MakeTaskHash(k)
		// 把 task 的 hash 存储为 key，值为 true
		hash.SetTaskHash(h)


		// 从 cache 中寻找一个值
		_, found := config.CacheTarget.Get(k)
		// 没有发现
		if !found {
			log.Log.Infof("Ip: %v, Port: %v, Protocol: [%v], Username: %v, Password: %v",
				result.TargetsModel.Ip,
				result.TargetsModel.Port,
				result.TargetsModel.Protocol,
				result.TargetsModel.Username,
				result.TargetsModel.Password)
		}
		// 将爆破的结果添加到缓存，替换任何现有项目。
		// 如果持续时间为 0（DefaultExpiration），则使用缓存的默认过期时间。
		// 如果为 -1（NoExpiration），则该数据永不过期。
		config.CacheTarget.Set(k, result, cache.NoExpiration)
	}
}

// 打印爆破的结果的状态信息，如爆破的用时、爆破得到的有效弱口令的总数
func ResultTotal() {
	config.ProcessBarScan.Finish()
	log.Log.Info(fmt.Sprintf("Finshed scan, total result: %v, used time: %v",
		config.CacheTarget.ItemCount(),
		time.Since(config.StartTime)))
}

// 将爆破结果保存到一个 DB 文件中，DB文件的格式为 go-cache 库定义的格式
func SaveResultToFile() error {
	return config.CacheTarget.SaveFile("crack_result.db")
}

// 将结果导出到一个 config 指定的 txt 文件中
func DumpToFile(filename string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	// 遍历爆破结果，将每一个爆破结果写入 file
	_, items := CacheStatus()
	for _, v := range items {
		result := v.Object.(TScanResult)
		_, _ = file.WriteString(fmt.Sprintf("%v:%v|%v,%v:%v\n",
			result.TargetsModel.Ip,
			result.TargetsModel.Port,
			result.TargetsModel.Protocol,
			result.TargetsModel.Username,
			result.TargetsModel.Password),
		)
	}

	return err
}

// 返回爆破结果的数目和详细内容
func CacheStatus() (count int, items map[string]cache.Item) {
	count = config.CacheTarget.ItemCount()
	items = config.CacheTarget.Items()
	return count, items
}
