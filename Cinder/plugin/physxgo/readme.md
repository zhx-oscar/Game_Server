# Linux编译教程
## 一、安装docker
https://g-cinder.github.io/Memoranda/project_specs/local_server_guide.html

## 二、下载PhysX-4.1
https://github.com/NVIDIAGameWorks/PhysX

## 三、在CentOS7上编译Physx（项目git上已经有编译好的PhysX，一般情况无需执行此步骤）
### 1. 安装physx依赖项
```
clang++ 4.0.1
cmake 3.17.3
python 2.7.5
```
安装的依赖项必须大于等于以上版本

### 2. 修改physx编译配置
```
cd /code/PhysX-4.1/physx/buildtools/presets/public
vi linux.xml
```
修改PX_BUILDSNIPPETS、PX_BUILDPUBLICSAMPLES为False，PX_GENERATE_STATIC_LIBRARIES为True（注意大小写）

### 3. 生成physx工程文件
```
cd /code/PhysX-4.1/physx
./generate_projects.sh
Preset parameter required, available presets:
(0) ios64 <--- iOS Xcode PhysX general settings
(1) linux-aarch64 <--- Linux-aarch64 gcc PhysX SDK general settings
(2) linux <--- Linux clang PhysX SDK general settings
(3) mac64 <--- macOS Xcode PhysX general settings
Enter preset number: 2
...
```
注意检查生成结果，确认所有的工程文件生成成功，可能会发生以下几种CMake错误：
- 变量空字符串问题，加上引号括起报错的变量即可
- CACHE INTERNAL写错为CACHE INTERAL，修改即可
- 输出文件生成在UNKNOWN目录下，未生成在linux.clang目录下，这种情况不用特殊处理
- 其他发生错误，将C:\code\PhysX-4.1\physx\compiler目录下  
linux-checked、linux-debug、linux-profile、linux-release删除再重试即可

### 4. 编译physx工程
release版本：
```
cd /code/PhysX-4.1/physx/compiler/linux-release
make clean
make
make install
...
```
debug版本：
```
cd /code/PhysX-4.1/physx/compiler/linux-debug
make clean
make
make install
...
```

### 5. 拷贝physx install目录至项目目录
在windows中将以下目录
```
C:\code\PhysX-4.1\physx\install
```
拷贝至
```
C:\code\newcldzz\Daisy\Server\src\physxgo\PhysX-4.1
```
注意如果库文件生成在UNKNOWN目录下，需要重命名为linux.clang，physx编译完成

## 四、在ubuntu上编译PhysX（目前项目在CentOS7上部署，ununtu编译仅做参考）
假设PhysX-4.1存放在C:\code\PhysX-4.1，项目工程存放在C:\code\newcldzz\Daisy目录下，执行以下步骤
### 1. 安装ubuntu镜像
```
docker run -itd -v /c/code:/code ubuntu /bin/bash
docker exec -it 容器ID /bin/bash
```

### 2. 安装physx依赖项
```
apt-get update
apt-get install libterm-readkey-perl python2.7 cmake clang-9 make
...
update-alternatives --install /usr/bin/python python /usr/bin/python2.7 2
update-alternatives --install /usr/bin/clang clang /usr/bin/clang-9 9 --slave /usr/bin/clang++ clang++ /usr/bin/clang++-9
```

### 3. 修改physx编译配置
```
cd /code/PhysX-4.1/physx/buildtools/presets/public
vi linux.xml
```
修改PX_BUILDSNIPPETS、PX_BUILDPUBLICSAMPLES为False，PX_GENERATE_STATIC_LIBRARIES为True（注意大小写）

### 4. 生成physx工程文件
```
cd /code/PhysX-4.1/physx
./generate_projects.sh
Preset parameter required, available presets:
(0) ios64 <--- iOS Xcode PhysX general settings
(1) linux-aarch64 <--- Linux-aarch64 gcc PhysX SDK general settings
(2) linux <--- Linux clang PhysX SDK general settings
(3) mac64 <--- macOS Xcode PhysX general settings
Enter preset number: 2
...
```
注意检查生成结果，确认所有的工程文件生成成功，可能会发生以下几种CMake错误：
- 变量空字符串问题，加上引号括起报错的变量即可
- CACHE INTERNAL写错为CACHE INTERAL，修改即可
- 输出文件生成在UNKNOWN目录下，未生成在linux.clang目录下，这种情况不用特殊处理
- 其他发生错误，将C:\code\PhysX-4.1\physx\compiler目录下  
linux-checked、linux-debug、linux-profile、linux-release删除再重试即可

### 5. 编译physx工程
release版本：
```
cd /code/PhysX-4.1/physx/compiler/linux-release
make clean
make
make install
...
```
debug版本：
```
cd /code/PhysX-4.1/physx/compiler/linux-debug
make clean
make
make install
...
```

### 6. 拷贝physx install目录至项目目录
在windows中将以下目录
```
C:\code\PhysX-4.1\physx\install
```
拷贝至
```
C:\code\newcldzz\Daisy\Server\src\physxgo\PhysX-4.1
```
注意如果库文件生成在UNKNOWN目录下，需要重命名为linux.clang，physx编译完成

## 五、编译physx c语言包装lib
### 1. 编译lib
release版本：
```
cd /c/code/newcldzz/Daisy/Server/src/physxgo/physxcwrap
make clean
make release
```
debug版本：
```
cd /c/code/newcldzz/Daisy/Server/src/physxgo/physxcwrap
make clean
make debug
```

### 2. 修改编译连接版本
注意go代码中默认编译连接release版本，使用debug版本调试时需要在physx_sdk.go修改代码

----

# Windows编译教程
## 一、安装visual studio 2017以上版本和msys2
https://visualstudio.microsoft.com/zh-hans/vs/
https://www.msys2.org/

注意msys安装后，在开始菜单中找到MSYS2 MinGW 64-bit并运行，并执行以下命令安装工具，中途可能需要重启命令行界面并重新执行命令
```
pacman -Syu
pacman -S msys2-devel
pacman -S mingw-w64-x86_64-toolchain
```

## 二、下载PhysX-4.1
https://github.com/NVIDIAGameWorks/PhysX

## 三、编译PhysX（项目git上已经有编译好的PhysX，一般情况无需执行此步骤）
假设PhysX-4.1存放在C:\code\PhysX-4.1，项目工程存放在C:\code\newcldzz\Daisy目录下，执行以下步骤
### 1. 生成physx工程文件
在windows中进入以下目录
```
cd /code/PhysX-4.1/physx
```
运行以下脚本
```
generate_projects.bat
(0) android-arm64-v8a <--- Android-19, arm64-v8a PhysX SDK
(1) android <--- Android-19, armeabi-v7a with NEON PhysX SDK
(2) vc12win32 <--- VC12 Win32 PhysX general settings
(3) vc12win64 <--- VC12 Win64 PhysX general settings
(4) vc14win32 <--- VC14 Win32 PhysX general settings
(5) vc14win64 <--- VC14 Win64 PhysX general settings
(6) vc15uwp32 <--- VC15 UWP 32bit PhysX general settings
(7) vc15uwp64 <--- VC15 UWP 64bit PhysX general settings
(8) vc15uwparm32 <--- VC15 UWP 32bit PhysX general settings
(9) vc15uwparm64 <--- VC15 UWP ARM 64bit PhysX general settings
(10) vc15win32 <--- VC15 Win32 PhysX general settings
(11) vc15win64 <--- VC15 Win64 PhysX general settings
(12) vc16win32 <--- VC16 Win32 PhysX general settings
(13) vc16win64 <--- VC16 Win64 PhysX general settings
Enter preset number: 13
...
```
注意检查生成结果，确认所有的工程文件生成成功，如果发生错误，将C:\code\PhysX-4.1\physx\compiler目录下vc16win64删除再重试即可

### 2. 修改physx编译配置
在windows中进入以下目录
```
/code/PhysX-4.1/physx/buildtools/presets/public
```
打开配置文件
```
vc16win64.xml
```
修改PX_BUILDSNIPPETS、PX_BUILDPUBLICSAMPLES、PX_GENERATE_STATIC_LIBRARIES为False（注意大小写）

### 3. 修改IDE平台工具集 
使用高于visual studio 2017以上版本IDE时，在打开工程（sln）后，选中所有项目(project)，右击打开属性页，将平台工具集修改为Visual Studio 2017(v141)

### 4. 编译
先编译ALL_BUILD项目，再编译INSTALL项目

### 5. 拷贝physx install目录至项目目录
在windows中将以下目录
```
C:\code\PhysX-4.1\physx\install
```
拷贝至
```
C:\code\newcldzz\Daisy\Server\src\physxgo\PhysX-4.1
```
physx编译完成

## 四、编译physx c语言包装lib
### 1. 编译lib
在windows中进入以下目录 
```
C:\code\newcldzz\Daisy\Server\src\physxgo\physxcwrap
```
打开工程（sln）开始编译

### 2. 执行安装脚本
release版本：
```
cd /c/code/newcldzz/Daisy/Server/src/physxgo
./install_release.bat
```
debug版本：
```
cd /c/code/newcldzz/Daisy/Server/src/physxgo
./install_debug.bat
```
注意go代码中默认编译连接release版本，使用debug版本时需要在physx_sdk.go修改代码