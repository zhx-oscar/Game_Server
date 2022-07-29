@echo off
if not exist ..\..\..\..\bin md ..\..\..\..\bin
copy /b /y .\physxcwrap\install\lib\debug\PhysXCWrap.dll ..\..\..\..\bin\PhysXCWrap.dll
copy /b /y .\PhysX-4.1\install\vc15win64\PhysX\bin\win.x86_64.vc142.mt\debug\PhysX_64.dll ..\..\..\..\bin\PhysX_64.dll
copy /b /y .\PhysX-4.1\install\vc15win64\PhysX\bin\win.x86_64.vc142.mt\debug\PhysXCooking_64.dll ..\..\..\..\bin\PhysXCooking_64.dll
copy /b /y .\PhysX-4.1\install\vc15win64\PhysX\bin\win.x86_64.vc142.mt\debug\PhysXCommon_64.dll ..\..\..\..\bin\PhysXCommon_64.dll
copy /b /y .\PhysX-4.1\install\vc15win64\PhysX\bin\win.x86_64.vc142.mt\debug\PhysXFoundation_64.dll ..\..\..\..\bin\PhysXFoundation_64.dll

pause