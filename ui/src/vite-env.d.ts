/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

import {Store} from "vuex";

declare module '@vue/runtime-core' {
  interface State {
    isLogin: boolean,
    userID: string,
    nickname: string,
    permission: number,

    ARankAdminPermission: number,
    SRankAdminPermission: number,
  }

  interface Filters {
    displayTime: Function,
    displayIsPublic: Function,
  }

  interface ComponentCustomProperties {
    $store: Store<State>,
    $filters: Filters,
  }
}
