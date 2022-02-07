import Vue from 'vue'
import VueRouter, { RouteConfig } from 'vue-router'

Vue.use(VueRouter)

const routes: Array<RouteConfig> = [
  {
    path: "/login",
    name: "login",
    component: () => import("@/views/home/login.vue")
  },
  {
    path: "/",
    name: "home",
    meta: { needLogin: true },
    component: () => import("@/views/home/home.vue"),
    children: [
      {
        path: "user/modify",
        name: "userModify",
        meta: { needLogin: true },
        component: () => import("@/views/user/modify.vue")
      },
      {
        path: "user/create",
        name: "userCreate",
        meta: { needLogin: true },
        component: () => import("@/views/user/create.vue")
      },
      {
        path: "user/list",
        name: "userList",
        meta: { needLogin: true },
        component: () => import("@/views/user/list.vue")
      },
      {
        path: "cloud-file/upload",
        name: "cloudFileUpload",
        meta: { needLogin: true },
        component: () => import("@/views/cloud_file/upload.vue")
      },
      {
        path: "cloud-file/list-by-uploader",
        name: "cloudFileListByUploader",
        meta: { needLogin: true },
        component: () => import("@/views/cloud_file/list_by_uploader.vue")
      },
      {
        path: "cloud-file/list-public",
        name: "cloudFileListPublic",
        meta: { needLogin: true },
        component: () => import("@/views/cloud_file/list_public.vue")
      },
      {
        path: "note/create",
        name: "noteCreate",
        meta: { needLogin: true },
        component: () => import("@/views/note/create.vue")
      },
      {
        path: "note/list-by-writer",
        name: "noteListByWriter",
        meta: { needLogin: true },
        component: () => import("@/views/note/list_by_writer.vue")
      },
      {
        path: "note/list-public",
        name: "noteListPublic",
        meta: { needLogin: true },
        component: () => import("@/views/note/list_public.vue")
      },
      {
        path: "task/create",
        name: "taskCreate",
        meta: { needLogin: true },
        component: () => import("@/views/task/create.vue")
      },
      {
        path: "task/list",
        name: "taskList",
        meta: { needLogin: true },
        component: () => import("@/views/task/list.vue")
      },
    ]
  },
  {
    path: "/404",
    name: "notFound",
    redirect: { name: "login" }
  },
  {
    path: "*",
    redirect: { name: "notFound" }
  }
]

const router = new VueRouter({
  routes
})

export default router

router.beforeEach((to, from, next) => {
  if (!to.meta || !to.meta.needLogin) { // default keyword is false
    next();
    return;
  }

  if (sessionStorage.getItem("auth")) {
    next();
    return;
  } else {
    next({ path: "/login" });
    return;
  }
})
