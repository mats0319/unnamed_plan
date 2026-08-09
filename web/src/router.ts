import { createRouter, createWebHashHistory, RouteRecordRaw } from "vue-router"
import { useUserStore } from "@/pinia/user.ts"

const routes: Array<RouteRecordRaw> = [
    {
        path: "/",
        name: "home",
        component: () => import("@/views/home.vue"),
    },
    {
        path: "/note",
        name: "note",
        component: () => import("@/views/note.vue"),
    },
    {
        path: "/personal-center",
        name: "personalCenter",
        meta: { requireLogin: true },
        component: () => import("@/views/personal_center/frame.vue"),
        children: [
            {
                path: "",
                name: "pDefault",
                component: () => import("@/views/personal_center/default.vue"),
            },
            {
                path: "list-user",
                name: "pListUser",
                component: () => import("@/views/personal_center/list_user.vue"),
            },
            {
                path: "modify-user",
                name: "pModifyUser",
                component: () => import("@/views/personal_center/modify_user.vue"),
            },
            {
                path: "set-mfa-status",
                name: "pSetMFAStatus",
                component: () => import("@/views/personal_center/set_mfa_status.vue"),
            },
            {
                path: "note",
                name: "pNote",
                component: () => import("@/views/personal_center/my_note.vue"),
            },
        ],
    },
    {
        path: "/game",
        name: "game",
        component: () => import("@/views/game/frame.vue"),
        children: [
            {
                path: "",
                name: "gDefault",
                component: () => import("@/views/game/default.vue"),
            },
            {
                path: "flip",
                name: "gFlip",
                component: () => import("@/views/game/flip.vue"),
            },
        ],
    },
    {
        path: "/404", // 考虑后续可能编写404页面，此处预留路由
        name: "notFound",
        redirect: { name: "home" },
    },
    {
        path: "/:pathMatch(.*)*", // 将匹配所有内容并将其放在 `$route.params.pathMatch` 下
        redirect: { name: "notFound" },
    },
]

export const router = createRouter({
    history: createWebHashHistory(),
    routes: routes,
})


if (__IsDev__) {
    router.addRoute({
        path: "/test",
        name: "test",
        meta: { hideTop: true },
        component: () => import("@/views/components/test_page.vue"),
        children: [{
            path: "lock-button",
            name: "TLockButton",
            component: () => import("@/components/lock_button_test.vue")
        }, {
            path: "lock-screen",
            name: "TLockScreen",
            component: () => import("@/components/lock_screen_test.vue")
        }]
    })
}

router.beforeEach((to, _from, next) => {
    const userStore = useUserStore()

    // 页面不需要登录，或者已经登录
    if (!(to.meta && to.meta.requireLogin) || userStore.isLogin()) {
        next()
        return
    }

    next({ path: "/" })
    return
})
