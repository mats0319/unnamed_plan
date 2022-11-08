<template>
  <div class="login">
    <div class="login-content">
      <el-card class="lc-card">
        <div>
          <div class="lcc-item">
            <span class="lcci-label">用户名&#58;</span>
            <el-input v-model="userName" placeholder="请输入用户名" />
          </div>

          <div class="lcc-item">
            <span class="lcci-label">密码&#58;</span>
            <el-input v-model="password" type="password" placeholder="请输入密码" clearable />
          </div>

          <div class="lcc-submit">
            <el-button type="primary" @click="auth">登录</el-button>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import homeAxios from "shared/ts/axios_wrapper/home";

@Component
export default class Login extends Vue {
  private userName = "";
  private password = "";

  private mounted() {
    // placeholder
  }

  private auth(): void { // not use 'login', because func name as 'homeAxios.login', idea can't distinguish them
    homeAxios.login(this.userName, this.password)
      .then(response => {
        if (response.data["hasError"]) {
          throw response.data["data"];
        }

        sessionStorage.setItem("auth", process.env.VUE_APP_axios_source_sign as string);

        const payload = JSON.parse(response.data["data"] as string);

        this.setLoginData(payload.userID, payload.nickname, payload.permission);

        this.$router.push({ name: "home" });
      })
      .catch(err => {
        this.$message.error("登录失败");
        console.log("> login failed.", err);
      })
      .finally(() => {
        this.password = "";
      });
  }

  private setLoginData(userID: string, nickname: string, permission: number): void {
    this.$store.state.isLogin = true;
    this.$store.state.userID = userID;
    this.$store.state.nickname = nickname;
    this.$store.state.permission = permission;
  }
}
</script>

<style lang="less">
.login {
  height: 100vh;

  .login-content {
    width: 40vw;
    height: 50vh;

    padding-top: 20vh;
    margin: auto;

    .lc-card {
      height: fit-content;
      font-size: 2.5rem;

      .lcc-item {
        display: flex;
        padding: 3vh 3vw;
        text-align: left;

        .lcci-label {
          width: 20%;
          align-self: center;
        }

        .el-input {
          width: 80%;
        }
      }

      .lcc-submit {
        padding: 5vh 0 3vh;
      }
    }
  }
}
</style>
