<template>
  <div class="thinking-note-create">
    <el-form label-position="left" label-width="15%">
      <el-form-item label="主题">
        <el-input v-model="topic" placeholder="请输入随想主题" />

        <el-popover trigger="hover" placement="top" :content="tips_ThinkingNote_Topic">
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
          placeholder="请输入随想内容"
        />
      </el-form-item>

      <el-form-item>
        <el-button type="info" plain @click="beforeCreateThinkingNote">记录随想</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { tips_IsPublic, tips_ThinkingNote_Topic } from "shared_ui/ts/const";
import thinkingNoteAxios from "shared_ui/ts/axios_wrapper/thinking_note";

@Component
export default class CreateThinkingNote extends Vue {
  private topic = "";
  private isPublic = false;
  private content = "";

  // const
  private tips_IsPublic = tips_IsPublic;
  private tips_ThinkingNote_Topic = tips_ThinkingNote_Topic;

  private mounted() {
    // placeholder
  }

  private createThinkingNote(): void {
    thinkingNoteAxios.create(this.$store.state.userID, this.topic, this.content, this.isPublic)
      .then(response => {
        if (response.data["hasError"]) {
          throw response.data["data"];
        }

        this.$message.success("记录随想成功");
      })
      .catch(err => {
        this.$message.error("记录随想失败，错误：" + err);
      })
  }

  private beforeCreateThinkingNote(): void {
    // null topic is valid
    if (this.content.length < 1) {
      this.$message.info("请输入随想内容");
      return;
    }

    this.createThinkingNote();
  }
}
</script>

<style lang="scss">
.thinking-note-create {
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
