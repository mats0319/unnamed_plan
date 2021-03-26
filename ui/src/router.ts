import Vue from "vue"
import VueRouter, { RouteConfig } from "vue-router"

Vue.use(VueRouter)

const routes: Array<RouteConfig> = [
  {
    path: "/",
    name: "home",
    component: () => import("@/views/home/home.vue")
  },
  {
    path: "/404",
    name: "notFound",
    component: () => import("@/views/home/home.vue")
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
