<template>
  <div class="container">
    <div class="page-layout">
      <!-- 左侧推荐课程 -->
      <div class="left-section">
        <h2 class="section-title">推荐课程</h2>
        <div class="course-list">
          <div v-for="course in recommendedCourses" :key="`rec-${course.id}`" class="course-item"
            @click="enrollCourse(course)">
            <div class="course-image-wrapper">
              <img :src="course.image" :alt="course.name" class="course-image" />
            </div>
            <div class="course-content">
              <div class="course-header">
                <h3 class="course-title">{{ course.name }}</h3>
                <div class="course-stats">
                  <i class="eye-icon">👁</i>
                  <span class="view-count">{{ course.views }}</span>
                  <span class="course-level">{{ course.progress }}</span>
                </div>
              </div>
              <div class="course-tags">
                <span v-for="tag in course.tags" :key="tag" class="course-tag" :class="getTagClass(tag)">
                  {{ tag }}
                </span>
              </div>
              <div class="course-description">
                <p class="course-info">开课学校:{{ course.school }} 教师:{{ course.teacher }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧历史课程 -->
      <div class="right-section">
        <h2 class="section-title">历史课程</h2>

        <!-- 搜索框 -->
        <div class="search-container">
          <input type="text" v-model="historySearchQuery" placeholder="搜索历史课程..." class="search-input">
          <div class="search-icon">🔍</div>
        </div>

        <div class="history-list">
          <div v-for="course in filteredHistoryCourses" :key="`hist-${course.id}`" class="history-item"
            @click="continueCourse(course)">
            <div class="history-header">
              <h3 class="history-title">{{ course.name }}</h3>
              <div class="history-image">
                <img :src="course.image" :alt="course.name" />
              </div>
            </div>
            <div class="history-content">
              <p class="history-category">课程类别：{{ course.category }}</p>
              <div class="history-status">
                <span class="status-badge" :class="getStatusClass(course.status)">
                  {{ getStatusText(course.status) }}
                </span>
                <span class="progress-text">{{ course.progress }}</span>
              </div>
            </div>
          </div>

          <!-- 无搜索结果提示 -->
          <div v-if="filteredHistoryCourses.length === 0 && historySearchQuery" class="no-results">
            <div class="no-results-icon">📚</div>
            <p class="no-results-text">没有找到匹配的历史课程</p>
            <p class="no-results-hint">试试其他关键词吧</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import rec_course_data from "@/data/rec_course_data.json";
import his_course_data from "@/data/his_course_data.json";

const router = useRouter()

const historySearchQuery = ref('')

// 推荐课程数据
const recommendedCourses = ref(rec_course_data)

// 历史课程数据
const historyCourses = ref(his_course_data)

// 计算属性：过滤历史课程
const filteredHistoryCourses = computed(() => {
  if (!historySearchQuery.value) return historyCourses.value
  return historyCourses.value.filter(course =>
    course.name.toLowerCase().includes(historySearchQuery.value.toLowerCase()) ||
    course.category.toLowerCase().includes(historySearchQuery.value.toLowerCase())
  )
})

const enrollCourse = (course) => {
  alert(`准备学习课程: ${course.name}`)
}

const continueCourse = (course) => {
  if (course.status === 'completed') {
    alert(`重新学习课程: ${course.name}`)
  } else if (course.status === 'in-progress') {
    alert(`继续学习课程: ${course.name}`)
  } else {
    alert(`开始学习课程: ${course.name}`)
  }
}

const getTagClass = (tag) => {
  const tagClasses = {
    '数据科学': 'tag-primary',
    '机器学习': 'tag-primary',
    '深度学习': 'tag-secondary',
    'TensorFlow': 'tag-secondary',
    '大数据': 'tag-accent',
    '集群': 'tag-accent',
    '网络安全': 'tag-primary',
    '防护': 'tag-primary',
    '云计算': 'tag-secondary',
    '架构设计': 'tag-secondary',
    '算法竞赛': 'tag-accent',
    '编程': 'tag-accent'
  }
  return tagClasses[tag] || 'tag-primary'
}

const getStatusClass = (status) => {
  return {
    'completed': 'status-completed',
    'in-progress': 'status-progress',
    'saved': 'status-saved'
  }[status]
}

const getStatusText = (status) => {
  return {
    'completed': '已完成',
    'in-progress': '学习中',
    'saved': '已收藏'
  }[status]
}
</script>

<style scoped>
:root {
  --primary-color: #6366f1;
  --secondary-color: #64748b;
  --primary-light: #e0e7ff;
  --secondary-light: #f1f5f9;
  --accent-light: #a4a4a4;
  --primary-dark: #4f46e5;
  --secondary-dark: #475569;
  --accent-dark: #059669;
}

.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
  background-color: #fafafa;
  min-height: 100vh;
}

.page-layout {
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 20px;
}

.section-title {
  font-size: 24px;
  font-weight: 600;
  color: var(--secondary-dark);
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 2px solid var(--secondary-light);
}

.left-section {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid var(--secondary-light);
}

.course-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.course-item {
  display: flex;
  padding: 20px;
  border: 1px solid var(--secondary-light);
  border-radius: 10px;
  transition: all 0.3s ease;
  cursor: pointer;
  background: #ffffff;
}

.course-item:hover {
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.1);
  transform: translateY(-2px);
  border-color: var(--primary-color);
}

.course-image-wrapper {
  margin-right: 16px;
  flex-shrink: 0;
}

.course-image {
  width: 120px;
  height: 80px;
  border-radius: 8px;
  object-fit: cover;
  border: 1px solid var(--secondary-light);
}

.course-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.course-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.course-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--secondary-dark);
  margin: 0;
  flex: 1;
  margin-right: 12px;
  line-height: 1.4;
}

.course-stats {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--secondary-color);
  font-size: 14px;
}

.eye-icon {
  font-size: 16px;
  color: var(--secondary-color);
}

.view-count {
  font-weight: 500;
  color: var(--secondary-color);
}

.course-level {
  background: var(--primary-light);
  color: var(--primary-dark);
  padding: 3px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.course-tags {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.course-tag {
  padding: 4px 12px;
  border-radius: 16px;
  font-size: 12px;
  font-weight: 500;
  color: white;
  transition: all 0.2s ease;
}

.tag-primary {
  background: var(--primary-color);
}

.tag-primary:hover {
  background: var(--primary-dark);
}

.tag-secondary {
  background: var(--secondary-color);
}

.tag-secondary:hover {
  background: var(--secondary-dark);
}

.tag-accent {
  background: var(--accent-color);
}

.tag-accent:hover {
  background: var(--accent-dark);
}

.course-info {
  color: var(--secondary-color);
  font-size: 14px;
  margin: 0;
  line-height: 1.4;
}

.right-section {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid var(--secondary-light);
}

.search-container {
  position: relative;
  margin-bottom: 20px;
}

.search-input {
  width: 100%;
  padding: 12px 40px 12px 16px;
  border: 1px solid var(--secondary-light);
  border-radius: 8px;
  font-size: 14px;
  outline: none;
  transition: all 0.3s ease;
  background: #ffffff;
}

.search-input:focus {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.search-input::placeholder {
  color: var(--secondary-color);
}

.search-icon {
  position: absolute;
  right: 14px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--secondary-color);
  font-size: 16px;
  pointer-events: none;
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.history-item {
  border: 1px solid var(--secondary-light);
  border-radius: 10px;
  padding: 18px;
  transition: all 0.3s ease;
  cursor: pointer;
  background: #ffffff;
}

.history-item:hover {
  box-shadow: 0 4px 12px rgba(100, 116, 139, 0.1);
  transform: translateY(-2px);
  border-color: var(--secondary-color);
}

.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.history-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--secondary-dark);
  margin: 0;
  flex: 1;
  line-height: 1.4;
}

.history-image {
  margin-left: 12px;
}

.history-image img {
  width: 80px;
  height: 60px;
  border-radius: 6px;
  object-fit: cover;
  border: 1px solid var(--secondary-light);
}

.history-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.history-category {
  color: var(--secondary-color);
  font-size: 13px;
  margin: 0;
}

.history-status {
  display: flex;
  align-items: center;
  gap: 10px;
}

.status-badge {
  padding: 3px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-completed {
  background: var(--accent-light);
  color: var(--accent-dark);
}

.status-progress {
  background: var(--primary-light);
  color: var(--primary-dark);
}

.status-saved {
  background: var(--secondary-light);
  color: var(--secondary-dark);
}

.progress-text {
  color: var(--secondary-color);
  font-size: 12px;
}

.no-results {
  text-align: center;
  padding: 40px 20px;
  color: var(--secondary-color);
}

.no-results-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.no-results-text {
  font-size: 16px;
  font-weight: 500;
  margin: 0 0 8px 0;
  color: var(--secondary-dark);
}

.no-results-hint {
  font-size: 14px;
  margin: 0;
  opacity: 0.8;
}

@media (max-width: 1024px) {
  .page-layout {
    grid-template-columns: 1fr;
  }

  .right-section {
    order: -1;
  }
}

@media (max-width: 768px) {
  .container {
    padding: 16px;
  }

  .left-section,
  .right-section {
    padding: 20px;
  }

  .search-container {
    margin-bottom: 16px;
  }

  .search-input {
    padding: 10px 36px 10px 14px;
    font-size: 16px;

  }

  .course-item {
    flex-direction: column;
    padding: 16px;
  }

  .course-image-wrapper {
    margin-right: 0;
    margin-bottom: 12px;
    align-self: center;
  }

  .course-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .course-stats {
    margin-top: 8px;
  }

  .history-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .history-image {
    margin-left: 0;
    margin-top: 10px;
    align-self: center;
  }

  .course-tags {
    gap: 6px;
  }
}
</style>
