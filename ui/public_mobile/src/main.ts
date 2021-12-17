import Vue from "vue"
import App from "./app.vue"
import router from "./router"
import store from "./store"

Vue.config.productionTip = false

// element-ui
import ElementUI from "element-ui";
import 'element-ui/lib/theme-chalk/index.css';

Vue.use(ElementUI);

// filters
import { displayIsPublic, displayTime } from "shared/ts/data";

Vue.filter("displayIsPublic", displayIsPublic);
Vue.filter("displayTime", displayTime);

// axios init interceptors
import { initInterceptors } from "shared/ts/axios_wrapper/config";

initInterceptors((): void => {
  // placeholder
})

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app")
