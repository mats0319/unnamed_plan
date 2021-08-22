<template>
  <div class="login">
    <div class="login-content">
      <el-card class="lc-card">
        <div>
          <div class="lcc-item">
            <span class="lcci-label">账号&#58;</span>
            <el-input v-model="userName" placeholder="请输入账号" />
          </div>

          <div class="lcc-item">
            <span class="lcci-label">密码&#58;</span>
            <el-input v-model="password" type="password" placeholder="请输入密码" clearable />
          </div>

          <div class="lcc-submit">
            <el-button type="primary" @click="login">登录</el-button>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import axios from "axios";

@Component
export default class Login extends Vue {
  private userName = "";
  private password = "";

  private mounted() {
    // placeholder
  }

  private login(): void {
    let pwd = this.password;
    this.password = "";

    let data: FormData = new FormData();
    data.append("userName", this.userName);
    data.append("password", pwd);

    axios.post(process.env.VUE_APP_login_url, data).then(
      response => {
        if (response.data.hasError) {
          throw response.data;
        }

        sessionStorage.setItem("auth", "passed");

        this.$store.state.userName = this.userName;
        this.$store.state.permission = JSON.parse(response.data.data as string).permission;

        this.$router.push({ name: "home" });
      }
    ).catch(
      err => {
        console.log("login failed, error:", err);
      }
    );
  }
}
</script>

<style lang="scss">
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
