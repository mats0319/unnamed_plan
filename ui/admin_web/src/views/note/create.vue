<template>
  <div class="note-create">
    <el-form label-position="left" label-width="15%">
      <el-form-item label="主题">
        <el-input v-model="topic" placeholder="请输入笔记主题" />

        <el-popover trigger="hover" placement="top" :content="tips_Note_Topic">
          <i slot="reference" class="el-icon-warning-outline" />
        </el-popover>
      </el-form-item>

      <el-form-item label="是否公开">
        <el-checkbox v-model="isPublic">公开</el-checkbox>

        <el-popover trigger="hover" placement="top" :content="tips_IsPublic">
          <i slot="reference" class="el-icon-warning-outline" />
        </el-popover>
      </el-form-item>

      <el-form-item label="内容">
        <el-input
          v-model="content"
          type="textarea"
          :autosize="{ minRows: 3, maxRows: 5 }"
          resize="none"
          placeholder="请输入笔记内容"
        />
      </el-form-item>

      <el-form-item>
        <el-button type="info" plain @click="beforeCreateNote">记录笔记</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { tips_IsPublic, tips_Note_Topic } from "shared/ts/const";
import noteAxios from "shared/ts/axios_wrapper/note";

@Component
export default class CreateNote extends Vue {
  private topic = "";
  private isPublic = false;
  private content = "";

  // const
  private tips_IsPublic = tips_IsPublic;
  private tips_Note_Topic = tips_Note_Topic;

  private mounted() {
    // placeholder
  }

  private createNote(): void {
    noteAxios.create(this.$store.state.userID, this.topic, this.content, this.isPublic)
      .then(response => {
        if (response.data["hasError"]) {
          throw response.data["data"];
        }

        this.$message.success("记录笔记成功");

        this.topic = "";
        this.isPublic = false;
        this.content = "";
      })
      .catch(err => {
        this.$message.error("记录笔记失败");
        console.log("> create note failed.", err);
      })
  }

  private beforeCreateNote(): void {
    // null topic is valid
    if (this.content.length < 1) {
      this.$message.info("请输入笔记内容");
      return;
    }

    this.createNote();
  }
}
</script>

<style lang="scss">
.note-create {
  padding: 7vh 15vw;
  text-align: left;

  .el-form-item {
    margin: 5vh 0;
  }

  .el-form-item__label {
    font-size: 2rem;
  }

  .el-input, .el-checkbox, .el-textarea {
    width: 60%;
  }

  .el-popover__reference-wrapper {
    margin-left: 5vh;
  }
}
</style>
