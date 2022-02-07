import Vue from "vue"
import Vuex from "vuex"

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    isLogin: false,
    userID: "",
    nickname: "",
    permission: 0,

    cloudFilePageType: "", // 0: list by uploader, 1: list public
    notePageType: "", // 0: list by writer, 1: list public
  }
})
