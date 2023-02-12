import { createApp } from "vue"
import "./style.less"
import App from "./app.vue"

import router from "./router"
import store from "./store"

// element message
import { ElMessage } from "element-plus"
import "element-plus/es/components/message/style/css"

// icons
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

// axios init interceptors
import { initInterceptors } from "@/ts/axios/config/config"

initInterceptors((): void => {
  router.replace({ name: "home", params: { v: "1" } })
})

// init
import { displayIsPublic, displayTime } from "@/ts/utils";

let app = createApp(App)
app.config.globalProperties.$filters = {
  displayTime: displayTime,
  displayIsPublic: displayIsPublic,
}

for (const [ key, component ] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(router).use(store).use(ElMessage).mount("#app")
