@echo off
if not exist ..\..\..\..\bin md ..\..\..\..\bin
copy /b /y .\physxcwrap\install\lib\release\PhysXCWrap.dll ..\..\..\..\bin\PhysXCWrap.dll
copy /b /y .\PhysX-4.1\install\vc15win64\PhysX\bin\win.x86_64.vc142.mt\release\PhysX_64.dll ..\..\..\..\bin\PhysX_64.dll
copy /b /y .\PhysX-4.1\install\vc15win64\PhysX\bin\win.x86_64.vc142.mt\release\PhysXCooking_64.dll ..\..\..\..\bin\PhysXCooking_64.dll
copy /b /y .\PhysX-4.1\install\vc15win64\PhysX\bin\win.x86_64.vc142.mt\release\PhysXCommon_64.dll ..\..\..\..\bin\PhysXCommon_64.dll
copy /b /y .\PhysX-4.1\install\vc15win64\PhysX\bin\win.x86_64.vc142.mt\release\PhysXFoundation_64.dll ..\..\..\..\bin\PhysXFoundation_64.dll

pause