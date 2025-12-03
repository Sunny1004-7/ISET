export default {
  path: "/Socratic",
  name: "Socratic",
  component: () => import("@/views/Socratic/index.vue"),
  meta: {
    icon: "homeFilled",
    title: "苏格拉底式教学",
    rank: 4
  }
} as RouteConfigsTable;