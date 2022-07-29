pushd ..
go test -c -covermode=count -coverpkg ./...
go test ./rpcproc -c -o Chat.rpcprocc.test.exe -covermode=count -coverpkg ./...

ping -n 1 127.1 >nul
copy *.test.exe ..\..\..\bin\
del *.test.exe
popd

REM cd ..\..\..\..\bin\
REM Chat.test.exe --systemTest --test.coverprofile ..\src\Cinder\Chat\coverage_test\Chat.cov
pause
