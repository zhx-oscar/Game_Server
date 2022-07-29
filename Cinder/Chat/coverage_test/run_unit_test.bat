pushd ..
go test ./... -covermode=count -coverpkg=./... -coverprofile coverage_test/Chat.cov
popd

pause
