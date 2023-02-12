import {
  createRouter,
  createWebHashHistory,
  NavigationGuardNext,
  RouteLocationNormalized,
  RouteLocationNormalizedLoaded,
  RouteRecordRaw
} from "vue-router"

let routes: RouteRecordRaw[] = [
  {
    path: "/",
    name: "home",
    component: () => import("@/views/home/home.vue"),
  },
  {
    path: "/user",
    name: "user",
    meta: { requireLogin: true },
    component: () => import("@/views/user/index.vue"),
    children: [
      {
        path: "modify-info",
        name: "userModifyInfo",
        component: () => import("@/views/user/modify_info.vue"),
      },
      {
        path: "list",
        name: "userList",
        component: () => import("@/views/user/list.vue"),
      },
      {
        path: "create",
        name: "userCreate",
        component: () => import("@/views/user/create.vue"),
      },
    ]
  },
  {
    path: "/cloud-file",
    name: "cloudFile",
    meta: { requireLogin: true },
    component: () => import("@/views/cloud_file/index.vue"),
    children: [
      {
        path: "list",
        name: "cloudFileList",
        component: () => import("@/views/cloud_file/list.vue"),
      },
      {
        path: "upload",
        name: "cloudFileUpload",
        component: () => import("@/views/cloud_file/upload.vue"),
      },
    ]
  },
  {
    path: "/404", // 考虑后续可能编写404页面，此处预留路由
    name: "notFound",
    redirect: { name: "home" }
  },
  {
    path: "/:pathMatch(.*)*", // 将匹配所有内容并将其放在 `$route.params.pathMatch` 下
    redirect: { name: "notFound" }
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes: routes,
})

router.beforeEach((to: RouteLocationNormalized, from: RouteLocationNormalizedLoaded, next: NavigationGuardNext) => {
  if (!to.meta || !to.meta.requireLogin) { // default keyword is false
    next();
    return;
  }

  if (sessionStorage.getItem("auth") === import.meta.env.Vite_axios_source_sign as string) {
    next();
    return;
  } else {
    next({ path: "/" });
    return;
  }
})

export default router
