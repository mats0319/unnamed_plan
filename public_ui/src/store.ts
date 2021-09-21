import Vue from "vue"
import Vuex from "vuex"

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    isLogin: false,
    userID: "",
    nickname: "",
    permission: 0,
  }
})
