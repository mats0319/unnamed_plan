<template>
  <el-form class="user-create" label-position="left" label-width="25%">
    <el-form-item label="用户名">
      <el-input v-model="userName" placeholder="请输入新用户的用户名" />
    </el-form-item>

    <el-form-item label="密码">
      <el-input v-model="password" type="password" placeholder="请输入新用户的密码" show-password clearable />
    </el-form-item>

    <el-form-item label="权限等级">
      <el-select v-model="permission">
        <el-option v-for="item in $store.state.permission-1" :key="item" :label="item" :value="item" />
      </el-select>
    </el-form-item>

    <el-divider />

    <el-form-item label="操作员密码">
      <el-input v-model="operatorPassword" type="password" placeholder="请输入操作员的密码" show-password clearable />
    </el-form-item>

    <el-form-item>
      <el-button type="info" plain @click="beforeCreateUser">创建</el-button>
    </el-form-item>
  </el-form>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import userAxios from "@/ts/axios/user";

export default defineComponent({
  name: "UserCreate",
  data() {
    return {
      userName: "",
      password: "",
      permission: 0,
      operatorPassword: "",
    }
  },
  mounted() {
    // placeholder
  },
  methods: {
    createUser(): void {
      userAxios.create(this.$store.state.userID, this.userName, this.password, this.permission, this.operatorPassword)
        .then(response => {
          if (response.err) {
            throw response.err
          }

          this.$message.success("创建新用户成功")
        })
        .catch(err => {
          this.$message.error("创建新用户失败")
          console.log("> create user failed, error: ", err)
        })
        .finally(() => {
          this.password = ""
        })
    },

    beforeCreateUser(): void {
      if (this.userName.length < 1 || this.password.length < 1) {
        this.$message.info("请输入新用户的用户名和密码")
        return
      } else if (this.operatorPassword.length < 1) {
        this.$message.info("请输入操作员密码")
        return
      }

      this.createUser()
    }
  }
})
</script>

<style lang="less">
.user-create {
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
