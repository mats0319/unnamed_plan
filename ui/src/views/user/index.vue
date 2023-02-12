<template>
  <content-template
    :index-page-content="indexPageContent"
    :index-page-router-name="indexPageRouterName"
    :sub-page-nav="subPageNav"
  />
</template>

<script lang="ts">
import { defineComponent } from "vue";
import ContentTemplate from "@/components/content_template.vue"
import { SubPageNavItem } from "@/ts/structure";

export default defineComponent({
  name: "UserIndex",
  components: { ContentTemplate },
  data() {
    return {
      indexPageRouterName: "user",
      indexPageContent: "用户",
      subPageNav: [
        {
          name: "修改我的信息",
          routerName: "userModifyInfo",
        },
        {
          name: "查看用户列表",
          routerName: "userList",
          permission: 6,
        },
        {
          name: "创建新用户",
          routerName: "userCreate",
          permission: 6,
        },
      ] as SubPageNavItem[],
    }
  },
})
</script>
