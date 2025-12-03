export default {
  path: "/courseRec",
  name: "courseRec",
  component: () => import("@/views/courseRec/index.vue"),
  children: [
    {
      path: "/courseRec/recommend",
      name: "courseRecRecommend",
      component: () => import("@/views/courseRec/Recommend.vue"),
      meta: {
        title: "推荐课程详情"
      }
    }
  ],
  meta: {
    icon: "homeFilled",
    title: "课程推荐",
    rank: 2
  }
} as RouteConfigsTable;