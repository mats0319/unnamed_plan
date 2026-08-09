<template>
  <div class="test-page color-bg-0">
    <div class="tp-left">
      <p class="tpl-item" @click="testHome">> Home</p>
      <p class="tpl-item center-hv" @click="testLockButton">测试锁定按钮</p>
      <p class="tpl-item center-hv" @click="testLockScreen">测试全局锁定</p>
    </div>

    <el-divider direction="vertical" />

    <div class="tp-right">
      <router-view />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { routerLink } from "@/ts/util.js"

function testHome() {
    routerLink("test")
}

function testLockButton() {
    routerLink("TLockButton")
}

function testLockScreen() {
    routerLink("TLockScreen")
}
</script>

<style lang="less" scoped>
.test-page {
  display: flex;
  height: 100vh;

  .tp-left {
    width: calc(20vw - 1px);
    margin: 6rem 2vw;
    font-size: 1.4rem;

    .tpl-item {
      margin: 3rem auto;
    }
    .tpl-item:hover {
      cursor: pointer;
      text-decoration-line: underline;
    }
  }

  .tp-right {
    width: 60vw;
    margin: 8rem 8vw;
  }

  .el-divider--vertical {
    width: 0;
    height: 100%;
    margin: 0;
    border-left: 1px solid black;
  }
}
</style>
