@echo off
chcp 65001 >nul
echo ========================================
echo   PBFT DID Chain - 启动所有节点
echo ========================================
echo.

REM 检查可执行文件是否存在
if not exist "pbftdidchain.exe" (
    echo [错误] 找不到 pbftdidchain.exe 文件！
    echo 请先编译项目: go build -o pbftdidchain.exe
    pause
    exit /b 1
)

echo [信息] 正在启动4个PBFT节点...
echo.

REM 启动节点0 (主节点)
echo [启动] 节点0 (主节点, 端口30000)...
start "PBFT-Node-0" cmd /k "title PBFT Node 0 && pbftdidchain.exe 0"
timeout /t 2 /nobreak >nul

REM 启动节点1
echo [启动] 节点1 (端口30001)...
start "PBFT-Node-1" cmd /k "title PBFT Node 1 && pbftdidchain.exe 1"
timeout /t 2 /nobreak >nul

REM 启动节点2
echo [启动] 节点2 (端口30002)...
start "PBFT-Node-2" cmd /k "title PBFT Node 2 && pbftdidchain.exe 2"
timeout /t 2 /nobreak >nul

REM 启动节点3
echo [启动] 节点3 (端口30003)...
start "PBFT-Node-3" cmd /k "title PBFT Node 3 && pbftdidchain.exe 3"
timeout /t 2 /nobreak >nul

echo.
echo ========================================
echo [完成] 所有节点已启动！
echo.
echo 提示：
echo   - 每个节点都在独立的窗口中运行
echo   - 等待所有节点连接成功后再运行客户端
echo   - 使用 stop_all.bat 可以停止所有节点
echo ========================================
echo.
pause

