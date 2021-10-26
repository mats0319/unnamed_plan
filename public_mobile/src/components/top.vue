<template>
  <div>
    <div class="top-title" @click="linkTo('/', 'home')">上弦月</div>

    <div class="top-user-name">
      <span v-show="$store.state.isLogin" @click="exit">{{ $store.state.nickname }}</span>
    </div>

    <div class="top-menu" @click="loginOrShowMenu">
      <i class="el-icon-s-operation" />
    </div>

    <el-dialog
      class="top-login-dialog"
      :visible.sync="loginDialogController"
      width="90vw"
      :modal-append-to-body="false"
      append-to-body
    >
      <div slot="title" class="tld-title">登录</div>

      <div class="tld-content">
        <div class="tldc-item">
          <span class="tldci-label">用户名&#58;</span>
          <el-input v-model="userName" placeholder="请输入用户名" />
        </div>

        <div class="tldc-item">
          <span class="tldci-label">密码&#58;</span>
          <el-input v-model="password" type="password" placeholder="请输入密码" clearable />
        </div>
      </div>

      <div slot="footer">
        <el-button @click="cancelLogin">取消</el-button>
        <el-button type="info" @click="auth">登录</el-button>
      </div>
    </el-dialog>

    <el-drawer
      class="top-menu-drawer"
      :visible.sync="menuDrawerController"
      append-to-body
      size="70%"
    >
      <el-menu background-color="#f0f8ff" unique-opened router>
        <el-submenu index="cloud-file">
          <template slot="title">
            <span>云文件</span>
          </template>

          <el-menu-item index="/cloud-file/list-by-uploader" @click="closeDrawer">我上传的文件</el-menu-item>
        </el-submenu>
      </el-menu>
    </el-drawer>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import homeAxios from "shared_ui/ts/axios_wrapper/home";

@Component
export default class Top extends Vue {
  private userName = "";
  private password = "";

  private loginDialogController = false;
  private menuDrawerController = false;

  private mounted() {
    // placeholder
  }

  private loginOrShowMenu(): void {
    if (!this.$store.state.isLogin) {
      this.openLoginDialog();
    } else {
      this.menuDrawerController = true;
    }
  }

  private auth(): void {
    homeAxios.login(this.userName, this.password)
      .then(response => {
        if (response.data["hasError"]) {
          throw response.data["data"];
        }

        sessionStorage.setItem("auth", "passed");

        const payload = JSON.parse(response.data["data"] as string);
        this.$store.state.userID = payload.userID;
        this.$store.state.nickname = payload.nickname.length > 15 ? payload.nickname.slice(0, 15) + "..." : payload.nickname;
        this.$store.state.permission = payload.permission;

        this.userName = "";

        this.$store.state.isLogin = true;
        this.loginDialogController = false;
      })
      .catch(err => {
        this.$message.error("登录失败，错误：" + err);
      })
      .finally(() => {
        this.password = "";
      });
  }

  private exit(): void {
    this.$store.state.isLogin = false;

    sessionStorage.removeItem("auth");

    this.linkTo("/", "home");
  }

  private linkTo(path: string, name: string): void {
    if (location.href.split("#")[1] !== path) {
      this.$router.push({ name: name });
    }
  }

  private openLoginDialog(): void {
    this.userName = "";
    this.password = "";

    this.loginDialogController = true;
  }

  private cancelLogin(): void {
    this.loginDialogController = false;
  }

  private closeDrawer(): void {
    this.menuDrawerController = false;
  }
}
</script>

<style lang="scss">
.top-title {
  width: 40vw;
  line-height: 10rem;

  font-size: 3rem;
  font-weight: 600;
}

.top-user-name {
  width: 40vw;
  line-height: 10rem;
  font-size: 2rem;
}

.top-menu {
  width: 20vw;
  line-height: 10rem;
  font-size: 3rem;
}

.top-login-dialog {
  text-align: left;

  .tld-title {
    font-size: 2.5rem;
    font-weight: 600;
  }

  .tld-content {
    padding: 0 10%;
    font-size: 1.8rem;

    .tldc-item {
      display: flex;
      padding: 2vh 0;

      .tldci-label {
        width: 30%;
        align-self: center;
      }

      .el-input {
        width: 70%;
      }
    }
  }
}

.top-menu-drawer {
  .el-drawer {
    background-color: aliceblue;
  }
}
</style>
