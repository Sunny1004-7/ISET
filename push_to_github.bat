@echo off
chcp 65001 >nul
echo ========================================
echo 正在推送代码到 GitHub...
echo ========================================
echo.

cd /d "%~dp0"

echo [1/4] 检查 Git 状态...
git status
echo.

echo [2/4] 添加所有更改...
git add .
echo.

echo [3/4] 创建提交...
git commit -m "Update: Sync ISET project to GitHub" 2>nul
if errorlevel 1 (
    echo 没有新的更改需要提交
)
echo.

echo [4/4] 推送到 GitHub...
echo 请在弹出的浏览器窗口中登录 GitHub 账号
echo.

git push -u origin main

if errorlevel 1 (
    echo.
    echo ========================================
    echo ❌ 推送失败！
    echo ========================================
    echo.
    echo 可能的原因：
    echo 1. 网络连接问题
    echo 2. 需要配置 GitHub 身份验证
    echo 3. Git 安装不完整
    echo.
    echo 请尝试以下解决方案：
    echo.
    echo 方案A - 使用 GitHub Desktop（推荐）：
    echo 1. 下载：https://desktop.github.com/
    echo 2. 登录账号后选择"Add Local Repository"
    echo 3. 选择此项目文件夹并点击"Publish"
    echo.
    echo 方案B - 使用 Git Bash：
    echo 1. 在项目文件夹右键选择"Git Bash Here"
    echo 2. 执行: git push -u origin main
    echo.
    echo 方案C - 配置 SSH：
    echo 运行 setup_ssh.bat 脚本
    echo.
) else (
    echo.
    echo ========================================
    echo ✅ 成功推送到 GitHub！
    echo ========================================
    echo.
    echo 访问您的仓库：
    echo https://github.com/Sunny1004-7/SUNNY-AITEACH
    echo.
)

pause

