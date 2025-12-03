<template>
  <div class="socratic-teach">
    <el-container class="teach-container">
      <el-header class="teach-header">
        <div class="header-content">
          <h1>苏格拉底式教学</h1>
          <div class="thinking-mode-switch">
            <el-switch v-model="isSlowThinking" size="large" active-text="慢思考" inactive-text="快思考"
              active-color="#2c5aa0" inactive-color="#13ce66" @change="onThinkingModeChange" />
          </div>
        </div>
      </el-header>

      <!-- 教学对话区域 -->
      <el-main class="teach-main">
        <!-- 慢思考模式显示处理信息 -->
        <div v-if="isSlowThinking && processingInfo" class="processing-info">
          <el-alert :title="processingInfo" type="info" :closable="false" show-icon />
        </div>


        <div class="dialogue-box">
          <div v-for="(msg, index) in messages" :key="index" class="message-container"
            :class="{ 'user-message': msg.role === 'user', 'mentor-message': msg.role === 'assistant' }">
            <!-- 引导者消息 -->
            <div v-if="msg.role === 'assistant'" class="mentor-info">
              <div class="avatar-container">
                <div class="avatar mentor-avatar">
                  <el-icon>
                    <User />
                  </el-icon>
                </div>
                <div class="username">
                  苏格拉底导师
                  <span v-if="msg.mode" class="mode-tag">[{{ msg.mode === 'slow' ? '慢思考' : '快思考' }}]</span>
                </div>
              </div>
              <div class="message-group">
                <div class="message-bubble mentor-bubble">
                  <pre>{{ msg.content }}</pre>
                </div>
              </div>
            </div>

            <!-- 用户消息 -->
            <div v-else class="user-info">
              <div class="avatar-container">
                <div class="avatar user-avatar">
                  <el-icon>
                    <User />
                  </el-icon>
                </div>
                <div class="username">学习者</div>
              </div>
              <div class="message-group">
                <div class="message-bubble user-bubble">
                  <pre>{{ msg.content }}</pre>
                </div>
              </div>
            </div>
          </div>
        </div>

      </el-main>

      <!-- 输入区域 -->
      <el-footer class="teach-footer">
        <div class="prompt-grid">
          <el-button v-for="(prompt, index) in prompts" :key="index" class="prompt-btn" @click="fillPrompt(prompt)">
            {{ prompt }}
          </el-button>
        </div>
        <div class="input-wrapper">
          <el-input v-model="userMessage" :placeholder="isSlowThinking ? '慢思考模式：请输入你的问题' : '快思考模式：请输入你的问题'" :rows="1"
            type="textarea" :disabled="isProcessing" @keyup.enter="sendMessage">
            <template #append>
              <el-button :loading="isProcessing" type="primary" @click="sendMessage" :disabled="!userMessage.trim()">
                {{ isProcessing ? '处理中...' : '发送' }}
              </el-button>
            </template>
          </el-input>
        </div>
      </el-footer>
    </el-container>
  </div>
</template>

<script lang="ts" setup>
import { nextTick, ref, onMounted } from 'vue';
import OpenAI from "openai";

// 状态管理
const isSlowThinking = ref(false);
const isProcessing = ref(false);
const processingInfo = ref('');
const conversationId = ref('');

// 原有状态
const isReviewing = ref(false);
const deanMessages = ref<Array<{ content: string, status: 'accept' | 'revise' }>>([]);

// 预设苏格拉底式提问模板
const prompts = ref([
  "这个结论的前提假设是什么？",
  "能否用其他方式解释这个现象？",
  "可以举个例子说明吗？",
]);

// OpenAI 配置（用于快思考模式）
const openai = new OpenAI({
  apiKey: "sk-a23ce2357e2c4432816c772fd7449498",
  baseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
  dangerouslyAllowBrowser: true
});

const userMessage = ref('');
const messages = ref([
  {
    role: 'assistant',
    content: '让我们通过提问来探索真理吧！\n请先分享你的观点或疑问，我将通过系列提问帮助你深化思考。',
    mode: 'fast'
  }
]);

// 后端API配置
const API_BASE_URL = 'http://localhost:5000/api';

onMounted(async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/conversation`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (response.ok) {
      const data = await response.json();
      if (data.success) {
        conversationId.value = data.conversation_id;
        console.log('会话创建成功:', conversationId.value);

        // 如果有知识状态摘要，显示给用户
        if (data.knowledge_summary) {
          messages.value.push({
            role: 'assistant',
            content: `系统知识状态摘要：${data.knowledge_summary}`,
            mode: 'system'
          });
        }
      }
    }
  } catch (error) {
    console.error('初始化会话失败:', error);
  }
});

// 教师prompt
const teacherSystemMessage = `"你是一名教师，采用苏格拉底教学法，不用考虑学生的情绪问题，专注于解题过程，通过提问引导学生自主发现答案。
核心规则：
1. 绝不直接提供完整答案。相反，通过提问促使学生分析问题。
2. 从开放式问题入手，帮助学生拆解问题。
3. 若学生答案不完整或错误，根据其回答提出后续问题——不给予直接提示或答案。
4. 确保每个问题在逻辑上承接上一个问题，构建连贯的推理链条。
5. 这些规则适用于所有学科（数学、科学、文学等）。
--------------------------";`;
// 思考模式切换
const onThinkingModeChange = async (value: boolean) => {
  console.log('思考模式切换:', value ? '慢思考' : '快思考');

  // 更新后端会话模式
  try {
    const response = await fetch(`${API_BASE_URL}/conversation/${conversationId.value}/mode`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        mode: value ? 'slow' : 'fast'
      })
    });

    if (response.ok) {
      const data = await response.json();
      console.log('模式更新成功:', data.mode);
    }
  } catch (error) {
    console.error('更新模式失败:', error);
  }

  // 添加模式切换提示消息
  messages.value.push({
    role: 'assistant',
    content: value
      ? '已切换到慢思考模式 \n我将使用完整的ICECoT智能体分析流程来处理你的问题，包括：\n• 情绪状态分析\n• 学习意图推断\n• 教学策略选择\n• 个性化响应生成\n\n这可能需要更多时间，但会提供更深入的分析和指导。'
      : '已切换到快思考模式 。',
    mode: value ? 'slow' : 'fast'
  });

  scrollToBottom();
};

// 填充预设提示词
const fillPrompt = (text: string) => {
  userMessage.value = text;
};

// 快思考模式处理
const handleFastThinking = async (tempMsg: string) => {
  try {

    const teacherResponse = await openai.chat.completions.create({
      model: "qwen-plus",
      messages: [
        { role: 'system', content: teacherSystemMessage },
        { role: 'user', content: tempMsg }
      ]
    });
    const finalReply = teacherResponse.choices[0].message.content;
    console.log("教师回复:", finalReply);


    messages.value.push({
      role: 'assistant',
      content: finalReply,
      mode: 'fast'
    });

  } catch (error) {
    console.error('快思考模式错误:', error);
    messages.value.push({
      role: 'assistant',
      content: `抱歉，处理过程中出现错误：${error.message}`,
      mode: 'fast'
    });
  }
};

// 慢思考模式处理
const handleSlowThinking = async (tempMsg: string) => {
  try {
    processingInfo.value = '正在启动ICECoT智能体分析流程...';

    const response = await fetch(`${API_BASE_URL}/chat`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        conversation_id: conversationId.value,
        message: tempMsg,
        mode: 'slow'
      })
    });

    if (!response.ok) {
      throw new Error(`API请求失败: ${response.status}`);
    }

    const data = await response.json();

    if (data.success) {
      processingInfo.value = data.processing_info || '';

      messages.value.push({
        role: 'assistant',
        content: data.response,
        mode: 'slow'
      });
    } else {
      throw new Error(data.error || '未知错误');
    }

  } catch (error) {
    console.error('慢思考模式错误:', error);
    messages.value.push({
      role: 'assistant',
      content: `抱歉，慢思考模式处理失败：${error.message}`,
      mode: 'slow'
    });
  } finally {
    processingInfo.value = '';
  }
};

// 发送消息主逻辑
const sendMessage = async () => {
  if (!userMessage.value.trim() || isProcessing.value) return;

  isProcessing.value = true;

  // 添加用户消息
  messages.value.push({
    role: 'user',
    content: userMessage.value,
    mode: ''
  });
  const tempMsg = userMessage.value;
  userMessage.value = '';

  try {
    if (isSlowThinking.value) {
      // 慢思考模式
      await handleSlowThinking(tempMsg);
    } else {
      // 快思考模式
      await handleFastThinking(tempMsg);
    }
  } catch (error) {
    console.error('发送消息错误:', error);
    messages.value.push({
      role: 'assistant',
      content: `系统错误：${error.message}`,
      mode: isSlowThinking.value ? 'slow' : 'fast'
    });
  } finally {
    isProcessing.value = false;
    scrollToBottom();
  }
};

// 滚动控制
const scrollToBottom = () => {
  nextTick(() => {
    const container = document.querySelector('.dialogue-box');
    if (container) {
      container.scrollTop = container.scrollHeight;
    }
  });
};
</script>

<style scoped lang="scss">
.socratic-teach {
  height: 85vh;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  padding: 20px;
}

.teach-container {
  max-width: 1400px;
  height: 80vh;
  margin: 0 auto;
  background: white;
  border-radius: 12px;
  overflow: hidden;
}

.teach-header {
  padding: 0rem;
  background: #5f8bb6;

  .header-content {
    text-align: center;
    color: white;
    padding: 1px;
    position: relative;

    h1 {
      font-size: 35px;
      margin: 0 0 10px 0;
      letter-spacing: 4px;
    }

    // 思考模式切换区域
    .thinking-mode-switch {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 8px;
      margin-top: 10px;

      .mode-description {
        font-size: 14px;
        opacity: 0.9;
        min-height: 20px;

        .slow-mode-desc {
          color: #ffd700;
        }

        .fast-mode-desc {
          color: #90ee90;
        }
      }
    }
  }
}

.processing-info {
  margin: 10px;
}

.mode-tag {
  font-size: 12px;
  color: #666;
  font-weight: normal;
  margin-left: 5px;
  padding: 2px 6px;
  background: #f0f0f0;
  border-radius: 3px;
}

.avatar-container {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;

  .mentor-info & {
    margin-right: 12px;
  }

  .user-info & {
    flex-direction: row;
    margin-left: 12px;
  }

  .username {
    font-size: 25px;
    font-weight: 600;
    color: #1e67b0;
    letter-spacing: 0.5px;
  }
}

.message-container {
  margin: 18px 0;
  display: flex;

  &.user-message {
    justify-content: flex-end;
  }

  &.mentor-message {
    justify-content: flex-start;
  }
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background-size: cover;

  .el-icon {
    font-size: 24px;
    color: white;
  }

  &.mentor-avatar {
    background-image: url('/images/苏格拉底.png');
  }

  &.user-avatar {
    background-image: url('/images/学生.jpg');
  }
}

.message-bubble {
  padding: 12px 16px;
  border-radius: 8px;
  position: relative;
  line-height: 1.6;

  pre {
    margin: 0;
    white-space: pre-wrap;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  }

  &.mentor-bubble {
    background: #F1F5FF;
    border: 1px solid #D3E0FD;

    &::before {
      content: "";
      position: absolute;
      left: -8px;
      top: 12px;
      border: 8px solid transparent;
      border-right-color: #D3E0FD;
    }
  }

  &.user-bubble {
    background: #E8F5E9;
    border: 1px solid #C8E6C9;

    &::before {
      content: "";
      position: absolute;
      right: -8px;
      top: 12px;
      border: 8px solid transparent;
      border-left-color: #C8E6C9;
    }
  }
}

.dean-panel {
  margin: 10px;
  border: 1px solid #eee;
  border-radius: 8px;
  background: #f8f9fa;
}

.dean-alert {
  padding: 12px;
  display: flex;
  align-items: center;
  gap: 8px;

  &.accept {
    color: #67c23a;
    background: #f0f9eb;
  }

  &.revise {
    color: #e6a23c;
    background: #fdf6ec;
  }

  pre {
    margin: 0;
    font-size: 0.9em;
  }
}

.dialogue-box {
  height: 700px;
  padding: 1rem;
  overflow-y: auto;
  background: #f8f9fa;

  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-track {
    background: #f1f1f1;
  }

  &::-webkit-scrollbar-thumb {
    background: #c1c1c1;
    border-radius: 3px;
  }

  &::-webkit-scrollbar-thumb:hover {
    background: #a8a8a8;
  }
}

.prompt-grid {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 20px;
  padding: 5px 20px;
  background: white;
  border-top: 1px solid #eee;
  max-height: 60px;
  overflow: hidden;

  .prompt-btn {
    font-size: 1rem;
    padding: 6px 12px;
    background: #f8f9fa;
    border: 1px solid #dee2e6;
    border-radius: 6px;
    color: #2c3e50;
    transition: all 0.2s;
    margin: 0 !important;
    white-space: nowrap;
    flex: 0 0 auto;

    &:hover {
      background: #99bee2;
      transform: translateY(-1px);
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }
  }
}

.teach-footer {
  padding: 0px 20px;
  background: white;
  border-top: 1px solid #eee;
  height: 120px;

  .input-wrapper {
    :deep(.el-textarea__inner) {
      border-radius: 10px;
      padding: 12px 16px;
      font-size: 20px;
      border: 1px solid #4ca0e1;

      &:focus {
        font-size: 20px;
        border-color: #409eff;
        box-shadow: 0 0 8px rgba(64, 158, 255, 0.3);
      }

      &:disabled {
        background-color: #f5f7fa;
        color: #909399;
      }
    }

    :deep(.el-input-group__append) {
      .el-button {
        border-radius: 0 12px 12px 0;
      }
    }
  }
}

@media (max-width: 768px) {
  .avatar {
    width: 36px;
    height: 36px;

    .el-icon {
      font-size: 20px;
    }
  }

  .message-bubble {
    padding: 10px 14px;
    font-size: 14px;
  }

  .thinking-mode-switch {
    .mode-description {
      font-size: 12px;
    }
  }

  .prompt-grid {
    gap: 6px;
    padding: 8px 15px;
    max-height: 50px;

    .prompt-btn {
      font-size: 0.7rem;
      padding: 4px 8px;
    }
  }

  .dialogue-box {
    height: calc(70vh - 80px);
  }
}
</style>