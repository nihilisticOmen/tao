chcp 65001
@echo off
:loop
@echo off&amp;color 0A
cls
echo,
echo Please select the system environment you want to compile：
echo,
echo 1. Windows_amd64
echo 2. linux_amd64

set/p action=pleaseSelect:
if %action% == 1 goto build_Windows_amd64
if %action% == 2 goto build_linux_amd64

:build_Windows_amd64
echo COMPILE_WINDOWS_VERSION_64_BIT
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o project-user/target/project-user.exe project-user/main.go
go build -o project-api/target/project-api.exe project-api/main.go
:build_linux_amd64
echo COMPILE_LINUX_VERSION_64_BIT
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o project-user/target/project-user project-user/main.go
go build -o project-api/target/project-api project-api/main.go