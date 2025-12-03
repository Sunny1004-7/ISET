<template>
  <div class="container">
    <div class="header">
      <h1>在线课程平台</h1>
      <p>发现知识，启发创新</p>
    </div>

    <div class="search-container">
      <input type="text" v-model="searchQuery" placeholder="搜索您感兴趣的课程..." class="search-input">
      <div class="search-icon">🔍</div>
    </div>

    <transition name="fade">
      <div v-if="filteredCourses.length > 0" class="courses-grid">
        <div v-for="course in filteredCourses" :key="course.id" class="course-card" @click="viewCourse(course)">
          <div class="course-image" :style="{ backgroundImage: `url(${course.image})` }"></div>
          <div class="course-content">
            <div class="course-title">{{ course.name }}</div>
            <div class="course-description">{{ course.description }}</div>
          </div>
        </div>
      </div>
      <div v-else class="no-results">
        没有找到匹配的课程，试试其他关键词吧！
      </div>
    </transition>

    <button @click="goToRecommendPage" class="recommend-btn">
      🎯 智能推荐课程
    </button>

    <!-- 加载动画 -->
    <div v-if="isLoading" class="loading-overlay">
      <div class="loading-container">
        <div class="spinner"></div>
        <div class="loading-text">智能推荐生成中...</div>
      </div>
    </div>

    <!-- 简化的推荐结果弹窗卡片 -->
    <div v-if="showRecommendCard" class="modal-overlay" @click="closeModal">
      <div class="recommend-card" @click.stop>
        <div class="card-header">
          <h2>推荐分析结果</h2>
          <button class="close-btn" @click="closeModal">✕</button>
        </div>

        <div class="card-content">
          <div class="info-section">
            <div class="section-title">知识水平评估</div>
            <div class="content-text">
              {{ recommendData.知识水平评估 }}
            </div>
          </div>

          <div class="info-section">
            <div class="section-title">推荐原因</div>
            <div class="content-text">
              {{ recommendData.推荐原因 }}
            </div>
          </div>
        </div>

        <div class="card-footer">
          <button class="cancel-btn" @click="closeModal">稍后查看</button>
          <button class="confirm-btn" @click="confirmAndNavigate">查看详细推荐</button>
        </div>
      </div>
    </div>

    <router-view></router-view>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
const router = useRouter()

// 响应式数据
const searchQuery = ref('')
const USER_ID = 'U_10258448'
const isLoading = ref(false)
const showRecommendCard = ref(false)
const recommendData = ref({})

const courses = ref([
  {
    id: 1,
    name: "Vue.js 全栈开发",
    image: "/images/vue.png",
    description: "从零基础到项目实战，掌握Vue.js前端开发技术栈，包括组件化开发、路由管理、状态管理等核心概念。"
  },
  {
    id: 2,
    name: "React 现代开发",
    image: "/images/react.png",
    description: "学习React Hooks、Context API、路由等现代React开发技术，构建高性能的单页面应用程序。"
  },
  {
    id: 3,
    name: "Node.js 后端开发",
    image: "/images/node.png",
    description: "使用Node.js构建RESTful API，学习Express框架、数据库集成、身份验证等后端开发技能。"
  },
  {
    id: 4,
    name: "Python 数据科学",
    image: "/images/python.png",
    description: "掌握Python在数据分析、机器学习领域的应用，包括pandas、numpy、scikit-learn等重要库。"
  },
  {
    id: 5,
    name: "UI/UX 设计基础",
    image: "/images/uiux.png",
    description: "学习用户界面和用户体验设计原理，掌握Figma、Sketch等设计工具的使用方法。"
  },
  {
    id: 6,
    name: "移动端开发",
    image: "/images/mobile.png",
    description: "React Native和Flutter跨平台移动应用开发，一套代码适配iOS和Android平台。"
  },
  {
    id: 7,
    name: "数据结构与算法",
    image: "/images/datastructure.png",
    description: "系统学习线性表、栈、队列、树、图等基础数据结构，掌握排序、查找、动态规划等核心算法设计与分析方法，培养高效解决复杂问题的能力。"
  },
  {
    id: 8,
    name: "计算机网络",
    image: "/images/computernetwork.png",
    description: "深入理解TCP/IP协议栈、网络分层模型，学习路由算法、拥塞控制、数据链路层协议等核心知识，掌握网络编程与网络故障排查的基本技能。"
  }
])

// 计算属性
const filteredCourses = computed(() => {
  if (!searchQuery.value) return courses.value
  return courses.value.filter(course =>
    course.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
    course.description.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

const loadRecommendData = async () => {
  try {
    const response = await fetch('/test.txt');
    if (!response.ok) {
      throw new Error(`文件未找到：${response.status}`);
    }
    const testText = await response.text();
    const data = JSON.parse(testText);
    // 只提取需要的字段
    return {
      知识水平评估: data.知识水平评估 || '暂无评估信息',
      推荐原因: data.推荐原因 || '暂无推荐原因'
    }
  } catch (error) {
    console.error('读取推荐数据失败:', error)
    return {
      知识水平评估: '数据加载失败，请稍后重试',
      推荐原因: '数据加载失败，请稍后重试'
    }
  }
}

// 方法
const viewCourse = (course) => {
  alert(`点击了课程: ${course.name}`)
}

// 关闭弹窗
const closeModal = () => {
  showRecommendCard.value = false
}

// 确认并跳转
const confirmAndNavigate = () => {
  showRecommendCard.value = false
  router.push('/courseRec/recommend')
}

// 跳转到推荐页面
const goToRecommendPage = async () => {
  try {
    isLoading.value = true

    await axios.post('http://localhost:7788/recommend', {
      user_id: USER_ID,
    })

    const localData = await loadRecommendData()
    recommendData.value = localData

    showRecommendCard.value = true

  } catch (error) {
    console.error('推荐失败:', error)
    alert(`推荐失败：${error.message || '请稍后重试'}`)
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.container {
  width: 100%;
  height: 100%;
  margin: 0 auto;
  padding: 20px;
  background: linear-gradient(135deg, #bfc6e9 0%, #7dc2e8 100%);
  min-height: 100vh;
}

.header {
  text-align: center;
  margin-bottom: 40px;
  color: white;
}

.header h1 {
  font-size: 2.5rem;
  margin-bottom: 10px;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
}

.header p {
  font-size: 1.1rem;
  opacity: 0.9;
}

.search-container {
  max-width: 600px;
  margin: 0 auto 40px auto;
  position: relative;
}

.search-input {
  width: 100%;
  padding: 15px 50px 15px 20px;
  border: none;
  border-radius: 25px;
  font-size: 16px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
  outline: none;
  transition: box-shadow 0.3s ease;
}

.search-input:focus {
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
}

.search-icon {
  position: absolute;
  right: 20px;
  top: 50%;
  transform: translateY(-50%);
  color: #666;
  font-size: 18px;
}

.courses-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 25px;
  margin-bottom: 40px;
}

.course-card {
  background: white;
  border-radius: 15px;
  overflow: hidden;
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  cursor: pointer;
}

.course-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.15);
}

.course-image {
  width: 100%;
  height: 150px;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
}

.course-content {
  padding: 20px;
}

.course-title {
  font-size: 1.2rem;
  font-weight: bold;
  color: #333;
  margin-bottom: 10px;
}

.course-description {
  color: #666;
  line-height: 1.5;
  font-size: 14px;
}

.recommend-btn {
  display: block;
  margin: 0 auto;
  padding: 15px 40px;
  background: linear-gradient(45deg, #ff6b6b, #ff8e8e);
  color: white;
  border: none;
  border-radius: 25px;
  font-size: 18px;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.3s ease;
  text-decoration: none;
  box-shadow: 0 4px 15px rgba(255, 107, 107, 0.3);
}

.recommend-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(255, 107, 107, 0.4);
}

/* 加载动画样式 */
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(3px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 9999;
  transition: opacity 0.3s ease;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.spinner {
  width: 50px;
  height: 50px;
  border: 5px solid #f0f0f0;
  border-top: 5px solid #ff6b6b;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }

  100% {
    transform: rotate(360deg);
  }
}

.loading-text {
  font-size: 18px;
  color: #333;
  font-weight: 500;
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {

  0%,
  100% {
    opacity: 1;
  }

  50% {
    opacity: 0.7;
  }
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 10000;
}

.recommend-card {
  background: white;
  border-radius: 15px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  max-width: 650px;
  width: 90%;
  max-height: 70vh;
  overflow-y: auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 25px;
  border-bottom: 1px solid #eee;
  background: #f8f9fa;
}

.card-header h2 {
  margin: 0;
  color: #333;
  font-size: 1.3rem;
  font-weight: 600;
}

.close-btn {
  background: #f0f0f0;
  border: none;
  border-radius: 50%;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #666;
  font-size: 14px;
}

.close-btn:hover {
  background: #e0e0e0;
}

.card-content {
  padding: 25px;
}

.info-section {
  margin-bottom: 20px;
}

.section-title {
  font-size: 1rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 10px;
  padding-bottom: 5px;
  border-bottom: 2px solid #e9ecef;
}

.content-text {
  color: #555;
  line-height: 1.6;
  font-size: 14px;
  background: #f8f9fa;
  padding: 15px;
  border-radius: 8px;
  border-left: 3px solid #007bff;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  padding: 15px 25px 20px;
  border-top: 1px solid #eee;
  gap: 10px;
}

.cancel-btn,
.confirm-btn {
  padding: 10px 20px;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  flex: 1;
  font-size: 14px;
}

.cancel-btn {
  background: #f0f0f0;
  color: #666;
  border: 1px solid #ddd;
}

.cancel-btn:hover {
  background: #e0e0e0;
}

.confirm-btn {
  background: #007bff;
  color: white;
  border: none;
}

.confirm-btn:hover {
  background: #0056b3;
}

.no-results {
  text-align: center;
  color: white;
  font-size: 1.2rem;
  margin: 40px 0;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@media (max-width: 768px) {
  .courses-grid {
    grid-template-columns: 1fr;
  }

  .header h1 {
    font-size: 2rem;
  }

  .container {
    padding: 15px;
  }

  .recommend-card {
    width: 95%;
    margin: 10px;
  }

  .card-header,
  .card-content,
  .card-footer {
    padding-left: 20px;
    padding-right: 20px;
  }

  .card-footer {
    flex-direction: column;
  }
}
</style>