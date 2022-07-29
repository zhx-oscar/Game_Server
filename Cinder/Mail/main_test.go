package main

/* 测试整个服务器。
	go test -c -covermode=count -coverpkg ./...
生成 Mail.test.exe. 复制到 bin 目录下运行：
	Mail.test.exe --systemTest --test.coverprofile cover.out
生成代码覆盖测试结果 cover.out, 需要在 go.mod 管理的目录下执行
	go tool cover -html=cover.out -o cover.html
打开 cover.html 查看结果。
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
