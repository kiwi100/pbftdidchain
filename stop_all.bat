@echo off
chcp 65001 >nul
echo ========================================
echo   PBFT DID Chain - 停止所有节点
echo ========================================
echo.

echo [信息] 正在查找并停止所有PBFT进程...

REM 停止pbftdidchain.exe进程
tasklist /FI "IMAGENAME eq pbftdidchain.exe" 2>NUL | find /I /N "pbftdidchain.exe">NUL
if "%ERRORLEVEL%"=="0" (
    echo [停止] 正在终止 pbftdidchain.exe 进程...
    taskkill /F /IM pbftdidchain.exe >nul 2>&1
    if errorlevel 1 (
        echo [警告] 部分进程可能无法终止，请手动检查
    ) else (
        echo [成功] pbftdidchain.exe 进程已停止
    )
) else (
    echo [信息] 未找到运行中的 pbftdidchain.exe 进程
)

REM 停止client.exe进程
tasklist /FI "IMAGENAME eq client.exe" 2>NUL | find /I /N "client.exe">NUL
if "%ERRORLEVEL%"=="0" (
    echo [停止] 正在终止 client.exe 进程...
    taskkill /F /IM client.exe >nul 2>&1
    if errorlevel 1 (
        echo [警告] 部分进程可能无法终止，请手动检查
    ) else (
        echo [成功] client.exe 进程已停止
    )
) else (
    echo [信息] 未找到运行中的 client.exe 进程
)

REM 关闭所有以"PBFT-"开头的窗口
echo.
echo [信息] 正在关闭PBFT相关窗口...
for /f "tokens=2" %%a in ('tasklist /FI "WINDOWTITLE eq PBFT-*" /FO LIST 2^>NUL ^| findstr /I "PID"') do (
    taskkill /F /PID %%a >nul 2>&1
)

echo.
echo ========================================
echo [完成] 清理完成！
echo ========================================
echo.
pause

