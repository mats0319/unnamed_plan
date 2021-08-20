<template>
  <div class="files">
    <el-menu unique-opened>
      <el-menu-item>翻牌游戏</el-menu-item>
    </el-menu>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class Files extends Vue {
  private mounted() {
    // placeholder
  }
}
</script>

<style lang="scss">
.files {
  width: inherit;
  height: calc(100vh - 20rem);
  font-size: 5rem;

  display: flex;
  justify-content: center;
  align-items: center;
  margin: auto;
}
</style>
