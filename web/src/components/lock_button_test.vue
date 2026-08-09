<template>
  <p>测试结论：对Promise不熟悉的，可以考虑使用async/await编写处理函数</p>

  <el-divider />

  <p>测试组件：锁定按钮；在一次点击之后、直到功能执行完成并允许下一次点击之前，禁止点击行为、不响应点击事件</p>
  <div class="wrapper">
    <lock-button :onClick="clickButton">测试功能</lock-button>
    <div :style="{ width: '4rem'}"></div>
    <lock-button :disabled="true" :on-click="clickButton">禁用的lock button</lock-button>
    <div :style="{ width: '4rem'}"></div>
    <outlined-button :disabled="true" :on-click="clickButton">禁用的outlined button</outlined-button>
    <div :style="{ width: '4rem'}"></div>
    <elevated-button :disabled="true" :on-click="clickButton">禁用的elevated button</elevated-button>
  </div><br/>
  <p>count:&emsp;{{ count }}</p>

  <el-divider />

  <p>测试样式：修改按钮组件的颜色和边框</p>
  <div class="wrapper">
    <outlined-button class="fat-button" :on-click="clickButton">测试样式-带边框的文字按钮</outlined-button>

    <div :style="{ width: '6rem' }"></div>

    <elevated-button class="fat-button" bg="rgb(240, 239, 226)" :on-click="clickButton">测试样式-增强按钮</elevated-button>
  </div>
</template>

<script lang="ts" setup>
import LockButton from "@/components/lock_button.vue"
import OutlinedButton from "@/components/outlined_button.vue"
import ElevatedButton from "@/components/elevated_button.vue"
import { ref } from "vue"

const count = ref<number>(0)

async function clickButton(): Promise<void> {
    await sleep(1_500)
    count.value++
}

// function clickButton(): Promise<void> { // 返回Promise，按钮状态在函数执行完成后变化
//     return sleep(1000).
//         then(() => {
//             console.log("in then")
//         }).finally(() => {
//             console.log("in finally")
//             count.value++
//         })
// }

// function clickButton() { // 没有返回Promise，按钮立刻解除锁定
//     sleep(1000).
//         then(() => {
//             console.log("in then")
//         }).finally(() => {
//             console.log("in finally")
//             count.value++
//         })
// }

// function clickButton(): void {
//     count.value++
// }

async function sleep(time: number): Promise<void> {
    // resolve通知Promise执行结束，必须要调，不然就会一直卡住、不会进入.then等后续流程
    await new Promise(resolve => setTimeout(resolve, time))
}
</script>

<style lang="less" scoped>
.el-divider {
  height: 0;
  border: 1px solid darkgray;
}

.wrapper{
  display: flex;

  .fat-button {
    width: 16rem;
    height: 4rem;
  }
}
</style>
