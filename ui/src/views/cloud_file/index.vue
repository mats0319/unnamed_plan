<template>
  <content-template
    :index-page-content="indexPageContent"
    :index-page-router-name="indexPageRouterName"
    :sub-page-nav="subPageNav"
  />
</template>

<script lang="ts">
import { defineComponent } from "vue";
import ContentTemplate from "@/components/content_template.vue";
import { SubPageNavItem } from "@/ts/structure";

export default defineComponent({
  name: "CloudFileIndex",
  components: { ContentTemplate },
  data() {
    return {
      indexPageRouterName: "cloudFile",
      indexPageContent: "云文件",
      subPageNav: [
        {
          name: "查看文件列表",
          routerName: "cloudFileList",
        },
        {
          name: "上传新文件",
          routerName: "cloudFileUpload",
        },
      ] as SubPageNavItem[],
    }
  }
})
</script>
