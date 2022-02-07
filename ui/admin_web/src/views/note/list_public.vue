<template>
  <div>
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
    </el-table>

    <el-pagination
      :total="total"
      :page-size="pageSize"
      :current-page="pageNum"
      layout="prev, pager, next, ->, total"
      @current-change="list"
    />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { Note } from "shared/ts/data";
import noteAxios from "shared/ts/axios_wrapper/note";

@Component
export default class ListPublicNote extends Vue {
  private notes: Array<Note> = new Array<Note>();

  private total = 0;
  private pageSize = 10;
  private pageNum = 1;

  private mounted() {
    this.list();
  }

  private list(currPage?: number): void {
    this.total = 0;
    this.notes = [];

    noteAxios.listPublic(this.$store.state.userID, this.pageSize, currPage ? currPage : 1)
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
}
</script>
