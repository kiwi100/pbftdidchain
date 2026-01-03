@echo off
chcp 65001 >nul
echo ========================================
echo   PBFT DID Chain - 一键启动
echo ========================================
echo.

REM 检查可执行文件
if not exist "pbftdidchain.exe" (
    echo [错误] 找不到 pbftdidchain.exe 文件！
    echo 请先编译项目: go build -o pbftdidchain.exe
    pause
    exit /b 1
)

echo [步骤1/2] 启动所有节点...
call start_nodes.bat

echo.
echo [步骤2/2] 等待节点连接（5秒）...
timeout /t 5 /nobreak >nul

echo.
echo [信息] 是否现在启动客户端？(Y/N)
set /p choice="请输入: "
if /i "%choice%"=="Y" (
    call start_client.bat
) else (
    echo [信息] 稍后可以使用 start_client.bat 启动客户端
)

echo.
echo ========================================
echo [完成] 启动流程结束！
echo ========================================
echo.

