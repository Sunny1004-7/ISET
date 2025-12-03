@echo off
chcp 65001 >nul
echo ========================================
echo GitHub SSH 配置助手
echo ========================================
echo.

echo 此脚本将帮助您配置 SSH 密钥以推送到 GitHub
echo.
echo 步骤：
echo 1. 生成 SSH 密钥
echo 2. 将公钥添加到 GitHub
echo 3. 切换远程 URL 为 SSH
echo 4. 测试连接并推送
echo.
pause

echo.
echo [1/4] 检查现有 SSH 密钥...
if exist "%USERPROFILE%\.ssh\id_rsa.pub" (
    echo ✓ 找到现有的 SSH 密钥
    type "%USERPROFILE%\.ssh\id_rsa.pub"
) else if exist "%USERPROFILE%\.ssh\id_ed25519.pub" (
    echo ✓ 找到现有的 SSH 密钥
    type "%USERPROFILE%\.ssh\id_ed25519.pub"
) else (
    echo × 未找到 SSH 密钥，需要生成新密钥
    echo.
    echo 请在 Git Bash 中运行以下命令：
    echo ssh-keygen -t ed25519 -C "your_email@example.com"
    echo.
    echo 然后重新运行此脚本
    pause
    exit /b
)

echo.
echo [2/4] 请将上面显示的公钥添加到 GitHub：
echo.
echo 1. 访问：https://github.com/settings/keys
echo 2. 点击 "New SSH key"
echo 3. 粘贴上面显示的公钥内容
echo 4. 点击 "Add SSH key"
echo.
pause

echo.
echo [3/4] 切换远程 URL 为 SSH...
cd /d "%~dp0"
git remote set-url origin git@github.com:Sunny1004-7/SUNNY-AITEACH.git
echo ✓ 远程 URL 已更新
git remote -v
echo.

echo [4/4] 测试 SSH 连接...
ssh -T git@github.com 2>&1 | findstr /C:"successfully authenticated"
if errorlevel 1 (
    echo.
    echo × SSH 连接失败，请检查：
    echo 1. SSH 密钥是否已添加到 GitHub
    echo 2. SSH 服务是否正常运行
    echo.
) else (
    echo ✓ SSH 连接成功！
    echo.
    echo 现在可以推送代码了...
    git push -u origin main
)

echo.
pause

