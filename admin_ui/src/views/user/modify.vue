<template>
  <div class="user-modify">
    <el-form label-position="left" label-width="15%">
      <el-form-item label="当前密码">
        <el-input
          v-model="currPassword"
          placeholder="请输入当前密码"
          type="password"
          show-password
          clearable
        />
      </el-form-item>

      <el-form-item label="新的昵称">
        <el-input v-model="newNickname" placeholder="请输入新的昵称" />

        <el-popover trigger="hover" placement="top" :content="tips_User_Create">
          <i slot="reference" class="el-icon-warning-outline" />
        </el-popover>
      </el-form-item>

      <el-form-item label="新的密码">
        <el-input
          v-model="newPassword"
          placeholder="请输入新的密码"
          type="password"
          show-password
          clearable
        />

        <el-popover trigger="hover" placement="top" :content="tips_User_Create">
          <i slot="reference" class="el-icon-warning-outline" />
        </el-popover>
      </el-form-item>

      <el-form-item>
        <el-button type="info" plain @click="beforeModifyUserInfo">修改</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { tips_User_Create } from "shared_ui/ts/const";
import userAxios from "shared_ui/ts/axios_wrapper/user";

@Component
export default class UserModifyInfo extends Vue {
  private currPassword = "";
  private newNickname = "";
  private newPassword = "";

  // const
  private tips_User_Create = tips_User_Create;

  private mounted() {
    // placeholder
  }

  private modifyUserInfo(): void {
    userAxios.modifyInfo(this.$store.state.userID, this.$store.state.userID, this.currPassword, this.newNickname, this.newPassword)
      .then(response => {
        if (response.data["hasError"]) {
          throw response.data["data"];
        }

        const payload = JSON.parse(response.data["data"] as string);
        if (payload.isSuccess) {
          this.$message.success("修改用户信息成功");
        } else {
          this.$message.error("修改用户信息失败");
        }

        if (this.newPassword.length > 0) {
          this.$store.state.isLogin = false;

          sessionStorage.removeItem("auth");

          this.$router.push({ name: "login" });
        }
      })
      .catch(err => {
        this.$message.error("修改用户信息失败，错误：" + err);
      })
      .finally(() => {
        this.currPassword = "";
      })
  }

  private beforeModifyUserInfo(): void {
    let isAllowed = true;
    let errMsg = "";
    if (this.currPassword.length < 1) {
      isAllowed = false;
      errMsg = "请输入用户当前密码";
    } else if (this.newNickname.length + this.newPassword.length < 1) {
      isAllowed = false;
      errMsg = "请修改用户名或密码后尝试提交";
    }

    if (!isAllowed) {
      this.$message.info(errMsg);
      return;
    }

    this.modifyUserInfo();
  }
}
</script>

<style lang="scss">
.user-modify {
  padding: 7vh 15vw;
  text-align: left;

  .el-form-item {
    margin: 5vh 0;
  }

  .el-form-item__label {
    font-size: 2rem;
  }

  .el-input {
    width: 60%;
  }

  .el-popover__reference-wrapper {
    margin-left: 5vh;
  }
}
</style>
