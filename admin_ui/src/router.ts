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
    meta: {needLogin: true},
    component: () => import("@/views/home/home.vue"),
    children: [
      {
        path: "/user/modify",
        name: "userModify",
        meta: {needLogin: true},
        component: () => import("@/views/user/modify.vue")
      },
      {
        path: "/user/create",
        name: "userCreate",
        meta: {needLogin: true},
        component: () => import("@/views/user/create.vue")
      },
      {
        path: "/user/list",
        name: "userList",
        meta: {needLogin: true},
        component: () => import("@/views/user/list.vue")
      }
    ]
  },
  {
    path: "/404",
    name: "notFound",
    redirect: {name: "login"}
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
  if (!to.meta || !to.meta.needLogin) { // default keyword is false
    next();
    return;
  }

  if (sessionStorage.getItem("auth")) {
    next();
    return;
  } else {
    next({path: "/login"});
    return;
  }
})
