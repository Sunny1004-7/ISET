export default {
  path: "/exercises",
  name: "exercises",
  redirect: "/exercises",
  meta: {
    icon: "homeFilled",
    title: "大模型自动出题",
    rank: 3
  },
  children: [
    {
      path: "/exercises/exercisesIndex-stu",
      name: "exercisesIndex-stu",
      component: () => import("@/views/exercises/exercisesIndex-stu.vue"),
      meta: {
        title: "AI出题-学生"
      }
    },
    {
      path: "/exercises/exercisesIndex-tea",
      name: "exercisesIndex-tea",
      component: () => import("@/views/exercises/exercisesIndex-tea.vue"),
      meta: {
        title: "AI出题-教师"
      }
    },
  ]
} as RouteConfigsTable;
