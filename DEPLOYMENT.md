# 部署指南

## ✅ 已完成的工作

1. ✅ 创建了专业的 `README.md` 文件（包含完整的项目说明）
2. ✅ 删除了原有的 `readme.txt` 文件
3. ✅ 创建了 `.gitignore` 文件（排除不必要的文件）
4. ✅ 初始化了 Git 仓库
5. ✅ 创建了初始提交
6. ✅ 配置了远程仓库：`https://github.com/Sunny1004-7/SUNNY-AITEACH.git`

## 🔄 需要完成的工作

由于 Git HTTPS helper 配置问题，需要手动推送代码到 GitHub。

### 方法一：使用 Git Bash（推荐）

1. 打开 **Git Bash**（不是 PowerShell）
2. 导航到项目目录：
   ```bash
   cd "D:\DONTSTOP\XJTU教改项目\机器学习教改-AI出题模块\ISET"
   ```
3. 运行推送脚本：
   ```bash
   bash push_to_github.sh
   ```

### 方法二：直接使用 Git Bash 命令

1. 打开 **Git Bash**
2. 导航到项目目录：
   ```bash
   cd "D:\DONTSTOP\XJTU教改项目\机器学习教改-AI出题模块\ISET"
   ```
3. 执行推送命令：
   ```bash
   git push -u origin main
   ```

### 方法三：使用 GitHub Desktop

1. 下载并安装 [GitHub Desktop](https://desktop.github.com/)
2. 在 GitHub Desktop 中选择 "Add an existing repository"
3. 选择项目目录
4. 点击 "Publish repository" 或 "Push origin"

### 方法四：修复 Git 配置

如果以上方法都不行，可能需要重新安装或修复 Git：

1. 下载最新版 Git：https://git-scm.com/download/win
2. 重新安装，确保选择 "Git from the command line and also from 3rd-party software"
3. 安装完成后重启命令行，再次尝试推送

## 📝 首次推送说明

如果这是第一次推送到这个仓库，可能需要：

1. **设置 Git 用户信息**（如果尚未设置）：
   ```bash
   git config --global user.name "Your Name"
   git config --global user.email "your.email@example.com"
   ```

2. **GitHub 身份验证**：
   - 方式1：使用 Personal Access Token（推荐）
     - 前往 GitHub Settings → Developer settings → Personal access tokens
     - 生成新的 token（权限选择 repo）
     - 推送时使用 token 作为密码
   
   - 方式2：使用 SSH Key
     ```bash
     # 生成 SSH key
     ssh-keygen -t rsa -b 4096 -C "your.email@example.com"
     
     # 将 SSH key 添加到 GitHub
     # 复制 ~/.ssh/id_rsa.pub 内容到 GitHub Settings → SSH keys
     
     # 更改远程仓库为 SSH URL
     git remote set-url origin git@github.com:Sunny1004-7/SUNNY-AITEACH.git
     git push -u origin main
     ```

## 🔍 验证推送成功

推送成功后，访问以下地址查看：
- 仓库主页：https://github.com/Sunny1004-7/SUNNY-AITEACH
- README 预览：https://github.com/Sunny1004-7/SUNNY-AITEACH/blob/main/README.md

## 📊 推送内容

本次推送包含：
- ✅ 217 个文件
- ✅ 25,597 行代码
- ✅ 完整的项目结构
- ✅ 前端、后端、AI 模块的所有源代码
- ✅ 配置文件和文档

## 🚀 后续维护

### 日常更新代码

```bash
# 1. 添加更改
git add .

# 2. 提交更改
git commit -m "描述你的更改"

# 3. 推送到 GitHub
git push origin main
```

### 拉取最新代码

```bash
git pull origin main
```

### 创建新分支

```bash
# 创建并切换到新分支
git checkout -b feature/new-feature

# 推送新分支到 GitHub
git push -u origin feature/new-feature
```

## ❓ 常见问题

### Q1: 推送时提示 "Authentication failed"
**A**: 需要配置 GitHub 身份验证，参考上面的 "GitHub 身份验证" 部分。

### Q2: 推送时提示 "remote: Repository not found"
**A**: 检查仓库 URL 是否正确，或者确认你有该仓库的访问权限。

### Q3: 推送时提示 "Updates were rejected"
**A**: 远程仓库有更新，需要先拉取：
```bash
git pull origin main --rebase
git push origin main
```

### Q4: 想撤销某次提交
**A**: 
```bash
# 撤销最后一次提交（保留更改）
git reset --soft HEAD~1

# 撤销最后一次提交（丢弃更改）
git reset --hard HEAD~1
```

## 📞 需要帮助？

如果遇到问题，可以：
1. 查看 Git 文档：https://git-scm.com/doc
2. 查看 GitHub 文档：https://docs.github.com/
3. 在项目 Issues 中提问

---

**祝推送顺利！** 🎉

