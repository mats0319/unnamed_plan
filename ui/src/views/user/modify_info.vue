<template>
  <el-form class="user-modify-info" label-position="left" label-width="25%">
    <el-form-item label="新的昵称">
      <el-input v-model="nickname" placeholder="请输入新的昵称" />

      <el-tooltip :content="tips.userModify" effect="light" placement="top">
        <el-icon size="2rem"><InfoFilled /></el-icon>
      </el-tooltip>
    </el-form-item>

    <el-form-item label="新的密码">
      <el-input v-model="password" type="password" placeholder="请输入新的密码" show-password clearable />

      <el-tooltip :content="tips.userModify" effect="light" placement="top">
        <el-icon size="2rem"><InfoFilled /></el-icon>
      </el-tooltip>
    </el-form-item>

    <el-form-item label="当前密码">
      <el-input v-model="currPassword" type="password" placeholder="请输入新的密码" show-password clearable />
    </el-form-item>

    <el-form-item>
      <el-button type="info" plain @click="beforeModifyUserInfo">修改</el-button>
    </el-form-item>
  </el-form>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { tips } from "@/ts/const";
import userAxios from "@/ts/axios/user";

export default defineComponent({
  name: "UserModifyInfo",
  data() {
    return {
      nickname: "",
      password: "",
      currPassword: "",

      // const
      tips: tips,
    }
  },
  mounted() {
    // placeholder
  },
  methods: {
    modifyUserInfo(): void {
      userAxios.modifyInfo(
        this.$store.state.userID,
        this.$store.state.userID,
        this.currPassword,
        this.nickname,
        this.password,
      )
        .then(response => {
          if (response.err) {
            throw response.err
          }

          this.$message.success("修改我的信息成功")

          if (this.nickname.length > 0) {
            this.$store.state.nickname = this.nickname
          }
        })
        .catch(err => {
          this.$message.error("修改我的信息失败")
          console.log("> modify user info failed, error: ", err)
        })
        .finally(() => {
          this.currPassword = ""
        })
    },

    beforeModifyUserInfo(): void {
      if (this.currPassword.length < 1) {
        this.$message.info("请输入密码")
        return
      } else if (this.nickname.length + this.password.length < 1) {
        this.$message.info("请输入新的昵称或密码后重试")
        return
      }

      this.modifyUserInfo()
    }
  }
})
</script>

<style lang="less">
.user-modify-info {
  margin: 7rem 20vw 0;

  .el-form-item {
    margin: 5vh 0;
  }

  .el-form-item__label {
    font-size: 2rem;
  }

  .el-input {
    width: 60%;
  }

  .el-icon {
    margin-left: 5%;
  }
}

.el-popper {
  font-size: 1.8rem;
}
</style>
