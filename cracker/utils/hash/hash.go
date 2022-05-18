package hash

import "github.com/toki-plus/tkscan/cracker/config"

// 将 task 转成哈希
func MakeTaskHash(k string) string {
	hash := MD5(k)
	return hash
}

// 检查特定的任务是否爆破成功，成功的话不再尝试破解该 target
func CheckTaskHash(hash string) bool {
	// SuccessHash 是一个 sync.Map
	// Load 加载一个存储在 map 里值作为 key，ok 代表着这个 key 存在于 map
	_, ok := config.SuccessHash.Load(hash)
	return ok
}

// 存储 task 的 hash 为 key，true 为值
func SetTaskHash(hash string) {
	config.SuccessHash.Store(hash, true)
}
