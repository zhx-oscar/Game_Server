pushd ..
go test -c -covermode=count -coverpkg ./...
go test ./rpcproc -c -o Mail.rpcproc.test.exe -covermode=count -coverpkg ./...

ping -n 1 127.1 >nul
copy *.test.exe ..\..\..\bin\
del *.test.exe
popd

REM cd ..\..\..\..\bin\
REM Mail.test.exe --systemTest --test.coverprofile ..\src\Cinder\Mail\coverage_test\cover.out
pause
