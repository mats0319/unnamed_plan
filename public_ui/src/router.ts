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
        component: () => import("@/views/cloud_file/list_by_uploader.vue")
      },
      {
        path: "cloud-file/list-public",
        name: "cloudFileListPublic",
        meta: { needLogin: true },
        component: () => import("@/views/cloud_file/list_public.vue")
      },
      {
        path: "thinking-note/list-by-writer",
        name: "thinkingNoteListByWriter",
        meta: { needLogin: true },
        component: () => import("@/views/thinking_note/list_by_writer.vue")
      },
      {
        path: "thinking-note/list-public",
        name: "thinkingNoteListPublic",
        meta: { needLogin: true },
        component: () => import("@/views/thinking_note/list_public.vue")
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
