<template>
  <div class="top">
    <span class="t-name" @click="linkTo('/', 'home')">上弦月</span>

    <span class="t-blank">&nbsp;</span>

    <span class="t-item">
      <a href="https://github.com/mats9693/unnamed_plan" target="_blank">本站代码</a>
    </span>

    <span v-if="!$store.state.isLogin" class="t-item" @click="openLoginDialog">登录</span>
    <el-dropdown v-else class="t-item" popper-class="t-dropdown-options" :hide-timeout="300">
      <span class="ti-fold-text">{{ $store.state.nickname }}</span>
      <template #dropdown>
        <el-dropdown-item @click="linkTo('/cloud-file', 'cloudFile')">
          <span class="tdo-item">云文件</span>
        </el-dropdown-item>

        <el-dropdown-item divided @click="linkTo('/user', 'user')">
          <span class="tdo-item">用户</span>
        </el-dropdown-item>

        <el-dropdown-item divided @click="exit">
          <span class="tdo-item">退出登录</span>
        </el-dropdown-item>
      </template>
    </el-dropdown>

    <!--  login  -->
    <el-dialog v-model="showLoginDialog" class="t-login-dialog">
      <template #header>
        <span class="tld-title">登录</span>
      </template>

      <el-form label-position="left" label-width="20%">
        <el-form-item label="用户名">
          <el-input v-model="userName" placeholder="请输入用户名" />
        </el-form-item>

        <el-form-item label="密码">
          <el-input v-model="password" type="password" placeholder="请输入密码" clearable show-password />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="closeLoginDialog">取消</el-button>
        <el-button type="info" @click="login">登录</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import userAxios from "@/ts/axios/user"

export default defineComponent({
  name: "TopComponent",
  data() {
    return {
      showLoginDialog: false,
      userName: "",
      password: "",
    }
  },
  mounted() {
    // placeholder
  },
  methods: {
    login(): void {
      if (this.userName.length < 1 || this.password.length < 1) {
        this.$message.info("请输入用户名和密码")
        return
      }

      userAxios.login(this.userName, this.password)
        .then(response => {
          if (response.err) {
            throw response.err
          }

          sessionStorage.setItem("auth", import.meta.env.Vite_axios_source_sign as string)

          this.userName = ""

          this.setLoginData(response.user_id, response.nickname, response.permission)

          this.showLoginDialog = false

          this.$message.success("登录成功")
        })
        .catch(err => {
          this.$message.error("登录失败")
          console.log("> login failed, error: ", err)
        })
        .finally(() => {
          this.password = ""
        });
    },

    exit(): void {
      this.emptyLoginData();

      sessionStorage.removeItem("auth");
      sessionStorage.removeItem("user");
      sessionStorage.removeItem("token");

      this.linkTo("/", "home");
    },

    linkTo(path: string, name: string): void {
      if (location.href.split("#")[1] !== path) {
        this.$router.push({ name: name })
      }
    },

    setLoginData(userID: string, nickname: string, permission: number): void {
      this.$store.state.isLogin = true;
      this.$store.state.userID = userID;
      this.$store.state.nickname = nickname;
      this.$store.state.permission = permission;
    },

    emptyLoginData(): void {
      this.$store.state.isLogin = false;
      this.$store.state.userID = "";
      this.$store.state.nickname = "";
      this.$store.state.permission = 0;
    },

    openLoginDialog(): void {
      this.showLoginDialog = true
    },

    closeLoginDialog(): void {
      this.showLoginDialog = false
    }
  }
})
</script>

<style lang="less">
.top {
  padding: 0 5vw;

  display: flex;
  align-items: center;

  .t-name {
    width: 10vw;

    font-size: 4rem;
    font-weight: 600;
    letter-spacing: .5rem;

    &:hover {
      cursor: pointer;
    }
  }

  .t-blank {
    width: 60vw;
  }

  .t-item {
    width: 7vw;
    margin: 0 1.5vw;

    font-size: 2.5rem;

    a {
      color: black;
      text-decoration: none;
    }

    .ti-fold-text {
      width: inherit;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    &:hover {
      cursor: pointer;
    }
  }

  .t-login-dialog {
    text-align: left;

    .tld-title {
      font-size: 3rem;
      font-weight: 600;
    }

    .el-form-item {
      margin-left: 20%;

      .el-input {
        width: 60%;
      }
    }

    .el-form-item__label {
      font-size: 2rem;
    }
  }
}

.t-dropdown-options {
  .tdo-item {
    font-size: 1.6rem;
  }
}
</style>
