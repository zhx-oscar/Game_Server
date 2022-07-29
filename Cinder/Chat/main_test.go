package main

/* 测试整个服务器。
	go test -c -covermode=count -coverpkg ./...
生成 Chat.test.exe. 复制到 bin 目录下运行：
	Chat.test.exe --systemTest --test.coverprofile Chat.cov
生成代码覆盖测试结果 Chat.cov, 需要在 go.mod 管理的目录下执行
	go tool cover -html=Chat.cov -o Chat.html
打开 Chat.html 查看结果。
*/

import (
	"flag"
	"fmt"
	"testing"
)

var systemTest *bool

func init() {
	systemTest = flag.Bool("systemTest", false, "Set to true when running system tests")
}

// Test started when the test binary is started. Only calls main.
func TestSystem(t *testing.T) {
	if *systemTest {
		fmt.Println("Test system...")
		main()
	}
}
