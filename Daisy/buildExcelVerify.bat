@echo off
set GOPROXY=https://mirrors.aliyun.com/goproxy/

cd /d ..\Daisy\Data\app
go build -o ..\..\..\..\..\Tools\ExcelVerify\excelverify.exe
if %errorlevel% neq 0 (
    echo Build Failed!
    pause
    exit /b
)

echo Build Finish!
pause