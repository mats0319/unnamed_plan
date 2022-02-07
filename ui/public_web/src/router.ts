import Vue from "vue"
import VueRouter, { RouteConfig } from "vue-router"

Vue.use(VueRouter)

const routes: Array<RouteConfig> = [
  {
    path: "/",
    name: "",
    component: () => import("@/views/home/home.vue"),
    children: [
      {
        path: "",
        name: "home",
        component: () => import("@/views/home/content.vue")
      },
      {
        path: "cloud-file/list-by-uploader",
        name: "cloudFileListByUploader",
        meta: { needLogin: true },
        component: () => import("@/views/cloud_file/list.vue")
      },
      {
        path: "cloud-file/list-public",
        name: "cloudFileListPublic",
        meta: { needLogin: true },
        component: () => import("@/views/cloud_file/list.vue")
      },
      {
        path: "note/list-by-writer",
        name: "noteListByWriter",
        meta: { needLogin: true },
        component: () => import("@/views/note/list.vue")
      },
      {
        path: "note/list-public",
        name: "noteListPublic",
        meta: { needLogin: true },
        component: () => import("@/views/note/list.vue")
      }
    ]
  },
  {
    path: "/404",
    name: "notFound",
    redirect: { name: "home" }
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
    next({ path: "/" });
    return;
  }
})
