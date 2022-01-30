<template>
  <div class="task-create">
    <el-form label-position="left" label-width="15%">
      <el-form-item label="任务名">
        <el-input v-model="taskName" placeholder="请输入任务名" />

        <el-popover trigger="hover" placement="top" :content="tips_Task_Name">
          <i slot="reference" class="el-icon-warning-outline" />
        </el-popover>
      </el-form-item>

      <el-form-item label="描述">
        <el-input
          v-model="description"
          type="textarea"
          :autosize="{ minRows: 3, maxRows: 5 }"
          resize="none"
          placeholder="请输入任务描述"
        />
      </el-form-item>

      <el-form-item label="前置任务">
        <el-select v-model="preTasks" multiple placeholder="请选择前置任务（未实现）" disabled />
      </el-form-item>

      <el-form-item>
        <el-button type="info" plain @click="beforePostTask">发布任务</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { tips_Task_Name } from "shared/ts/const";
import taskAxios from "shared/ts/axios_wrapper/task";

@Component
export default class TaskCreate extends Vue {
  private taskName = "";
  private description = "";
  private preTasks: Array<string> = []; // todo: unimplemented

  // const
  private tips_Task_Name = tips_Task_Name;

  private mounted() {
    // placeholder
  }

  private postTask(): void {
    taskAxios.create(this.$store.state.userID, this.taskName, this.description, this.preTasks)
      .then(response => {
        if (response.data["hasError"]) {
          throw response.data["data"];
        }

        this.$message.success("发布任务成功");

        this.taskName = "";
        this.description = "";
        this.preTasks = [];
      })
      .catch(err => {
        this.$message.error("发布任务失败");
        console.log("> create task failed.", err);
      })
  }

  private beforePostTask(): void {
    if (this.taskName.length < 1) {
      this.$message.info("请输入任务名");
      return;
    }

    this.postTask();
  }
}
</script>

<style lang="scss">
.task-create {
  padding: 7vh 15vw;
  text-align: left;

  .el-form-item {
    margin: 5vh 0;
  }

  .el-form-item__label {
    font-size: 2rem;
  }

  .el-input, .el-select, .el-textarea {
    width: 60%;
  }

  .el-popover__reference-wrapper {
    margin-left: 5vh;
  }
}
</style>
