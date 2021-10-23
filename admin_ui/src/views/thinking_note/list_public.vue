<template>
  <div class="thinking-note-list-public">
    <el-table :data="notes" height="calc(80vh - 32px)" stripe highlight-current-row>
      <el-table-column label="主题" prop="topic" min-width="3" show-overflow-tooltip />
      <el-table-column label="内容" prop="content" min-width="5" show-overflow-tooltip />
      <el-table-column label="是否公开" prop="isPublicDisplay" min-width="1" show-overflow-tooltip />
      <el-table-column label="修改时间" prop="updateTimeDisplay" min-width="2" show-overflow-tooltip />
      <el-table-column label="上传时间" prop="createdTimeDisplay" min-width="2" show-overflow-tooltip />
    </el-table>

    <el-pagination
      :total="total"
      :page-size="pageSize"
      :current-page="pageNum"
      layout="prev, pager, next, ->, total"
      @current-change="listPublic"
    />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { displayIsPublic, displayTime, ThinkingNote } from "@/ts/data";
import axios from "_axios@0.21.4@axios";

@Component
export default class ListPublicThinkingNote extends Vue {
  private notes: Array<ThinkingNote> = new Array<ThinkingNote>();

  private total = 0;
  private pageSize = 10;
  private pageNum = 1;

  private mounted() {
    this.listPublic();
  }

  private listPublic(currPage?: number): void {
    this.total = 0;
    this.notes = [];

    let data: FormData = new FormData();
    data.append("operatorID", this.$store.state.userID);
    data.append("pageSize", this.pageSize.toString());
    data.append("pageNum", currPage ? currPage.toString() : "1");

    axios.post(process.env.VUE_APP_thinking_note_list_public_url, data)
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
}
</script>

<style lang="scss">
.thinking-note-list-public {

}
</style>
