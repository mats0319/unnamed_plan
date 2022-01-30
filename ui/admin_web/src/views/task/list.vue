<template>
  <div class="task-list">
    <el-table :data="tasks" height="calc(80vh - 32px)" stripe highlight-current-row>
      <el-table-column label="任务名" prop="taskName" min-width="2" show-overflow-tooltip />
      <el-table-column label="描述" prop="description" min-width="5" show-overflow-tooltip />

      <el-table-column label="前置任务" min-width="5" show-overflow-tooltip>
        <template slot-scope="scope">
          {{ scope.row.preTasks }}
        </template>
      </el-table-column>

      <el-table-column label="状态" min-width="2" show-overflow-tooltip>
        <template slot-scope="scope">
          {{ scope.row.status | displayTaskStatus }}
        </template>
      </el-table-column>

      <el-table-column label="修改时间" min-width="3" show-overflow-tooltip>
        <template slot-scope="scope">
          {{ scope.row.updateTime | displayTime }}
        </template>
      </el-table-column>

      <el-table-column label="发布时间" min-width="3" show-overflow-tooltip>
        <template slot-scope="scope">
          {{ scope.row.createdTime | displayTime }}
        </template>
      </el-table-column>

      <el-table-column label="操作" min-width="2">
        <template slot-scope="scope">
          <el-button type="info" size="mini" plain @click="beforeModifyTask(scope.row)">修改</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
      class="tl-dialog"
      :visible.sync="modifyDialogController"
      title="修改任务"
      :before-close="resetModifyDialog"
    >
      <div class="tld-content">
        <el-form label-position="left" label-width="20%">

          <el-form-item label="任务ID">{{ selectedTask.taskID }}</el-form-item>
          <el-form-item label="原任务名">{{ selectedTask.taskName }}</el-form-item>

          <hr />

          <el-form-item label="任务名">
            <el-input v-model="modifyingTask.taskName" placeholder="任务名" />

            <el-popover trigger="hover" placement="top" :content="tips_Task_Name">
              <i slot="reference" class="el-icon-warning-outline" />
            </el-popover>
          </el-form-item>

          <el-form-item label="描述">
            <el-input
              v-model="modifyingTask.description"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 5 }"
              resize="none"
              placeholder="请输入任务描述"
            />
          </el-form-item>

          <el-form-item label="前置任务">
            {{ modifyingTask.preTaskIDs }}
          </el-form-item>

          <el-form-item label="状态">
            <el-select v-model="modifyingTask.status" placeholder="请选择任务状态" clearable>
              <el-option
                v-for="(item, index) in taskStatus"
                :key="index"
                :label="item"
                :value="index"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="密码">
            <el-input
              v-model="password"
              type="password"
              placeholder="请输入密码"
              show-password
              clearable
            />
          </el-form-item>
        </el-form>
      </div>

      <div slot="footer">
        <el-button type="info" plain @click="modifyTask">修改</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { Task, newTask, deepCopyTask } from "shared/ts/data";
import { tips_Task_Name, taskStatus } from "shared/ts/const";
import { compareOnStringSliceNotStrict } from "shared/ts/utils";
import taskAxios from "shared/ts/axios_wrapper/task";
import {displayPreTasks} from "@/ts/utils";

@Component
export default class TaskList extends Vue {
  private tasks: Array<Task> = new Array<Task>();
  private selectedTask = newTask();

  private modifyDialogController = false;
  private modifyingTask = newTask();
  private password = "";

  private total = 0;

  // const
  private tips_Task_Name = tips_Task_Name;
  private taskStatus = taskStatus;

  private mounted() {
    this.list();
  }

  private list(): void {
    this.total = 0;
    this.tasks = [];

    taskAxios.list(this.$store.state.userID)
      .then(response => {
        if (response.data["hasError"]) {
          throw response.data["data"];
        }

        const payload = JSON.parse(response.data["data"] as string);

        this.total = payload.total;
        for (let i = 0; i < payload.tasks.length; i++) {
          const item = payload.tasks[i];

          this.tasks.push({
            taskID: item.taskID,
            taskName: item.taskName,
            description: item.description,
            preTaskIDs: item.preTaskIDs,
            preTasks: displayPreTasks(item.preTaskIDs, payload.tasks as Array<Task>),
            status: item.status,
            updateTime: item.updateTime,
            createdTime: item.createdTime,
          });
        }
      })
      .catch(err => {
        this.$message.error("获取当前用户发布的任务列表失败");
        console.log("> get task by poster failed.", err);
      });
  }

  private modifyTask(): void {
    if (this.selectedTask.taskName == this.modifyingTask.taskName &&
      this.selectedTask.description == this.modifyingTask.description &&
      compareOnStringSliceNotStrict(this.selectedTask.preTaskIDs, this.modifyingTask.preTaskIDs) &&
      this.selectedTask.status == this.modifyingTask.status) {
      this.$message.info("当前未执行任何修改");
      return;
    }

    taskAxios.modify(this.$store.state.userID,
      this.modifyingTask.taskID,
      this.password,
      this.modifyingTask.taskName,
      this.modifyingTask.description,
      this.modifyingTask.preTaskIDs,
      this.modifyingTask.status,
    )
      .then(response => {
        if (response.data["hasError"]) {
          throw response.data["data"];
        }

        this.$message.success("修改任务成功");
        this.modifyDialogController = false;

        this.list();
      })
      .catch(err => {
        this.$message.error("修改任务失败");
        console.log("> modify task failed.", err);
      })
      .finally(() => {
        this.password = "";
      })
  }

  private beforeModifyTask(item: Task): void {
    this.selectedTask = deepCopyTask(item);
    this.modifyingTask = deepCopyTask(item);

    this.modifyDialogController = true;
  }

  private resetModifyDialog(): void {
    this.password = "";
  }
}
</script>

<style lang="scss">
.task-list {
  .el-table, .cell.el-tooltip {
    white-space: pre-wrap;
  }

  .tl-dialog {
    text-align: left;

    .tld-content {
      padding: 2vh 15%;
    }

    .el-input, .el-select, textarea {
      width: 60%;
    }

    .el-popover__reference-wrapper {
      margin-left: 5vh;
    }
  }
}
</style>
