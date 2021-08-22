<template>
  <div class="home">
    <div class="home-navigation">
      <div class="hn-user">
        <div class="hnu-name">{{userName}}</div>

        <el-button type="info" size="mini" @click="exit" plain>退出登录</el-button>
      </div>

      <el-menu class="hn-items" router unique-opened>
        <el-menu-item index="1" :route="{ name: 'user' }">用户管理</el-menu-item>
      </el-menu>
    </div>

    <div class="home-content">
      <router-view />
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class Home extends Vue {
  private userName = "";
  private permission = 0;

  private mounted() {
    this.userName = this.$store.state.userName;
    this.permission = this.$store.state.permission;

    if (this.userName.length < 1 || !this.permission) {
      this.$router.push({ name: "login" });
    }
  }

  private exit(): void {
    this.$store.state.userName = "";
    this.$store.state.permission = 0;

    sessionStorage.removeItem("auth");

    this.$router.push({ name: "login" });
  }
}
</script>

<style lang="scss">
.home {
  height: 100vh;

  display: flex;

  .home-navigation {
    width: 20rem;
    height: inherit;

    .hn-user {
      height: 19.9rem;

      font-size: 2rem;

      border-right: 1px solid rgba(230, 230, 230, 1);
      border-bottom: 1px solid rgba(230, 230, 230, 1);

      .hnu-name {
        height: 10rem;
        padding-top: 3rem;
      }
    }

    .hn-items {
      height: calc(100vh - 20rem);

      .el-menu-item {
        font-size: 2rem;
      }
    }
  }

  .home-content {
    width: calc(100% - 20rem);
    height: inherit;
  }
}
</style>
