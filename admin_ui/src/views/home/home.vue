<template>
  <div class="home">
    <div class="home-navigation">
      <div class="hn-user">
        <div class="hnu-name">
          <span class="hnu-link-home" @click="linkHome">{{ $store.state.nickname }}</span>
        </div>

        <div class="hnu-permission">
          <span class="hnu-link-home" @click="linkHome">权限等级：{{ $store.state.permission }}</span>
        </div>

        <el-button type="info" size="mini" @click="exit" plain>退出登录</el-button>
      </div>

      <el-menu class="hn-items" router unique-opened>
        <el-submenu index="user">
          <template slot="title">用户管理</template>

          <el-menu-item index="modify" :route="{ name: 'userModify' }">修改当前用户信息</el-menu-item>
          <el-menu-item index="create" :route="{ name: 'userCreate' }">创建新用户</el-menu-item>
          <el-menu-item index="list" :route="{ name: 'userList' }">查看其它用户</el-menu-item>
        </el-submenu>
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
  private mounted() {
    if (!this.$store.state.isLogin) {
      this.exit();
    }
  }

  private linkHome(): void {
    if (location.href.split("#/")[1].length > 0) {
      this.$router.push({ name: "home" });
    }
  }

  private exit(): void {
    this.$store.state.userID = "";

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

      border-right: 1px solid rgba(230, 230, 230, 1);
      border-bottom: 1px solid rgba(230, 230, 230, 1);

      .hnu-name {
        height: 3rem;
        padding-top: 5rem;
        font-size: 2rem;
      }

      .hnu-permission {
        height: 5rem;
        font-size: 1.6rem;
      }

      .hnu-link-home:hover {
        cursor: pointer;
      }
    }

    .hn-items {
      height: calc(100vh - 20rem);
      text-align: left;

      .el-submenu__title {
        font-size: 1.8rem;
      }

      .el-menu-item {
        font-size: 1.6rem;
      }
    }
  }

  .home-content {
    width: calc(90% - 20rem);
    height: 80vh;
    padding: 10vh 5%;
  }
}
</style>
