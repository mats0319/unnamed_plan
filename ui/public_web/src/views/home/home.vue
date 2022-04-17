<template>
  <div class="home">
    <top class="top" />

    <div class="content">
      <router-view />
    </div>

    <div class="bottom">
      <a href="https://beian.miit.gov.cn/" target="_blank">&#91;假装这里有一个域名的备案号&#93;</a>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Top from "@/components/top.vue";

@Component({
  components: { Top }
})
export default class Home extends Vue {
  private mounted() {
    // placeholder
  }
}
</script>

<style lang="less">
.home {
  width: 100vw;
  height: 100vh;

  .top {
    height: 10rem;
    padding-left: 10vw;
    display: flex;

    position: sticky;
    z-index: 10;
    top: 0;

    background-color: aliceblue; // for test
  }

  .content {
    height: calc(100vh - 20rem);
  }

  .bottom {
    display: flex;
    height: 10rem;
    background-color: aliceblue;

    a {
      margin: auto;
      color: black;
      font-size: 2rem;
    }
  }
}
</style>
