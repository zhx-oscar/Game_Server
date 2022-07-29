pushd ..
go test ./... -covermode=count -coverpkg=./... -coverprofile coverage_test/cover.out
popd

pause
