<template>
  <div class="home">
    <div class="top">
      <div class="title" @click="routerLink('default')">有凤来仪</div>

      <div class="links">
        <span class="link-item" @click="routerLink('games')">小游戏</span>
        <span class="link-item" @click="externalLink('https://github.com/mats9693/unnamed_plan')">Github</span>
      </div>
    </div>

    <router-view />

    <div class="bottom">
      <a href="https://beian.miit.gov.cn/" target="_blank">吉ICP备2021002242号-1</a>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class Home extends Vue {
  private mounted() {
    // placeholder
  }

  private routerLink(name: string): void {
    this.$router.push({ name: name });
  }

  private externalLink(url: string): void {
    window.open(url, "_blank");
  }
}
</script>

<style lang="scss">
.home {
  width: 100vw;
  height: 100vh;

  .top {
    height: 10rem;
    padding: 0 10vw;
    display: flex;

    position: sticky;
    top: 0;

    background-color: aliceblue; // for test

    .title {
      width: 10vw;
      line-height: 10rem;

      font-size: 4rem;
      font-weight: 600;
    }
    .title:hover {
      cursor: pointer;
    }

    .links {
      width: 70vw;
      height: 10rem;
      display: flex;
      justify-content: flex-end;

      font-size: 2rem;

      .link-item {
        width: 7vw;
        margin: auto 1vw;
      }
      .link-item:hover {
        cursor: pointer;
      }
    }
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
