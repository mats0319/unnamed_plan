<template>
  <div class="user-list">
    <el-table class="ul-table" :data="users" stripe highlight-current-row>
      <el-table-column prop="user_name" label="用户名" min-width="2" />
      <el-table-column prop="nickname" label="昵称" min-width="2" />
      <el-table-column prop="permission" label="权限等级" min-width="1" />
      <el-table-column v-if="$store.state.permission >= $store.state.ARankAdminPermission" label="操作" min-width="3">
        <template #default="scope">
          <el-button type="info" plain @click="beforeLockOrUnlockUser(scope.$index, !scope.row.is_locked)">
            {{ !scope.row.is_locked ? "锁定" : "解锁" }}
          </el-button>

          <el-button
            v-if="$store.state.permission >= $store.state.SRankAdminPermission"
            type="info"
            plain
            @click="beforeModifyUserPermission(scope.$index)"
          >
            修改权限等级
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      :current-page="pageNum"
      :page-size="pageSize"
      :total="total"
      :layout="layout"
      @update:current-page="onPageChange"
    />

    <!--  lock or unlock  -->
    <el-dialog v-model="showLockOrUnlockDialog" class="ul-dialog">
      <template #header>
        <span class="uld-title">{{ wantLock ? "锁定" : "解锁" }}</span>
      </template>

      <el-form label-position="left" label-width="30%">
        <el-form-item label="用户ID">{{ users[userIndex].user_id }}</el-form-item>

        <el-form-item label="用户名">{{ users[userIndex].user_name }}</el-form-item>

        <el-form-item label="用户状态">{{ users[userIndex].is_locked ? "已锁定" : "未锁定" }}</el-form-item>

        <el-divider />

        <el-form-item label="密码">
          <el-input v-model="password" type="password" placeholder="请输入密码" show-password clearable />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showLockOrUnlockDialog = false">取消</el-button>
        <el-button type="info" @click="lockOrUnlockUser">{{ wantLock ? "锁定" : "解锁" }}</el-button>
      </template>

    </el-dialog>

    <!--  modify user permission  -->
    <el-dialog v-model="showModifyPermissionDialog" class="ul-dialog">
      <template #header>
        <span class="uld-title">修改权限等级</span>
      </template>

      <el-form label-position="left" label-width="30%">
        <el-form-item label="用户ID">{{ users[userIndex].user_id }}</el-form-item>

        <el-form-item label="用户名">{{ users[userIndex].user_name }}</el-form-item>

        <el-form-item label="原权限等级">{{ users[userIndex].permission }}</el-form-item>

        <el-divider />

        <el-form-item label="新权限等级">
          <el-select v-model="permission">
            <el-option v-for="item in 7" :key="item" :label="item" :value="item" />
          </el-select>
        </el-form-item>

        <el-form-item label="密码">
          <el-input v-model="password" type="password" placeholder="请输入密码" show-password clearable />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showModifyPermissionDialog = false">取消</el-button>
        <el-button type="info" @click="modifyUserPermission">修改</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { User } from "@/ts/axios/proto/1_user.pb";
import userAxios from "@/ts/axios/user";

export default defineComponent({
  name: "UserList",
  data() {
    return {
      users: [] as User.Data[],
      pageSize: 10,
      pageNum: 1,
      total: 0,
      layout: "prev,pager,next,->,total",

      showLockOrUnlockDialog: false,
      wantLock: false,
      userIndex: 0, // for update data when operate success
      password: "",

      showModifyPermissionDialog: false,
      permission: 0,
    }
  },
  mounted() {
    this.listUser()
  },
  methods: {
    listUser(): void {
      userAxios.list(this.$store.state.userID, this.pageSize, this.pageNum)
        .then(response => {
          if (response.err) {
            throw response.err
          }

          this.users = response.users ? response.users : new Array<User.Data>()
          this.total = response.total ? response.total : 0
        })
        .catch(err => {
          this.users = new Array<User.Data>()
          this.total = 0

          this.$message.error("获取用户列表失败")
          console.log("> list user failed, error: ", err)
        })
    },

    onPageChange(newValue: number): void {
      this.pageNum = newValue
      this.listUser()
    },

    lockOrUnlockUser(): void {
      if (this.password.length < 1) {
        this.$message.info("请输入密码")
        return
      }

      let res: Promise<User.LockRes | User.UnlockRes>
      if (this.wantLock) {
        res = userAxios.lock(this.$store.state.userID, this.users[this.userIndex].user_id, this.password)
      } else {
        res = userAxios.unlock(this.$store.state.userID, this.users[this.userIndex].user_id, this.password)
      }

      res
        .then(response => {
          if (response.err) {
            throw response.err
          }

          this.$message.success(this.wantLock ? "锁定" : "解锁" + "用户成功")
          this.users[this.userIndex].is_locked = this.wantLock

          this.showLockOrUnlockDialog = false // close dialog only on success
        })
        .catch(err => {
          this.$message.error(this.wantLock ? "锁定" : "解锁" + "用户失败")
          console.log("> " + this.wantLock ? "锁定" : "解锁" + " user failed, error: ", err)
        })
        .finally(() => {
          this.password = ""
        })
    },

    modifyUserPermission(): void {
      if (this.password.length < 1) {
        this.$message.info("请输入密码")
        return
      }

      userAxios.modifyPermission(this.$store.state.userID, this.users[this.userIndex].user_id, this.permission, this.password)
        .then(response => {
          if (response.err) {
            throw response.err
          }

          this.$message.success("修改用户权限等级成功")
          this.users[this.userIndex].permission = this.permission

          this.showModifyPermissionDialog = false // close dialog only on success
        })
        .catch(err => {
          this.$message.error("修改用户权限等级失败")
          console.log("> modify user permission failed, error: ", err)
        })
        .finally(() => {
          this.password = ""
        })
    },

    beforeLockOrUnlockUser(index: number, wantLock: boolean): void {
      this.wantLock = wantLock
      this.userIndex = index

      this.showLockOrUnlockDialog = true
    },

    beforeModifyUserPermission(index: number): void {
      this.permission = 0
      this.userIndex = index

      this.showModifyPermissionDialog = true
    },
  }
})
</script>

<style lang="less">
.user-list {
  height: inherit;
  padding-right: 2rem;

  .ul-table {
    margin-top: 7rem;
    margin-bottom: 2rem;
    height: calc(100% - 13rem);
  }

  .ul-dialog {
    text-align: left;

    .uld-title {
      font-size: 3rem;
      font-weight: 600;
    }

    .el-form-item, .el-divider {
      margin-left: 20%;
    }

    .el-divider {
      width: 60%;
    }

    .el-input {
      width: 60%;
    }

    .el-form-item__label {
      font-size: 2rem;
    }
  }
}
</style>
