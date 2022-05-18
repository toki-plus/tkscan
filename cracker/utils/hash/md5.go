package hash

import (
	"crypto/md5"
	"fmt"
	"io"
)

func MD5(s string) (m string) {
	h := md5.New()
	// 将字符串 s 写入 io.Writer
	_, _ = io.WriteString(h, s)
	// Sum 将当前的 hash 添加到 []byte，并将 []byte 返回
	return fmt.Sprintf("%x", h.Sum(nil))
}
