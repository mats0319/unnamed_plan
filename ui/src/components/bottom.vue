<template>
  <div class="bottom">
    <a href="https://beian.miit.gov.cn/" target="_blank">&#91;{{ domainNumber }}&#93;</a>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "BottomComponent",
  data() {
    return {
      domainNumber: "假装这里有一个域名的备案号",
    }
  }
})
</script>

<style lang="less">
.bottom {
  display: flex;

  a {
    margin: auto;
    color: black;
    font-size: 2rem;
  }
}
</style>
