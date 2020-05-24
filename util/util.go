package util

import (
	"math/rand"
	"time"
)

func RandomString(n int) string {
	var letters = []byte("abcdefghtyuikllwetyibvhjigkgggkgju")
	result := make([]byte, n)
	// 设置种子，不然每次都会随机成0
	rand.Seed(time.Now().Unix())

	for i := range result {
		// 生成 letters的长度(n) 到 (n-1) 的随机数（输出 int类型）
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}
