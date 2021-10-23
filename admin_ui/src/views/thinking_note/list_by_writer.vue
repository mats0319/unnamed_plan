<template>
  <div class="thinking-note-list-by-writer">
    <el-table :data="notes" height="calc(80vh - 32px)" stripe highlight-current-row>
      <el-table-column label="主题" prop="topic" min-width="3" show-overflow-tooltip />
      <el-table-column label="内容" prop="content" min-width="5" show-overflow-tooltip />
      <el-table-column label="是否公开" prop="isPublicDisplay" min-width="1" show-overflow-tooltip />
      <el-table-column label="修改时间" prop="updateTimeDisplay" min-width="2" show-overflow-tooltip />
      <el-table-column label="上传时间" prop="createdTimeDisplay" min-width="2" show-overflow-tooltip />
      <el-table-column label="操作" min-width="2">
        <template slot-scope="scope">
          <el-button
            type="info"
            size="mini"
            plain
            @click="beforeModifyNote(scope.row.noteID, scope.row.topic, scope.row.content, scope.row.isPublic)"
          >
            修改
          </el-button>

          <el-button
            type="info"
            size="mini"
            plain
            @click="beforeDeleteNote(scope.row.noteID, scope.row.topic)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      :total="total"
      :page-size="pageSize"
      :current-page="pageNum"
      layout="prev, pager, next, ->, total"
      @current-change="listByWriter"
    />

    <el-dialog
      class="tnlbw-dialog"
      :visible.sync="modifyDialogController"
      title="修改随想"
      :before-close="resetModifyDialogData"
    >
      <div class="tnlbwd-content">
        <el-form label-position="left" label-width="20%">

          <el-form-item label="随想ID">{{ noteID }}</el-form-item>
          <el-form-item label="原主题">{{ oldTopic }}</el-form-item>

          <hr />

          <el-form-item label="主题">
            <el-input v-model="topic" placeholder="主题" />

            <el-popover trigger="hover" placement="top" :content="tips_ThinkingNote_Topic">
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

          <el-form-item label="是否公开">
            <el-checkbox v-model="isPublic">公开</el-checkbox>

            <el-popover trigger="hover" placement="top" :content="tips_IsPublic">
              <i slot="reference" class="el-icon-warning-outline" />
            </el-popover>
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
        <el-button type="info" plain @click="modifyNote">修改</el-button>
      </div>
    </el-dialog>

    <el-dialog
      class="tnlbw-dialog"
      :visible.sync="deleteDialogController"
      title="删除随想"
      :before-close="resetDeleteDialogData"
    >
      <div class="tnlbwd-content">
        <el-form label-position="left" label-width="20%">

          <el-form-item label="随想ID">{{ noteID }}</el-form-item>
          <el-form-item label="主题">{{ topic }}</el-form-item>

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
        <el-button type="info" plain @click="deleteNote">删除</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { displayIsPublic, displayTime, ThinkingNote } from "@/ts/data";
import axios from "axios";
import { tips_IsPublic, tips_ThinkingNote_Topic } from "@/ts/const";
import { calcSHA256 } from "@/ts/utils";

@Component
export default class ListThinkingNoteByWriter extends Vue {
  private notes: Array<ThinkingNote> = new Array<ThinkingNote>();

  private total = 0;
  private pageSize = 10;
  private pageNum = 1;

  private modifyDialogController = false;
  private noteID = "";
  private oldTopic = "";
  private oldContent = "";
  private oldIsPublic = false;

  private deleteDialogController = false;
  private topic = "";
  private content = "";
  private isPublic = false;
  private password = "";

  // const
  private tips_IsPublic = tips_IsPublic;
  private tips_ThinkingNote_Topic = tips_ThinkingNote_Topic;

  private mounted() {
    this.listByWriter();
  }

  private listByWriter(currPage?: number): void {
    this.total = 0;
    this.notes = [];

    let data: FormData = new FormData();
    data.append("operatorID", this.$store.state.userID);
    data.append("pageSize", this.pageSize.toString());
    data.append("pageNum", currPage ? currPage.toString() : "1");

    axios.post(process.env.VUE_APP_thinking_note_list_by_writer_url, data)
      .then(response => {
        if (response.data.hasError) {
          throw response.data.data;
        }

        if (currPage) {
          this.pageNum = currPage;
        }

        const payload = JSON.parse(response.data.data as string);

        this.total = payload.total;
        for (let i = 0; i < payload.notes.length; i++) {
          const item = payload.notes[i];

          this.notes.push({
            noteID: item.noteID,
            writeBy: item.writeBy,
            topic: item.topic,
            content: item.content,
            isPublic: item.isPublic,
            isPublicDisplay: displayIsPublic(item.isPublic),
            updateTime: item.updateTime,
            updateTimeDisplay: displayTime(item.updateTime),
            createdTime: item.createdTime,
            createdTimeDisplay: displayTime(item.createdTime)
          });
        }
      })
      .catch(err => {
        this.$message.error("获取当前用户记录的随想列表失败，错误：" + err);
      });
  }

  private modifyNote(): void {
    if (!this.isValidModifyParams()) {
      this.$message.info("当前未执行任何修改");
      return;
    }

    const pwd = calcSHA256(this.password);
    this.password = "";

    let data: FormData = new FormData();
    data.append("operatorID", this.$store.state.userID);
    data.append("noteID", this.noteID);
    data.append("password", pwd);
    data.append("topic", this.topic);
    data.append("content", this.content);
    data.append("isPublic", this.isPublic.toString());

    axios.post(process.env.VUE_APP_cloud_file_modify_url, data)
      .then(response => {
        if (response.data.hasError) {
          throw response.data.data;
        }

        const payload = JSON.parse(response.data.data as string);
        if (payload.isSuccess) {
          this.$message.success("修改随想成功");

          this.listByWriter(this.pageNum);
        } else {
          this.$message.error("修改随想失败");
        }
      })
      .catch(err => {
        this.$message.error("修改随想失败，错误：" + err);
      })
  }

  private deleteNote(): void {
    const pwd = calcSHA256(this.password);
    this.password = "";

    let data: FormData = new FormData();
    data.append("operatorID", this.$store.state.userID);
    data.append("noteID", this.noteID);
    data.append("password", pwd);

    axios.post(process.env.VUE_APP_thinking_note_delete_url, data)
      .then(response => {
        if (response.data.hasError) {
          throw response.data.data;
        }

        const payload = JSON.parse(response.data.data as string);
        if (payload.isSuccess) {
          this.$message.success("删除随想成功");

          this.listByWriter(this.pageNum);
        } else {
          this.$message.error("删除随想失败");
        }
      })
      .catch(err => {
        this.$message.error("删除随想失败，错误：" + err);
      })
  }

  private isValidModifyParams(): boolean {
    return this.topic != this.oldTopic || this.content != this.oldContent || this.isPublic != this.oldIsPublic
  }

  private beforeModifyNote(noteID: string, topic: string, content: string, isPublic: boolean): void {
    this.noteID = noteID;
    this.oldTopic = topic;
    this.topic = topic;
    this.oldContent = content;
    this.content = content;
    this.oldIsPublic = isPublic;
    this.isPublic = isPublic;

    this.modifyDialogController = true;
  }

  private beforeDeleteNote(noteID: string, topic: string): void {
    this.noteID = noteID;
    this.topic = topic;

    this.deleteDialogController = true;
  }

  private resetModifyDialogData(): void {
    this.noteID = "";
    this.oldTopic = "";
    this.topic = "";
    this.isPublic = false;

    this.modifyDialogController = false;
  }

  private resetDeleteDialogData(): void {
    this.noteID = "";
    this.topic = "";
    this.password = ""

    this.deleteDialogController = false;
  }
}
</script>

<style lang="scss">
.thinking-note-list-by-writer {
  .tnlbw-dialog {
    text-align: left;

    .tnlbwd-content {
      padding: 2vh 15%;
    }

    .el-input, .el-checkbox, textarea {
      width: 60%;
    }

    .el-popover__reference-wrapper {
      margin-left: 5vh;
    }
  }
}
</style>
