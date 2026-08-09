<template>
  <!--  锁屏组件在实际注册时你可以注册到任何位置，只要注册过就行-->
  <lock-screen />

  <p>
    结论：<br/>
    如果一个功能执行过程中，不希望用户进行其他操作，则应考虑全屏锁定，例如登录；<br/>
    如果一个功能执行过程中，允许用户进行其他操作，则应考虑局部锁定（锁按钮、表格），例如查询下一页。<br/>
    （锁定：锁定期间禁止点击、不响应点击事件）
  </p>

  <el-divider />

  <p>点击按钮锁定整个屏幕（实际上按钮只修改了信号，锁屏由锁屏组件实现）</p>
  <elevated-button :onClick="testFullScreenLoading">测试全屏锁定</elevated-button>

  <el-divider />

  <p>点击按钮锁定按钮本身和表格</p>
  <elevated-button :disabled="tableLoading" :onClick="testTableLoading">测试局部锁定</elevated-button>
  <br />
  <el-table v-loading="tableLoading" :data="tableData">
    <el-table-column prop="name" label="姓名" />
    <el-table-column prop="score" label="分数" />
  </el-table>
</template>

<script lang="ts" setup>
import { useFlagStore } from "@/pinia/flag.ts"
import LockScreen from "@/components/lock_screen.vue"
import { ref } from "vue"
import ElevatedButton from "@/components/elevated_button.vue"

const flagStore = useFlagStore()

const tableLoading = ref<boolean>(false)
const tableData = ref<Array<{ name: string, score: number }>>([
    { name: "admin", score: 100 },
    { name: "mario", score: 99 },
    { name: "mats0319", score: 98 },
    { name: "matongshuai", score: 97 },
])

function testFullScreenLoading() {
    flagStore.setLoading(true)
    setTimeout(() => {
        flagStore.setLoading(false)
    }, 1500)
}

function testTableLoading() {
    tableLoading.value = true

    tableData.value = tableData.value.reverse()
    setTimeout(() => {
        tableLoading.value = false
    }, 2000)
}
</script>

<style lang="less" scoped>
.el-divider {
  height: 0;
  border: 1px solid darkgray;
}
</style>
