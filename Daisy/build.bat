@echo off
set GOPROXY=https://goproxy.cn

set servers=..\Cinder\Chat;..\Cinder\Login;..\Cinder\Agent;..\Daisy\DBAgent;..\Daisy\Battle;..\Daisy\Game;..\Cinder\Mail

for %%I in (%servers%) do (
	echo build %%I
	@echo on
	cd %%I
	go build -o %~dp0/../../bin
	cd /d %~dp0
	@echo off
)

echo Build Finish!
pause