import { createStore } from "vuex"

let store = createStore({
  state: () => ({
    isLogin: false,
    userID: "",
    nickname: "",
    permission: 0,

    ARankAdminPermission: import.meta.env.Vite_A_rank_admin_permission,
    SRankAdminPermission: import.meta.env.Vite_S_rank_admin_permission,
  })
})

export default store
