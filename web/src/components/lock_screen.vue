<template>
  <!--  使用前请将该组件注册到项目中，注册到app.vue都可以-->
  <div></div>
</template>

<script lang="ts" setup>
import { useFlagStore } from "@/pinia/flag.ts"
import { ElLoading, LoadingInstance } from "element-plus"
import { ref, watch } from "vue"

const flagStore = useFlagStore()

const fullScreenLoadingFlag = ref<LoadingInstance>()

function fullScreenLoading() {
    fullScreenLoadingFlag.value = ElLoading.service({
        lock: true,
        text: "loading...",
        background: "rgba(240, 239, 226, 0.5)",
    })
}

watch(() => flagStore.loading, (newValue: boolean) => {
    if (newValue) {
        fullScreenLoading()
    } else {
        fullScreenLoadingFlag.value?.close()
    }
})
</script>

<style lang="scss" scoped>

</style>
