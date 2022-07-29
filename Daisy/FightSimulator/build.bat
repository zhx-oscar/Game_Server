@echo off
set GOPROXY=https://mirrors.aliyun.com/goproxy/

for /F %%i in ('git rev-parse HEAD') do ( set commitid=%%i)

go build -o ..\..\..\bin\FightSimulator.exe -ldflags="-X main.CommitID=%commitid%"
if %errorlevel% neq 0 (
    echo Build Failed!
    pause
    exit /b
)

copy /b /y ..\..\..\bin\FightSimulator.exe ..\..\..\..\Client\Assets\Tools\FightSimulator\Editor\FightSimulator.exe

echo Build Finish!
pause