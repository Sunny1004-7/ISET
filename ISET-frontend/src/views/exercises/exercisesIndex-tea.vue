<template>
  <div class="page-container">
    <!-- 左侧主要区域 -->
    <div class="left-section">
      <!-- 智能体交互可视化窗口 -->
      <div class="interaction-visualization">
        <h2 class="text-center">智能体交互可视化窗口</h2>
        <pre v-if="!haschatHistory" class="formattedMessages">{{ formattedMessages }}</pre>
        <!-- 对话式内容展示 -->
        <div v-else class="chat-container">
          <div v-for="(msg, index) in chatHistory" :key="index" :class="[
            'chat-bubble',
            msg.role,
            getBubbleClass(msg.name)
          ]">
            <img :src="getAvatar(msg.name)" class="avatar" :alt="msg.name">
            <div class="bubble-content">
              <div class="agent-name">{{ formatName(msg.name) }}</div>
              <pre class="message">{{ filterContent(msg.content) }}</pre>
            </div>
          </div>
        </div>
      </div>
      <!-- 习题生成结果（可编辑窗口） -->
      <div class="generated-exercises">
        <h2 class="text-center">习题生成结果（可编辑窗口）</h2>
        <textarea v-model="exerciseContent" placeholder="习题内容" class="exerciseContent-textarea"></textarea>
      </div>
      <div class="button-area">
        <!-- 参数设置区域 -->
        <div class="param-container">
          <!-- 难度、语言、科目 -->
          <div class="param-row compact">
            <div class="param-group">
              <select v-model="difficulty" class="param-select">
                <option value="">题型难度</option>
                <option value="简单">简单</option>
                <option value="中等">中等</option>
                <option value="困难">困难</option>
              </select>
            </div>
            <div class="param-group">
              <select v-model="language" class="param-select">
                <option value="">语言</option>
                <option value="汉语">汉语</option>
                <option value="英语">英语</option>
              </select>
            </div>
            <div class="param-group">
              <select v-model="subject" class="param-select">
                <option value="">科目</option>
                <option value="数学">数学</option>
                <option value="计算机科学与技术">计算机科学与技术</option>
                <option value="机器学习">机器学习</option>
              </select>
            </div>
          </div>
          <div class="param-row">
            <input v-model="knowledgePoint" type="text" placeholder="请输入知识点" class="knowledge-input">
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧选择和设置区域 -->
    <div class="right-section">
      <!-- 智能体选择区 -->
      <div class="agent-selection">
        <div class="section-header">
          <img src="/images/箭头.png" class="title-icon">
          <h2 class="section-title">智能体选择区</h2>
        </div>
        <div class="selection-grid">
          <div v-for="(agent, index) in agents" :key="index" class="selection-wrapper">
            <div class="selection-item" :style="{ backgroundImage: 'url(' + agent.imageUrl + ')' }"
              :class="{ checked: agent.checked }">
              <div class="image-placeholder"></div>
            </div>
            <div class="label-container">
              <span class="selection-name">{{ agent.name }}</span>
              <input type="checkbox" v-model="agent.checked" class="checkbox" :disabled="agent.disabled">
            </div>
          </div>
        </div>
      </div>

      <!-- 大模型选择区 -->
      <div class="model-selection">
        <div class="section-header">
          <img src="/images/箭头.png" class="title-icon">
          <h2 class="section-title">大模型选择区</h2>
        </div>
        <div class="selection-grid">
          <div v-for="(model, index) in models" :key="index" class="selection-wrapper">
            <div class="selection-item" :style="{ backgroundImage: 'url(' + model.imageUrl + ')' }"
              :class="{ checked: selectedModel === model.name }">
              <div class="image-placeholder"></div>
            </div>
            <div class="label-container">
              <span class="selection-name">{{ model.name }}</span>
              <input type="radio" v-model="selectedModel" :value="model.name" class="radio">
            </div>
          </div>
        </div>
      </div>

      <!-- 参数设置区 -->
      <div class="parameter-setting">
        <div class="section-header">
          <img src="/images/箭头.png">
          <h2 class="section-title">参数设置区</h2>
        </div>
        <div class="parameter-row">
          <div class="parameter-item">
            <select v-model="exerciseType" class="exercise-type-select" style="width: 100%;">
              <option value="">下拉框 题型</option>
              <option value="单选题">单选题</option>
              <option value="多选题">多选题</option>
              <option value="判断题">判断题</option>
            </select>
          </div>
          <div class="parameter-item">
            <input v-model="exerciseQuantity" type="text" placeholder="输入框 数量" class="quantity-input"
              style="width: 100%;">
          </div>
        </div>
      </div>

      <hr class="separator-line" />

      <!-- 操作按钮区域 -->
      <div class="operButton-area">
        <button class="regenerate-button" @click="submitParameters">生成习题</button>
        <button class="export-button" @click="exportExercises">导出按钮</button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue';
import { saveAs } from 'file-saver';
import axios from 'axios';
import { ElMessageBox, ElMessage } from 'element-plus';

const haschatHistory = ref(false)
const chatHistory = ref([])

const filterContent = (text) => {
  return text
    .replace(/```json|```/g, '')
    .replace(/[{}[\]]/g, '')
    .replace(/\*\*(.*?)\*\*/g, '$1')
    .replace(/\n{2,}/g, '\n')
    .replace(/(\n\s*){2,}/g, '\n')
    .trim();
}
const getBubbleClass = (name) => {
  const map = {
    'agent_host': 'supervise-bubble',
    'agent_generator': 'exercise-bubble',
    'Linguistic_Fluency_discriminator': 'evaluator-bubble',
    'Knowledge_Concept_Coverage_discriminator': 'evaluator-bubble',
    'Correctness_and_Reasonableness_discriminator': 'evaluator-bubble'
  }
  return map[name] || ''
}
// 头像配置
const avatarConfig = {
  agent_host: '/images/监管.png',
  agent_generator: '/images/习题生成.png',
  Linguistic_Fluency_discriminator: '/images/质量评估.png',
  Knowledge_Concept_Coverage_discriminator: '/images/质量评估.png',
  Correctness_and_Reasonableness_discriminator: '/images/质量评估.png'
}

const isEvaluator = (name) => {
  const evaluators = [
    'Linguistic_Fluency_discriminator',
    'Knowledge_Concept_Coverage_discriminator',
    'Correctness_and_Reasonableness_discriminator'
  ];
  return evaluators.includes(name);
}

const getAvatar = (name) => {
  if (isEvaluator(name)) {
    return '/images/质量评估.png';
  }
  return avatarConfig[name] || '/images/质量评估.png';
}

const formatName = (name) => {
  const nameMap = {
    agent_host: 'Generation Manager',
    agent_generator: 'Exercise Generator',
    Linguistic_Fluency_discriminator: 'Evaluation Experts',
    Knowledge_Concept_Coverage_discriminator: 'Evaluation Experts',
    Correctness_and_Reasonableness_discriminator: 'Evaluation Experts'
  }
  return nameMap[name] || name;
}

const formattedMessages = computed(() => {
  return interactionMessages.value
    .map(m => m.content)
    .join('\n')
})

// WebSocket连接
const socket = ref(null);


const interactionMessages = ref([]);


const exerciseContent = ref('');

const difficulty = ref('');
const knowledgePoint = ref('');
const language = ref('');
const subject = ref('');
const exerciseType = ref('');
const exerciseQuantity = ref('');
const selectedModel = ref('');

// 智能体列表
const agents = ref([
  { id: 1, name: '生成监管者', checked: true, imageUrl: '/images/监管.png', disabled: false },
  { id: 2, name: '习题生成者', checked: true, imageUrl: '/images/习题生成.png', disabled: false },
  { id: 3, name: '评估专家组', checked: true, imageUrl: '/images/质量评估.png', disabled: false }
]);
// 大模型列表
const models = ref([
  { id: 1, name: 'LLaMa', imageUrl: '/images/llama.png' },
  { id: 2, name: 'DeepSeek', imageUrl: '/images/deepseek.png' },
  { id: 3, name: 'Qwen', imageUrl: '/images/Qwen.png' },
  { id: 4, name: 'GPT-4o', imageUrl: '/images/gpt.png' }
]);
// WebSocket 连接函数
const connectWebSocket = () => {

  if (socket.value) {
    return;
  }
  const ws = new WebSocket('ws://localhost:8010/api/ws/history');

  ws.onopen = () => {
    console.log('WebSocket连接成功');
  };

  ws.onmessage = (event) => {
    const message = event.data;
    if (message.trim() === "") return;

    console.log("收到消息:", message);
    interactionMessages.value.push({
      timestamp: new Date().toLocaleTimeString(),
      content: message
    });

    // 自动滚动到底部
    const textarea = document.querySelector('.interaction-visualization textarea');
    if (textarea) {
      textarea.scrollTop = textarea.scrollHeight;
    }
  };

  // WebSocket 连接关闭
  ws.onclose = () => {
    console.log('WebSocket连接关闭');

    setTimeout(() => {
      connectWebSocket();
    }, 1000);
  };


  // WebSocket 连接错误处理
  ws.onerror = (error) => {
    console.error("WebSocket错误:", error);
  };

  socket.value = ws;
};

// 页面销毁时关闭 WebSocket 连接
onUnmounted(() => {
  if (socket.value) {
    socket.value.close();
  }
});

// 导出功能
const exportExercises = () => {
  const blob = new Blob([exerciseContent.value], { type: 'text/plain;charset=utf-8' });
  saveAs(blob, 'exercises.txt');
};

// 提交参数
const submitParameters = async () => {
  const selected = agents.value.filter(agent => agent.checked).map(agent => agent.name);

  const exerciseTypeMapping = {
    "单选题": "Single_Choice",
    "多选题": "Multiple_Choice",
    "判断题": "True_or_False"
  };

  const mappedExerciseType = exerciseTypeMapping[exerciseType.value] || "";

  const agentAliasMap = {
    "生成监管者": "SuperviseAgent",
    "习题生成者": "ExeAgent",
    "评估专家组": "DiscriminatorAgent"
  };

  const mappedselectedAgents = selected
    .map(agent => agentAliasMap[agent] || agent)
    .filter(Boolean);

  const infoMessage = {
    exerciseType: exerciseType.value,
    exerciseQuantity: exerciseQuantity.value,
    selectedAgents: selected.join(', '),
    selectedModel: selectedModel.value,
    difficulty: difficulty.value,
    knowledgePoint: knowledgePoint.value,
    language: language.value,
    subject: subject.value
  };

  // 弹出卡片并显示各个信息
  ElMessageBox.confirm(
    `<div style="display: flex; flex-direction: column; gap: 10px;">
        <span>题型: ${infoMessage.exerciseType}</span>
        <span>数量: ${infoMessage.exerciseQuantity}</span>
        <span>难度: ${infoMessage.difficulty}</span>
        <span>知识点: ${infoMessage.knowledgePoint}</span>
        <span>语言: ${infoMessage.language}</span>
        <span>科目: ${infoMessage.subject}</span>
        <span>智能体: ${infoMessage.selectedAgents}</span>
        <span>大模型: ${infoMessage.selectedModel}</span>
    </div>`,
    '生成习题',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      dangerouslyUseHTMLString: true,
      type: 'info'
    }
  ).then(async () => {

    if (!exerciseType.value || !exerciseQuantity.value || !selectedModel.value ||
      !difficulty.value || !language.value || !subject.value) {
      ElMessageBox.alert('请确保所有字段已填写', '提示', {
        confirmButtonText: '确定'
      });
      return;
    }
    try {
      const response = await axios.get('http://localhost:8010/api/ExeAlgorithm-tea', {
        params: {
          exerciseType: mappedExerciseType,
          exerciseQuantity: exerciseQuantity.value,
          selectedModel: selectedModel.value,
          selectedAgents: mappedselectedAgents.join(', '),
          difficulty: difficulty.value,
          knowledgePoint: knowledgePoint.value,
          language: language.value,
          subject: subject.value
        }
      });

      const jsonString = response.data.data;
      const parsedData = JSON.parse(jsonString);
      const exercises = parsedData.exercises || [];
      chatHistory.value = parsedData.chatHistory || []

      if (chatHistory.value.length > 0) {
        haschatHistory.value = true
      }

      if (Array.isArray(exercises)) {
        exerciseContent.value = exercises.map((exercise, index) => {
          console.log(exercise);
          return `
      Exercise ${index + 1}
      Content：${exercise.content}
      Option：${formatOptions(exercise.option)}
      Answer：${exercise.answer}
      Knowledge-Concept：${exercise.concept_id}
      Course-Name：${exercise.course_name}
      Explanation：：${exercise.explanation}
      `;
        }).join('');
      } else {
        exerciseContent.value = '后端返回的数据不符合预期';
      }
    } catch (error) {
      console.error("Network Error:", error);
      ElMessageBox.alert(`发送失败: ${error.message}`, '错误', {
        confirmButtonText: '确定'
      });
    }
  }).catch(() => {
    ElMessageBox.alert('生成习题已取消', '提示', {
      confirmButtonText: '确定'
    });
  });
};
const formatOptions = (options) => {
  try {
    const obj = typeof options === 'string' ? JSON.parse(options) : options;
    return Object.entries(obj).map(([k, v]) => `${k}: ${v}`).join("\n      ");
  } catch {
    return options;
  }
};
// 生命周期
onMounted(connectWebSocket);
onUnmounted(() => socket.value?.close());
</script>

<style scoped>
.page-container {
  display: flex;
  width: 100%;
  min-height: 90%;

}

.left-section {
  flex: 2.5;
  display: flex;
  flex-direction: column;
  padding: 5px;
  background-color: #ffffff;
  height: 90vh;
  border-radius: 10px;
  border: 1px solid #ccc;

  gap: 0px;
}

.right-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 0px;
  background-color: #f9f9f9;
  gap: 0px;
  height: 90vh;
  border-radius: 10px;
  border: 1px solid #ccc;

}

.interaction-visualization,
.generated-exercises {
  flex: 1;
  padding: 10px;
  border: none;
  background-color: white;
  width: 100%;
  height: 100%;
  margin-bottom: 4px;
}

.chat-container {
  height: 40vh;

  padding: 16px;
  border: 2px solid #f0f0f0;
  border-radius: 12px;
  background: #ffffff;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
  overflow-y: auto;
  scroll-behavior: smooth;

  display: flex;
  flex-direction: column;
  gap: 20px;
  margin-top: 16px;
}


.chat-container::-webkit-scrollbar {
  width: 8px;
  background-color: #f8f9fa;
}

.chat-container::-webkit-scrollbar-thumb {
  background-color: #d4d4d4;
  border-radius: 8px;
}

.chat-container::-webkit-scrollbar-track {
  background: transparent;
}

.chat-bubble {
  display: flex;
  margin: 16px 0;
  gap: 12px;
}

.evaluator-bubble .bubble-content {
  background: #fff3e0 !important;
  border: 2px solid #ffb74d;
  border-radius: 12px;
}

.evaluator-bubble .agent-name {
  color: #ef6c00;
  font-weight: 600;
}


.supervise-bubble .bubble-content {
  background: #e3f2fd !important;
  border: 2px solid #64b5f6;
  border-radius: 12px 12px 12px 4px;
}

.supervise-bubble .agent-name {
  color: #1e88e5;
}


.knowledge-bubble .bubble-content {
  background: #e8f5e9 !important;
  border: 2px solid #81c784;
  border-radius: 12px 12px 4px 12px;
}

.knowledge-bubble .agent-name {
  color: #2e7d32;
}

.exercise-bubble .bubble-content {
  background: #f3e5f5 !important;
  border: 2px solid #ba68c8;
  border-radius: 4px 12px 12px 12px;
}

.exercise-bubble .agent-name {
  color: #9c27b0;
}

.bubble-content {
  max-width: 75%;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  word-wrap: break-word;
  overflow-wrap: break-word;
}

.agent-name {
  font-size: 20px;
  margin-bottom: 4px;
}

.message {
  white-space: pre-wrap;
  font-family: system-ui;
  line-height: 1.6;
  color: #070707;
  margin: 0;
  word-wrap: break-word;
  overflow-wrap: break-word;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 6px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.button-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.quality-dialog :deep(.el-dialog__header) {
  border-bottom: 1px solid #ebeef5;
  padding: 20px;
  margin: 0;
}

.quality-dialog :deep(.el-dialog__title) {
  font-size: 18px;
  font-weight: 500;
  color: #303133;
}

.quality-dialog :deep(.el-dialog__body) {
  padding: 20px;
}

.quality-options {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.quality-option {
  display: flex;
  align-items: center;
  padding: 16px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  background: #ffffff;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.quality-option:hover {
  border-color: #409eff;
  box-shadow: 0 2px 12px rgba(32, 107, 196, 0.12);
  transform: translateY(-2px);
}

.quality-option.is-selected {
  border-color: #409eff;
  background-color: #f0f7ff;
}

.quality-content {
  flex: 1;
  margin-right: 20px;
}

.quality-title {
  margin: 0 0 8px 0;
  font-size: 15px;
  font-weight: 500;
  color: #303133;
}

.quality-desc {
  margin: 0;
  font-size: 13px;
  line-height: 1.5;
  color: #606266;
}

.quality-checkbox :deep(.el-checkbox__input) .el-checkbox__inner {
  width: 18px;
  height: 18px;
  border-radius: 4px;
}

.quality-checkbox :deep(.el-checkbox__input) .el-checkbox__inner::after {
  left: 6px;
  top: 2px;
  width: 5px;
  height: 10px;
}

.quality-checkbox :deep(.el-checkbox__label) {
  display: none;
}

.knowledge-input {
  width: 100%;
  padding: 8px 12px;
  border: 2px solid #c5e7fd;
  border-radius: 4px;
  font-size: 20px;
}


.param-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
}

.param-row {
  display: flex;
  gap: 20px;
  height: 50px;
}

.param-group {
  flex: 1;
  min-width: 120px;
}

.param-select,
.param-input {
  width: 100%;
  padding: 6px;
  border: 2px solid #c5e7fd;
  border-radius: 4px;
  font-size: 14px;
  height: 40px;
}

.operButton-area {
  display: flex;
  justify-content: space-around;
  margin-top: 0px;
  height: 30%;
  background-color: #ffffff;
}

.regenerate-button,
.export-button {
  padding: 15px 25px;
  border: 2px solid #80e27e;
  border-radius: 5px;
  cursor: pointer;
  font-size: 20px;
  font-weight: bold;
  transition: background-color 0.3s ease, color 0.3s ease;
  height: 40%;
  margin-top: 40px;
}

.regenerate-button {
  background-color: white;
  color: black;
  border-color: #80e27e;
  margin-bottom: 50px;
}

.regenerate-button:hover {
  background-color: #80e27e;
  color: white;
}

.export-button {
  background-color: #4CAF50;
  color: white;
  border: none;
}

.export-button:hover {
  background-color: #45a049;
  color: white;
}

.action-buttons {
  display: flex;
  justify-content: space-around;
  margin-top: 0px;
  height: 30%;
  background-color: #ffffff;
}

.agent-selection,
.model-selection {
  padding: 10px;
  border-radius: 5px;
  background-color: #ffffff;
  margin-bottom: 0px;
  height: 50%;
}

.parameter-setting {
  padding: 10px;
  border-radius: 5px;
  background-color: #ffffff;
  margin-bottom: 0px;
  height: 30%;
}

.exercise-type-select,
.quantity-input {
  margin-top: 5px;
  padding: 5px;
  border: 1px solid #ffffff;
  border-radius: 3px;
  height: 80%;
  width: 100%;
}

.selection-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  row-gap: 25px;
  column-gap: 15px;
  justify-items: stretch;
  align-items: center;
  background-color: #ffffff;
  margin-top: 10px;
  margin-bottom: 20px;
}

.selection-wrapper {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-around;
  gap: 5px;
  background-color: #b8d0f6;
  border-radius: 10px;
  border: 0px solid #ffffff;
  height: 80px;
  box-sizing: border-box;
}

.agent-selection .selection-item {
  color: #333;
  border-radius: 5px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-between;
  padding: 10px;
  aspect-ratio: 2/3;
  width: 70px;
  height: 70px;
  background-size: cover;
  background-position: center;
  position: relative;
  margin-left: 10px;
}

.model-selection .selection-item {
  color: #333;
  border-radius: 5px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-between;
  padding: 10px;
  aspect-ratio: 2/1;
  width: 65px;
  height: 55px;
  background-size: cover;
  background-position: center;
  position: relative;
  margin-left: 30px;
}

.selection-name {
  font-size: 20px;
  color: #333;
  margin-right: 10px;
}

.image-placeholder {
  width: 100%;
  height: 100%;
  background-color: transparent;
  margin-bottom: 5px;
}

.checkbox {
  margin-left: 5px;
  width: 20px;
  height: 20px;
  margin-left: 0px;
  margin-right: 10px;
}

.radio {
  margin-left: 5px;
  width: 20px;
  height: 20px;
  margin-left: 10px;
  margin-right: 10px;
}

.parameter-row {
  display: flex;
  justify-content: space-between;
  gap: 3px;
  margin-top: 15px;
}

.parameter-item {
  background-color: #fafafa;
  color: #333;
  padding: 5px;
  border-radius: 5px;
  flex: 1;
  border: 2px solid #addaf5;
  font-size: 20px;
}


.interaction-textarea {
  width: 100%;
  height: 400px;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
  background-color: #ffffff;
  font-family: Arial, sans-serif;
  font-size: 14px;
  color: #333;
  resize: vertical;
  overflow-y: auto;
  white-space: pre-wrap;
}

.exerciseContent-textarea {
  width: 100%;
  height: 300px;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
  background-color: #ffffff;
  font-family: Arial, sans-serif;
  font-size: 14px;
  color: #333;
  resize: vertical;
  overflow-y: auto;
  white-space: pre-wrap;
}

.separator-line {
  border-top: 1px solid #ffffff;
  margin: 3px 0;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 10px;
}

.section-title {
  font-size: 30px;
}

h2.text-center {
  text-align: left;
  font-size: 35px;
}
</style>