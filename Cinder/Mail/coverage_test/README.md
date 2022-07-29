# 代码覆盖测试

测试时需要 etcd, mongodb, nsq 正常运行。

## 单元测试

运行 run_unity_test.bat，运行所有单元测试，结束后在当前目录下生成 cover.out,
然后运行 generate_coverage.bat 生成 cover.html.

## 系统测试

用 build_test.bat 编译 Mail.test.exe 到 bin 目录，
用 run_system_test.bat 在bin目录运行 Mail.text.exe，
其行为是正常的Mail服。
Ctrl-C结束并在当前目录下生成 cover.out。

## 单元测试说明

run_unit_test.bat 运行全部子包测试，但无法编译成一个可执行文件。

`build_test.bat`中`-c`参数能将 main 包的测试编译成执行文件，但不能包含子包的测试。
如果将编译命令改为
```
go test ./... -c -covermode=count -coverpkg ./...
```
则报错：
```
cannot use -c flag with multiple packages
```

可以编译每个子包并运行，这样需要手工合并cover.out
```
go test rpcproc -o Mail.rpcproc.test.exe -c -covermode=count -coverpkg ./...
```

可以在Mail目录下直接使用 go test:
```
go test ./... -covermode=count -coverpkg=./... -coverprofile coverage_test/cover.out
```