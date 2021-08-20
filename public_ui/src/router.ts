import Vue from "vue"
import VueRouter, {RouteConfig} from "vue-router"

Vue.use(VueRouter)

const routes: Array<RouteConfig> = [
  {
    path: "/",
    name: "",
    meta: {needLogin: false},
    component: () => import("@/views/home/home.vue"),
    children: [
      {
        path: "",
        name: "home",
        meta: {needLogin: false},
        component: () => import("@/views/home/content.vue")
      },
      {
        path: "files",
        name: "files",
        meta: {needLogin: true},
        component: () => import("@/views/files/files.vue")
      }
    ]
  },
  {
    path: "/404",
    name: "notFound",
    redirect: {name: "home"}
  },
  {
    path: "*",
    redirect: {name: "notFound"}
  }
]

const router = new VueRouter({
  routes
})

export default router

router.beforeEach((to, from, next) => {
  if (!to.meta!.needLogin) {
    next();
    return;
  }

  if (sessionStorage.getItem("auth")) {
    next();
    return;
  } else {
    next({path: "/"});
    return;
  }
})
