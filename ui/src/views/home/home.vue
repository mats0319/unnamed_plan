<template>
  <div class="home">
    <top class="top" />

    <router-view />

    <div class="bottom">
      <a href="https://beian.miit.gov.cn/" target="_blank">吉ICP备2021002242号-1</a>
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

<style lang="scss">
.home {
  width: 100vw;
  height: 100vh;

  .top {
    height: 10rem;
    padding-left: 10vw;
    display: flex;

    position: sticky;
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
