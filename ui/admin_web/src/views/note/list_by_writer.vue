<template>
  <div class="note-list-by-writer">
    <el-table :data="notes" height="calc(80vh - 32px)" stripe highlight-current-row>
      <el-table-column label="主题" prop="topic" min-width="3" show-overflow-tooltip />
      <el-table-column label="内容" prop="content" min-width="5" show-overflow-tooltip />

      <el-table-column label="是否公开" min-width="1" show-overflow-tooltip>
        <template slot-scope="scope">
          {{ scope.row.isPublic | displayIsPublic }}
        </template>
      </el-table-column>

      <el-table-column label="修改时间" min-width="3" show-overflow-tooltip>
        <template slot-scope="scope">
          {{ scope.row.updateTime | displayTime }}
        </template>
      </el-table-column>

      <el-table-column label="上传时间" min-width="3" show-overflow-tooltip>
        <template slot-scope="scope">
          {{ scope.row.createdTime | displayTime }}
        </template>
      </el-table-column>
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
      @current-change="list"
    />

    <el-dialog
      class="nlbw-dialog"
      :visible.sync="modifyDialogController"
      title="修改笔记"
      :before-close="resetModifyDialogData"
    >
      <div class="nlbwd-content">
        <el-form label-position="left" label-width="20%">

          <el-form-item label="笔记ID">{{ noteID }}</el-form-item>
          <el-form-item label="原主题">{{ oldTopic }}</el-form-item>

          <hr />

          <el-form-item label="主题">
            <el-input v-model="topic" placeholder="主题" />

            <el-popover trigger="hover" placement="top" :content="tips_Note_Topic">
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

          <el-form-item label="是否公开">
            <el-checkbox v-model="isPublic">公开</el-checkbox>

            <el-popover trigger="hover" placement="top" :content="tips_IsPublic">
              <i slot="reference" class="el-icon-warning-outline" />
            </el-popover>
          </el-form-item>

          <hr />

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
      class="nlbw-dialog"
      :visible.sync="deleteDialogController"
      title="删除笔记"
      :before-close="resetDeleteDialogData"
    >
      <div class="nlbwd-content">
        <el-form label-position="left" label-width="20%">

          <el-form-item label="笔记ID">{{ noteID }}</el-form-item>
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
import { Note } from "shared/ts/data";
import { tips_IsPublic, tips_Note_Topic } from "shared/ts/const";
import noteAxios from "shared/ts/axios_wrapper/note";

@Component
export default class ListNoteByWriter extends Vue {
  private notes: Array<Note> = new Array<Note>();

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
  private tips_Note_Topic = tips_Note_Topic;

  private mounted() {
    this.list();
  }

  private list(currPage?: number): void {
    this.total = 0;
    this.notes = [];

    noteAxios.listByWriter(this.$store.state.userID, this.pageSize, currPage ? currPage : 1)
      .then(response => {
        if (response.data["hasError"]) {
          throw response.data["data"];
        }

        if (currPage) {
          this.pageNum = currPage;
        }

        const payload = JSON.parse(response.data["data"] as string);

        this.total = payload.total;
        for (let i = 0; i < payload.notes.length; i++) {
          const item = payload.notes[i];

          this.notes.push({
            noteID: item.noteID,
            writeBy: item.writeBy,
            topic: item.topic,
            content: item.content,
            isPublic: item.isPublic,
            updateTime: item.updateTime,
            createdTime: item.createdTime
          });
        }
      })
      .catch(err => {
        this.$message.error("获取当前用户记录的笔记列表失败");
        console.log("> get note by writer failed.", err);
      });
  }

  private modifyNote(): void {
    if (this.topic == this.oldTopic && this.content == this.oldContent && this.isPublic == this.oldIsPublic) {
      this.$message.info("当前未执行任何修改");
      return;
    }

    noteAxios.modify(this.$store.state.userID, this.noteID, this.password, this.topic, this.content, this.isPublic)
      .then(response => {
        if (response.data["hasError"]) {
          throw response.data["data"];
        }

        this.$message.success("修改笔记成功");
        this.modifyDialogController = false;

        this.list(this.pageNum);
      })
      .catch(err => {
        this.$message.error("修改笔记失败");
        console.log("> modify note failed.", err);
      })
      .finally(() => {
        this.password = "";
      })
  }

  private deleteNote(): void {
    noteAxios.delete(this.$store.state.userID, this.noteID, this.password)
      .then(response => {
        if (response.data["hasError"]) {
          throw response.data["data"];
        }

        this.$message.success("删除笔记成功");
        this.deleteDialogController = false;

        this.list(); // 防止删除一页的最后一条时，再次查询该页数据导致异常，故此处重新查询
      })
      .catch(err => {
        this.$message.error("删除笔记失败");
        console.log("> delete note failed.", err);
      })
      .finally(() => {
        this.password = "";
      })
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
.note-list-by-writer {
  .nlbw-dialog {
    text-align: left;

    .nlbwd-content {
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
