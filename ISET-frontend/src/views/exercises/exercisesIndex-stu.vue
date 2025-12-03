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
      <div class="button-area">
        <button class="regenerate-button" @click="submitParameters">生成习题</button>
        <button class="export-button" @click="exportExercises">导出按钮</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { saveAs } from 'file-saver';
import axios from 'axios';
import { ElMessageBox, ElMessage } from 'element-plus';
import { computed } from 'vue'

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
    'agent_kt': 'knowledge-bubble',
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
  agent_kt: '/images/知识感知.png',
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

// 头像获取逻辑
const getAvatar = (name) => {
  if (isEvaluator(name)) {
    return '/images/质量评估.png';
  }
  return avatarConfig[name] || '/images/质量评估.png';
}

// 名称格式化
const formatName = (name) => {
  const nameMap = {
    agent_host: '生成监管者',
    agent_kt: '知识感知器',
    agent_generator: '习题生成器',
    Linguistic_Fluency_discriminator: '流畅性评估者',
    Knowledge_Concept_Coverage_discriminator: '知识概念覆盖评估者',
    Correctness_and_Reasonableness_discriminator: '正确性与合理性评估者'
  }
  return nameMap[name] || name;
}

const formattedMessages = computed(() => {
  return interactionMessages.value
    .map(m => m.content)
    .join('\n')
})
const socket = ref(null);

const interactionMessages = ref([]);

// 生成的习题内容
const exerciseContent = ref('');
const selectedModel = ref('');
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
  };


  ws.onclose = () => {
    console.log('WebSocket连接关闭');

    setTimeout(() => {
      connectWebSocket();
    }, 1000);
  };



  ws.onerror = (error) => {
    console.error("WebSocket错误:", error);
  };

  socket.value = ws;
};


onUnmounted(() => {
  if (socket.value) {
    socket.value.close();
  }
});

// 导出函数
const exportExercises = () => {
  const data = exerciseContent.value;
  const filename = 'exercises.txt';
  const blob = new Blob([data], { type: 'text/plain;charset=utf-8' });
  saveAs(blob, filename);
};

// 选中的题型
const exerciseType = ref('');
// 习题数量
const exerciseQuantity = ref('');
// 智能体列表
const agents = ref([
  { id: 1, name: '生成监管者', checked: true, imageUrl: '/images/监管.png', disabled: false },
  { id: 2, name: '知识感知器', checked: true, imageUrl: '/images/知识感知.png', disabled: false },
  { id: 3, name: '习题生成器', checked: true, imageUrl: '/images/习题生成.png', disabled: false },
  { id: 4, name: '评估专家', checked: true, imageUrl: '/images/质量评估.png', disabled: false }
]);
// 大模型列表
const models = ref([
  { id: 1, name: 'LLaMa', imageUrl: '/images/llama.png' },
  { id: 2, name: 'DeepSeek', imageUrl: '/images/deepseek.png' },
  { id: 3, name: 'Qwen', imageUrl: '/images/Qwen.png' },
  { id: 4, name: 'GPT-4o', imageUrl: '/images/gpt.png' }
]);

const submitParameters = async () => {
  // 获取选中的智能体
  const selected = agents.value.filter(agent => agent.checked).map(agent => agent.name);

  // 映射题型名称
  const exerciseTypeMapping = {
    "单选题": "Single_Choice",
    "多选题": "Multiple_Choice",
    "判断题": "True_or_False"
  };

  // 将中文题型转换为英文格式
  const mappedExerciseType = exerciseTypeMapping[exerciseType.value] || "";

  const agentAliasMap = {
    "生成监管者": "SuperviseAgent",
    "习题生成器": "ExeAgent",
    "知识感知器": "KnowledgeAwareAgent",
    "评估专家": "DiscriminatorAgent"
  };

  const mappedselectedAgents = selected
    .map(agent => agentAliasMap[agent] || agent)
    .filter(Boolean);

  // 设置要显示的信息，使用原始的中文名称
  const infoMessage = {
    exerciseType: exerciseType.value,
    exerciseQuantity: exerciseQuantity.value,
    selectedAgents: selected.join(', '),
    selectedModel: selectedModel.value,
  };

  // 弹出卡片并显示各个信息
  ElMessageBox.confirm(
    `<div style="display: flex; flex-direction: column; gap: 10px;">
      <span>题型: ${infoMessage.exerciseType}</span>
      <span>数量: ${infoMessage.exerciseQuantity}</span>
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

    if (!exerciseType.value || !exerciseQuantity.value || !selectedModel.value || selected.length === 0) {
      ElMessageBox.alert('请确保所有字段已填写', '提示', {
        confirmButtonText: '确定'
      });
      return;
    }
    try {

      const response = await axios.get('http://localhost:8010/api/ExeAlgorithm-stu', {
        params: {
          exerciseType: mappedExerciseType,
          exerciseQuantity: exerciseQuantity.value,
          selectedModel: selectedModel.value,
          selectedAgents: mappedselectedAgents.join(', ')
        }
      });

      const jsonString = response.data.data;
      const parsedData = JSON.parse(jsonString);

      const exercises = parsedData.exercises || [];
      chatHistory.value = parsedData.chatHistory || []

      if (chatHistory.value.length > 0) {
        haschatHistory.value = true
      }
      // 格式化并更新习题内容
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
      Explanation：${exercise.explanation}
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

// 页面加载时连接WebSocket
onMounted(() => {
  if (!socket.value) {
    connectWebSocket();
  }
});

</script>

<style scoped>
.page-container {
  display: flex;
  width: 90%;
  height: 90vh;
  min-height: 600px;
}

.left-section {
  flex: 2.5;
  display: flex;
  flex-direction: column;
  padding: 5px;
  background-color: #ffffff;
  height: 100%;
  border-radius: 10px;
  border: 1px solid #ccc;
  gap: 0px;
}

.right-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 0px;
  background-color: #fafafa;
  gap: 0px;
  height: 100%;
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
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 6px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}


.agent-selection,
.model-selection {
  padding: 10px;
  border-radius: 5px;
  background-color: #ffffff;
  margin-bottom: 0px;
  height: 30vh;
}

.parameter-setting {
  padding: 10px;
  border-radius: 5px;
  background-color: #ffffff;
  border-radius: 10px;
  margin-bottom: 0px;
  height: 20vh;
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
  background-color: #ffffff;
  margin-top: 15px;
}

.parameter-item {
  background-color: #ffffff;
  color: #333;
  padding: 5px;
  border: 2px solid #a1d1e3;
  flex: 1;
}

.regenerate-button,
.export-button {
  padding: 15px 25px;
  border: 2px solid #80e27e;
  border-radius: 5px;
  cursor: pointer;
  font-size: 25px;
  font-weight: bold;
  transition: background-color 0.3s ease, color 0.3s ease;
  height: 50%;
  margin-left: 20px;
  margin-right: 20px;
  margin-top: 40px;
}

.regenerate-button {
  background-color: white;
  color: black;
  border-color: #80e27e;
  margin-bottom: 50px;
  width: 20vh;
}

.regenerate-button:hover {
  background-color: #80e27e;
  color: white;
}

.export-button {
  background-color: #4CAF50;
  color: white;
  border: none;
  width: 20vh;
}

.export-button:hover {
  background-color: #45a049;
  color: white;
}

.button-area {
  display: flex;
  justify-content: space-around;
  margin-top: 0px;
  background-color: #ffffff;
  border-radius: 10px;
}

.separator-line {
  border-top: 1px solid #ffffff;
  margin: 3px 0;
}

.formattedMessages {
  width: 100%;
  height: 40vh;
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