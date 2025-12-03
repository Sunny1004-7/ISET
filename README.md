# ISET - Intelligent Smart Education Technology

<div align="center">

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Python](https://img.shields.io/badge/python-3.7-blue.svg)
![Go](https://img.shields.io/badge/go-1.22.5-00ADD8.svg)
![Vue](https://img.shields.io/badge/vue-3.2.45-4FC08D.svg)
![TensorFlow](https://img.shields.io/badge/TensorFlow-1.14.0-FF6F00.svg)

**一个基于AI多智能体的智慧教育系统**

[功能特性](#功能特性) • [系统架构](#系统架构) • [快速开始](#快速开始) • [技术栈](#技术栈) • [文档](#文档)

</div>

---

## 📖 项目简介

ISET（Intelligent Smart Education Technology）是面向西安交通大学教改项目开发的智慧教育系统，专注于AI出题模块。系统整合了**知识图谱**、**大语言模型**和**多智能体技术**，为教师和学生提供个性化的教学与学习支持。

### 核心价值

- 🎯 **个性化习题生成**：基于学生知识状态，智能生成针对性练习题
- 📊 **智能课程推荐**：分析学习历史，推荐最适合的课程路径
- 💬 **苏格拉底式问答**：情绪感知的智能教学对话系统
- 🗺️ **知识图谱可视化**：直观展示课程知识结构与关联

---

## ✨ 功能特性

### 学生端
- ✅ **课程学习管理**：课程列表、成绩查询、学情统计
- ✅ **知识图谱展示**：大纲视图、思维导图、知识图谱
- ✅ **智能习题生成**：基于个人知识状态的个性化习题
- ✅ **AI学习助手**：智能问答、学习建议、情绪感知对话

### 教师端
- ✅ **课程管理**：学生信息、成绩统计、学情分析
- ✅ **知识图谱编辑**：课程知识结构构建与维护
- ✅ **智能出题**：支持指定知识点、难度、题型的习题生成
- ✅ **教学辅助**：备课支持、教学策略推荐

---

## 🏗️ 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                        前端应用层                              │
│   ┌──────────────────────┐  ┌──────────────────────┐        │
│   │   学生学习系统         │  │   教师管理系统         │        │
│   │  (Vue3 + TypeScript) │  │  (Vue3 + TypeScript) │        │
│   └──────────────────────┘  └──────────────────────┘        │
└─────────────────────────────────────────────────────────────┘
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      业务逻辑层 (Go + Gin)                     │
│   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│   │  Controller │→ │   Service   │→ │     DAO     │        │
│   │  请求处理    │  │  业务逻辑     │  │  数据访问    │        │
│   └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      AI算法层 (Python)                         │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐        │
│  │ 习题生成模块  │ │ 课程推荐模块  │ │ 问答系统模块  │        │
│  │ Multi-Agent │ │ Multi-Agent │ │   AutoGen    │        │
│  └──────────────┘ └──────────────┘ └──────────────┘        │
└─────────────────────────────────────────────────────────────┘
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                        数据层                                 │
│  SQLite  │  Neo4j  │  Redis  │  File Storage  │  MOOCCubeX │
└─────────────────────────────────────────────────────────────┘
```

---

## 🛠️ 技术栈

### 前端
- **框架**: Vue 3.2.45 + TypeScript
- **构建工具**: Vite 4.0.2
- **UI组件**: Element Plus 2.10.6
- **状态管理**: Pinia 2.0.28
- **路由**: Vue Router 4.1.6
- **图表**: ECharts 5.4.1
- **知识图谱**: jsmind 0.7.1, vue-konva 3.0.2

### 后端
- **语言**: Go 1.22.5
- **框架**: Gin (HTTP Web框架)
- **ORM**: GORM
- **数据库**: SQLite
- **实时通信**: WebSocket (gorilla/websocket)

### AI算法模块
- **语言**: Python 3.7
- **深度学习**: TensorFlow 1.14.0
- **多智能体框架**: AutoGen 0.7.0
- **大语言模型**: OpenAI 1.99.6
- **数据处理**: Pandas, NumPy
- **GPU加速**: CUDA 10.2, cuDNN 7.6.5

### 数据库
- **关系型**: SQLite (课程、用户、习题数据)
- **图数据库**: Neo4j (知识图谱)
- **缓存**: Redis (数据缓存)

---

## 🚀 快速开始

### 环境要求

#### 硬件要求
- **GPU**: NVIDIA RTX 3080 10GB (或同等性能)
- **CPU**: Intel Xeon Silver 4210R (或同等性能)
- **内存**: 128GB DDR4 3200MHz

#### 软件要求
- **操作系统**: Windows 10/11, Linux, macOS
- **Python**: 3.7
- **Go**: 1.22.5 或更高
- **Node.js**: 16.0 或更高
- **CUDA**: 10.2
- **cuDNN**: 7.6.5
- **Anaconda**: 推荐使用

### 安装步骤

#### 1. 克隆项目

```bash
git clone https://github.com/Sunny1004-7/SUNNY-AITEACH.git
cd SUNNY-AITEACH
```

#### 2. 前端安装

```bash
cd ISET-frontend

# 安装 pnpm (如果尚未安装)
npm install -g pnpm

# 安装依赖
pnpm install

# 启动开发服务器
pnpm serve
```

前端将在 `http://localhost:5173` 运行

#### 3. 后端安装

```bash
cd backend

# 安装Go依赖
go mod download

# 编译
go build -o backend.exe main.go

# 运行
./backend
```

后端API将在 `http://localhost:8010` 运行

#### 4. Python算法模块配置

```bash
cd exeGen

# 创建虚拟环境
python -m venv venv

# 激活虚拟环境
# Windows:
venv\Scripts\activate
# Linux/Mac:
source venv/bin/activate

# 安装依赖
pip install -r requirements.txt
```

#### 5. 课程推荐模块配置

```bash
cd courseRec/Multi-Agent-for-Course-Rec

# 下载MOOCCubeX数据集 (由于数据集过大，需本地导入)
# 将数据集放置在 data/MOOCCubeX 目录下

# 安装依赖
pip install -r requirements.txt
```

---

## 📁 项目结构

```
ISET/
├── backend/                      # Go后端服务
│   ├── lib/                      # 公共库
│   │   ├── model.go             # 数据模型
│   │   ├── service.go           # 业务服务
│   │   └── misc.go              # 工具函数
│   ├── main.go                  # 入口文件
│   ├── config.yaml              # 配置文件
│   └── data.db                  # SQLite数据库
│
├── ISET-frontend/               # Vue前端应用
│   ├── src/
│   │   ├── api/                 # API接口
│   │   ├── views/               # 页面组件
│   │   ├── components/          # 公共组件
│   │   ├── router/              # 路由配置
│   │   ├── store/               # 状态管理
│   │   └── utils/               # 工具函数
│   ├── public/                  # 静态资源
│   ├── package.json             # 项目配置
│   └── vite.config.ts           # Vite配置
│
├── exeGen/                      # 习题生成模块
│   ├── process.py               # 完整版算法 (4个智能体)
│   ├── process_no_discriminator.py  # 无质量评估版本
│   ├── process_no_kt.py         # 无知识追踪版本
│   ├── process_no_generator.py  # 无习题生成版本
│   ├── tea_process.py           # 教师端版本
│   ├── create_prompt.py         # 提示词生成
│   ├── dataSet.py               # 数据集处理
│   ├── knowledge_extraction.py  # 知识提取
│   ├── subject_data(1)/         # 训练数据
│   ├── bert/                    # BERT模型
│   └── requirements.txt         # Python依赖
│
├── courseRec/                   # 课程推荐模块
│   └── Multi-Agent-for-Course-Rec/
│       ├── main_handoffs.py     # 主程序
│       ├── agents/              # 智能体定义
│       ├── systems/             # 系统配置
│       ├── tools/               # 工具函数
│       ├── config/              # 配置文件
│       └── data/MOOCCubeX/      # 数据集
│
├── Socratic_QA_System/          # 苏格拉底式问答系统
│   ├── main.py                  # 主程序
│   ├── core.py                  # 核心组件
│   ├── teacher_agent.py         # 教师智能体
│   ├── monitor_agent.py         # 监控智能体
│   ├── knowledge_state_agent.py # 知识状态智能体
│   ├── config.py                # 配置管理
│   └── backend.py               # 后端接口
│
├── start.ps1                    # Windows启动脚本
├── 命令.txt                      # 命令参考
├── 示例输入输出.txt               # 示例数据
└── README.md                    # 项目说明
```

---

## 🤖 多智能体架构

### 习题生成模块 (exeGen)

系统采用4个专业智能体协作生成高质量习题：

#### 1. Knowledge Tracking Agent (知识追踪智能体)
- **职责**: 分析学生答题记录，评估知识掌握情况
- **输入**: 学生历史答题数据
- **输出**: 知识状态摘要（已掌握/薄弱知识点）

#### 2. Exercise Generator Agent (习题生成智能体)
- **职责**: 根据学生薄弱知识点生成针对性习题
- **输入**: 知识状态 + 题型要求
- **输出**: 习题初稿（题目、选项、答案、解析）

#### 3. Discriminator Agents (质量评估智能体 x3)
- **Linguistic Fluency**: 评估语言流畅性、表达清晰度
- **Knowledge Coverage**: 检查知识点覆盖准确性
- **Correctness & Reasonableness**: 验证答案正确性和难度合理性

#### 4. Supervisor Agent (监督智能体)
- **职责**: 协调各智能体，管理迭代优化流程
- **流程**: Generator → Discriminators → Regenerate → Final Output

### 工作流程

```
学生答题记录 → Knowledge Tracking Agent
                    ↓
              知识状态分析
                    ↓
          Exercise Generator Agent
                    ↓
              生成习题初稿
                    ↓
         ┌─────────┴─────────┐
         ↓          ↓         ↓
    语言评估    知识评估    正确性评估
         ↓          ↓         ↓
         └─────────┬─────────┘
                   ↓
            质量检查通过？
              ↓ No    ↓ Yes
           重新生成   输出最终习题
```

---

## 📊 数据集说明

### MOOCCubeX
- **用途**: 课程推荐模块训练数据
- **位置**: `courseRec/Multi-Agent-for-Course-Rec/data/MOOCCubeX/`
- **说明**: 由于数据集较大，需要手动下载并放置到指定目录

### 习题数据集
- **位置**: `exeGen/subject_data(1)/`
- **文件**:
  - `stuRec_1000.csv`: 学生答题记录
  - `examples_with_explanation_with_tokens_KG.csv`: 示例习题库
  - `concept_relationship_filtered.csv`: 知识点关系
  - `problem_concept.csv`: 题目-知识点映射

---

## 🔌 API接口

### 习题生成接口

#### 学生端
```http
GET /api/ExeAlgorithm-stu
```

**参数**:
- `exerciseType`: 题型 (Single_Choice/Multiple_Choice/True_or_False)
- `exerciseQuantity`: 数量
- `selectedModel`: 模型选择
- `selectedAgents`: 智能体组合 (逗号分隔)

**智能体组合**:
- 完整版: `KnowledgeAwareAgent,DiscriminatorAgent,SuperviseAgent,ExeAgent`
- 无质量评估: `KnowledgeAwareAgent,SuperviseAgent,ExeAgent`
- 无知识追踪: `DiscriminatorAgent,SuperviseAgent,ExeAgent`

#### 教师端
```http
GET /api/ExeAlgorithm-tea
```

**额外参数**:
- `knowledgePoint`: 指定知识点
- `difficulty`: 难度等级
- `language`: 编程语言
- `subject`: 学科类别

### WebSocket接口
```
ws://localhost:8010/api/ws/history
```
实时推送算法执行进度和智能体对话历史

### 知识图谱接口
- `GET /graph`: 获取所有图谱
- `POST /graph`: 创建新图谱
- `GET /node/by_graph/:gid`: 获取图谱节点
- `GET /link/by_graph/:gid`: 获取图谱边

---

## 🎯 使用示例

### 示例1：学生端习题生成

**输入**:
```json
{
  "exerciseType": "Multiple_Choice",
  "exerciseQuantity": "10",
  "selectedModel": "qwen-plus-latest",
  "selectedAgents": "KnowledgeAwareAgent,DiscriminatorAgent,SuperviseAgent,ExeAgent"
}
```

**输出** (部分):
```json
{
  "exercises": [
    {
      "content": "阿里云伏羲分布式调度系统与社区Hadoop MR的处理能力对比，哪个系统更容易扩展？",
      "option": {
        "A": "伏羲",
        "B": "Hadoop",
        "C": "两者相同",
        "D": "无法比较"
      },
      "answer": ["A"],
      "knowledge_concept": "K_调度策略_计算机科学与技术",
      "explanation": "伏羲采用独立的两层架构资源与任务调度，具有更强的横向扩展性..."
    }
  ],
  "chatHistory": [...]
}
```

### 示例2：课程推荐

**输入**: 学生ID `U_10258447`

**输出**:
```json
{
  "user_id": "U_10258447",
  "历史课程": "大学计算机 — 计算思维导论, 人工智能概论, C语言程序设计精髓...",
  "知识水平评估": "目前处于计算机科学基础认知向技术深化过渡的阶段...",
  "学习偏好": "学习路径呈现系统性递进、理论与实践结合的特点...",
  "推荐课程名称列表": [
    "大学计算机基础 CAP",
    "网络技术与应用",
    "计算机组成原理（上）",
    "C语言程序设计进阶",
    "Java程序设计",
    "数据库系统原理"
  ],
  "推荐原因": "学生已具备计算机科学的基础认知，需在硬件原理、编程进阶等领域重点发力..."
}
```

---

## 🔧 配置说明

### 后端配置 (backend/config.yaml)

```yaml
databaseconfig:
  type: "sqlite"
  db: "data.db"
  migrate: true
port: ":8010"
```

### LLM配置 (exeGen/process.py)

```python
llm_config = {
    "cache_seed": None,
    "config_list": [{
        "model": "qwen-plus-latest",
        "base_url": "https://dashscope.aliyuncs.com/compatible-mode/v1",
        "api_key": "your-api-key",
        "price": [0, 0]
    }]
}
```

**支持的模型**:
- Qwen (qwen-plus-latest, qwen-turbo, qwen-max)
- OpenAI GPT (gpt-4o, gpt-4o-mini, gpt-3.5-turbo)
- DeepSeek
- Llama

---

## 📖 文档

### 设计文档
- [系统架构说明](ISET-frontend/doc/ISET-sys-arch.md)
- [视频演示脚本](ISET-frontend/doc/ISET-vid-note.md)

### README文件
- [前端开发指南](ISET-frontend/README.md)
- [习题生成模块](exeGen/README.md)
- [课程推荐模块](courseRec/Multi-Agent-for-Course-Rec/README.md)

---

## 🎓 学术价值

### 创新点

1. **多智能体协作机制**
   - 首次将AutoGen框架应用于教育领域习题生成
   - 实现知识追踪、生成、评估的智能体分工协作

2. **个性化学习路径**
   - 基于知识图谱的学习状态建模
   - 动态生成适配学生能力的习题

3. **情绪感知教学**
   - 苏格拉底式问答融合情绪识别
   - ICECoT流程：情绪分析→意图推断→策略选择→响应生成

4. **大规模知识图谱应用**
   - Neo4j存储课程知识关系
   - 支持多维度知识图谱可视化

### 应用场景

- ✅ 在线教育平台智能出题
- ✅ 高校教学辅助系统
- ✅ 个性化学习路径推荐
- ✅ 知识掌握度评估
- ✅ 自适应学习系统

---

## 🤝 贡献指南

我们欢迎所有形式的贡献！

### 如何贡献

1. **Fork** 本仓库
2. **创建**特性分支 (`git checkout -b feature/AmazingFeature`)
3. **提交**更改 (`git commit -m 'Add some AmazingFeature'`)
4. **推送**到分支 (`git push origin feature/AmazingFeature`)
5. **提交** Pull Request

### 代码规范

- **Python**: 遵循 PEP 8 规范
- **Go**: 遵循 Go 官方代码规范
- **JavaScript/TypeScript**: 遵循 ESLint 规则
- **提交信息**: 使用 Conventional Commits 格式

---

## 📄 许可证

本项目采用 [MIT License](LICENSE) 开源协议。

---

## 👥 团队

### 项目开发
- **西安交通大学** - 教改项目组
- **课题**: 机器学习教改-AI出题模块

### 致谢
- AutoGen团队提供的多智能体框架
- MOOCCubeX数据集支持
- Pure Admin前端模板

---

## 📞 联系方式

- **项目主页**: [https://github.com/Sunny1004-7/SUNNY-AITEACH](https://github.com/Sunny1004-7/SUNNY-AITEACH)
- **问题反馈**: [Issues](https://github.com/Sunny1004-7/SUNNY-AITEACH/issues)
- **讨论区**: [Discussions](https://github.com/Sunny1004-7/SUNNY-AITEACH/discussions)

---

## ⭐ Star History

如果这个项目对你有帮助，请给我们一个 ⭐️ Star！

[![Star History Chart](https://api.star-history.com/svg?repos=Sunny1004-7/SUNNY-AITEACH&type=Date)](https://star-history.com/#Sunny1004-7/SUNNY-AITEACH&Date)

---

<div align="center">

**让教育更智慧，让学习更个性化**

Made with ❤️ by XJTU Team

</div>
