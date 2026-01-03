@echo off
chcp 65001 >nul
echo ========================================
echo   PBFT DID Chain - 启动客户端
echo ========================================

REM 检查客户端可执行文件是否存在
if not exist "test\client.exe" (
    echo [警告] 找不到 test\client.exe，尝试编译...
    cd test
    go build -o client.exe client.go
    if errorlevel 1 (
        echo [错误] 编译失败！
        pause
        exit /b 1
    )
    cd ..
)

REM 获取主节点ID参数（默认为0）
set PRIMARY_ID=0
if not "%1"=="" (
    set PRIMARY_ID=%1
)

echo [信息] 正在启动客户端，连接到主节点 %PRIMARY_ID%...
echo.

cd test
start "PBFT-Client" cmd /k "title PBFT Client && client.exe %PRIMARY_ID%"
cd ..

echo.
echo [完成] 客户端已启动！
echo.
echo 提示：
echo   - 客户端将向主节点发送请求
echo   - 等待收到至少2个节点的回复后认为共识成功
echo ========================================
echo.
pause

