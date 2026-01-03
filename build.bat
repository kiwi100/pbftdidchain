@echo off
chcp 65001 >nul
echo ========================================
echo   PBFT DID Chain - 编译项目
echo ========================================
echo.

echo [编译] 正在编译主程序...
go build -o pbftdidchain.exe
if errorlevel 1 (
    echo [错误] 主程序编译失败！
    pause
    exit /b 1
)
echo [成功] pbftdidchain.exe 编译完成
echo.

echo [编译] 正在编译客户端...
cd test
go build -o client.exe client.go
if errorlevel 1 (
    echo [错误] 客户端编译失败！
    cd ..
    pause
    exit /b 1
)
cd ..
echo [成功] client.exe 编译完成
echo.

echo ========================================
echo [完成] 所有程序编译完成！
echo ========================================
echo.
pause

