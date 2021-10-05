<template>
  <div class="user-list">
    <el-table :data="users" height="calc(80vh - 32px)" stripe highlight-current-row>
      <el-table-column label="用户名" prop="userName" min-width="2" show-overflow-tooltip />
      <el-table-column label="昵称" prop="nickname" min-width="2" show-overflow-tooltip />
      <el-table-column label="锁定状态" prop="isLockedDisplay" min-width="2" show-overflow-tooltip />
      <el-table-column label="权限等级" prop="permission" min-width="2" show-overflow-tooltip />
      <el-table-column label="操作" min-width="3">
        <template slot-scope="scope">
          <div class="ul-operate">
            <div v-show="$store.state.permission >= 6" class="ulo-item">
              <el-button
                v-show="!scope.row.isLocked"
                type="info"
                size="mini"
                plain
                @click="lockOrUnlockUser(scope.row.userID, true)"
              >
                锁定
              </el-button>

              <el-button
                v-show="scope.row.isLocked"
                type="info"
                size="mini"
                plain
                @click="lockOrUnlockUser(scope.row.userID, false)"
              >
                解锁
              </el-button>
            </div>

            <div v-show="$store.state.permission >= 8" class="ulo-item">
              <el-button
                type="info"
                size="mini"
                plain
                @click="beforeModifyUserPermission(scope.row.userID)"
              >
                修改权限
              </el-button>
            </div>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      :total="total"
      :page-size="pageSize"
      :current-page="pageNum"
      layout="prev, pager, next, ->, total"
      @current-change="listUser"
    />

    <el-dialog
      class="ul-modify-permission-dialog"
      :visible.sync="modifyPermissionDialogController"
      title="修改权限等级"
      :before-close="resetDialogData"
    >
      <div>
        目标用户ID&#58;&nbsp;{{ userID }}<br />

        <el-select
          class="ulmpd-permission"
          v-model="permissionStr"
          placeholder="请选择目标用户新的权限等级"
          clearable
        >
          <el-option
            v-for="item in 7"
            :key="item"
            :label="item"
            :value="item.toString()"
          />
        </el-select>
      </div>

      <div slot="footer">
        <el-button type="info" plain @click="modifyUserPermission">修改</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { User, displayUserIsLocked } from "@/ts/data";
import axios from "axios";

@Component
export default class UserList extends Vue {
  private users: Array<User> = new Array<User>();

  private total = 0;
  private pageSize = 10;
  private pageNum = 1;

  private modifyPermissionDialogController = false;
  private userID = "";
  private permissionStr = "";

  private mounted() {
    this.listUser();
  }

  private listUser(currPage?: number): void {
    this.total = 0;
    this.users = [];

    let data: FormData = new FormData();
    data.append("operatorID", this.$store.state.userID);
    data.append("pageSize", this.pageSize.toString());
    data.append("pageNum", currPage ? currPage.toString() : "1");

    axios.post(process.env.VUE_APP_user_list_url, data)
      .then(response => {
        if (response.data.hasError) {
          throw response.data.data;
        }

        if (currPage) {
          this.pageNum = currPage;
        }

        const payload = JSON.parse(response.data.data as string);
        this.total = payload.total;
        for (let i = 0; i < payload.users.length; i++) {
          const item = payload.users[i];

          this.users.push({
            userID: item.userID,
            userName: item.userName,
            nickname: item.nickname,
            isLocked: item.isLocked,
            isLockedDisplay: displayUserIsLocked(item.isLocked),
            permission: item.permission,
            createdBy: item.createdBy
          });
        }
      })
      .catch(err => {
        console.log("list user failed, error:", err);
      });
  }

  private lockOrUnlockUser(userID: string, wantLock: boolean): void {
    let data: FormData = new FormData();
    data.append("operatorID", this.$store.state.userID);
    data.append("userID", userID);

    let url = wantLock ? process.env.VUE_APP_user_lock_url : process.env.VUE_APP_user_unlock_url;

    axios.post(url, data)
      .then(response => {
        if (response.data.hasError) {
          throw response.data.data;
        }

        const payload = JSON.parse(response.data.data as string);
        if (payload.isSuccess) {
          this.$message.success(wantLock ? "锁定用户成功" : "解锁用户成功");

          this.listUser(this.pageNum);
        } else {
          this.$message.error(wantLock ? "锁定用户失败" : "解锁用户失败");
        }
      })
      .catch(err => {
        console.log("lock or unlock failed, want lock:", wantLock, ", error:", err);
      });
  }

  private modifyUserPermission(): void {
    let data: FormData = new FormData();
    data.append("operatorID", this.$store.state.userID);
    data.append("userID", this.userID);
    data.append("permission", this.permissionStr);

    axios.post(process.env.VUE_APP_user_modify_permission_url, data)
      .then(response => {
        if (response.data.hasError) {
          throw response.data.data;
        }

        const payload = JSON.parse(response.data.data as string);
        if (payload.isSuccess) {
          this.$message.success("修改用户权限成功");

          this.listUser(this.pageNum);
        } else {
          this.$message.error("修改用户权限失败");
        }
      })
      .catch(err => {
        console.log("modify user permission failed, error:", err);
      })
  }

  private beforeModifyUserPermission(userID: string): void {
    this.userID = userID;
    this.modifyPermissionDialogController = true;
  }

  private resetDialogData(): void {
    this.userID = "";
    this.permissionStr = "";
  }
}
</script>

<style lang="scss">
.user-list {
  .ul-operate {
    display: flex;

    .ulo-item {
      margin: auto 1vw;

      .el-button + .el-button {
        margin-left: 0;
      }
    }
  }

  .ul-modify-permission-dialog {
    text-align: left;

    .ulmpd-permission {
      margin-top: 5vh;

      .el-input {
        width: 150%;
      }
    }

    .el-dialog__title {
      font-size: 3rem;
      line-height: 5rem;
    }

    .el-dialog__body {
      padding: 5vh 10%;
      font-size: 2rem;
    }
  }
}
</style>
