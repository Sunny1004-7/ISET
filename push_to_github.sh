#!/bin/bash
# 推送代码到GitHub的脚本
# 请在Git Bash中运行此脚本

echo "========================================"
echo "正在推送代码到GitHub..."
echo "========================================"

# 检查是否有未提交的更改
if [[ -n $(git status -s) ]]; then
    echo "检测到未提交的更改，正在添加..."
    git add .
    git commit -m "Update: sync latest changes"
fi

# 推送到GitHub
echo "推送到 GitHub..."
git push -u origin main

if [ $? -eq 0 ]; then
    echo "========================================"
    echo "✅ 成功推送到 GitHub!"
    echo "访问: https://github.com/Sunny1004-7/SUNNY-AITEACH"
    echo "========================================"
else
    echo "========================================"
    echo "❌ 推送失败，请检查："
    echo "1. 网络连接是否正常"
    echo "2. GitHub凭据是否正确"
    echo "3. 仓库权限是否足够"
    echo "========================================"
fi

