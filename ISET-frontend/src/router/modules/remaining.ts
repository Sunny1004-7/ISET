const Layout = () => import("@/layout/index.vue");

export default [
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/login/index.vue"),
    meta: {
      title: "登录",
      showLink: false,
      rank: 101
    }
  },
  {
    path: "/redirect",
    component: Layout,
    meta: {
      icon: "homeFilled",
      title: "首页",
      showLink: false,
      rank: 104
    },
    children: [
      {
        path: "/redirect/:path(.*)",
        name: "Redirect",
        component: () => import("@/layout/redirect.vue")
      }
    ]
  },
  {
    path: "/error/403",
    name: "403",
    component: () => import("@/views/error/403.vue"),
    meta: {
      showLink: false,
      title: "403"
    }
  },
  {
    path: "/error/404",
    name: "404",
    component: () => import("@/views/error/404.vue"),
    meta: {
      showLink: false,
      title: "404"
    }
  },
  {
    path: "/error/500",
    name: "500",
    component: () => import("@/views/error/500.vue"),
    meta: {
      showLink: false,
      title: "500"
    }
  }
] as Array<RouteConfigsTable>;
