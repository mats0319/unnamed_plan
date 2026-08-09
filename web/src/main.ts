import { createApp } from "vue"
import { createPinia } from "pinia"

import App from "./app.vue"
import { router } from "./router.ts"

import "element-plus/dist/index.css"
import "./index.less"

// axios init interceptors
import { initInterceptors } from "@/axios/ts/config_extend.ts"
import { UserStatus, useUserStore } from "@/pinia/user.ts"

initInterceptors((): void => {
    const userStore = useUserStore()
    userStore.setUserStatus(UserStatus.exit)

    router.replace({ name: "home" })
})

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount("#app")
