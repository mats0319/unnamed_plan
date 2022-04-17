<template>
  <div id="app" class="unsupported">
    <span class="text">
      该网站暂不支持<br />
      移动端设备访问
    </span>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class Unsupported extends Vue {
  private mounted() {
    // placeholder
  }
}
</script>

<style lang="less">
html {
  font-size: 62.5%; // now, 1 rem = 10 px
}

body {
  margin: 10rem 0;
  padding: 0 5vw;
  overflow-x: hidden;
}

#app {
  font-family: KaiTi, STKaiti, Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.unsupported {
  width: 90vw;
  height: 10rem;
  display: flex;

  .text {
    width: inherit;
    font-size: 3rem;
    text-align: center;
  }
}
</style>
