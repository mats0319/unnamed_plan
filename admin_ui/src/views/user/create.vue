<template>
  <div class="user-create">
    <el-form label-position="left" label-width="15%">
      <el-form-item label="用户名">
        <el-input v-model="userName" placeholder="请输入新用户的用户名" />
      </el-form-item>

      <el-form-item label="密码">
        <el-input
          v-model="password"
          placeholder="请输入新用户的密码"
          type="password"
          show-password
          clearable
        />
      </el-form-item>

      <el-form-item label="权限等级">
        <el-select v-model="permissionStr" placeholder="请选择新用户的权限等级" clearable>
          <el-option
            v-for="item in $store.state.permission-1"
            :key="item"
            :label="item"
            :value="item.toString()"
          />
        </el-select>
      </el-form-item>

      <el-form-item>
        <el-button type="info" plain @click="beforeCreateUser">
          创建新用户
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import axios from "axios";
import { calcSHA256 } from "@/ts/sha256";

@Component
export default class UserCreate extends Vue {
  private userName = "";
  private password = "";
  private permissionStr = "";

  private mounted() {
    // placeholder
  }

  private createUser(): void {
    const pwd = calcSHA256(this.password);
    this.password = "";

    let data: FormData = new FormData();
    data.append("operatorID", this.$store.state.userID);
    data.append("userName", this.userName);
    data.append("password", pwd);
    data.append("permission", this.permissionStr);

    axios.post(process.env.VUE_APP_user_create_url, data)
      .then(response => {
        if (response.data.hasError) {
          throw response.data.data;
        }

        const payload = JSON.parse(response.data.data as string);
        if (payload.isSuccess) {
          this.$message.success("创建新用户成功");

          this.userName = "";
          this.permissionStr = "";
        } else {
          this.$message.error("创建新用户失败");
        }
      })
      .catch(err => {
        console.log("create user failed, error:", err);
      })
  }

  private beforeCreateUser(): void {
    let isAllowed = true;
    let errMsg = "";
    if (this.userName.length < 1) {
      isAllowed = false;
      errMsg = "请填写新用户的用户名";
    } else if (this.password.length < 1) {
      isAllowed = false;
      errMsg = "请填写新用户的密码";
    } else if (this.permissionStr.length < 1) {
      isAllowed = false;
      errMsg = "请选择新用户的权限等级";
    }

    if (!isAllowed) {
      this.$message.info(errMsg);
      return;
    }

    this.createUser();
  }
}
</script>

<style lang="scss">
.user-create {
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

  .el-select {
    width: 50%;
  }
}
</style>
