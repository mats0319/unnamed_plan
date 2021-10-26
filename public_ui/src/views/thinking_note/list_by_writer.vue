<template>
  <div class="list-thinking-note-by-writer">
    <el-timeline class="ltnbw-content">
      <el-timeline-item
        v-for="item in notes"
        :key="item['noteID']"
        placement="top"
        hide-timestamp
      >
        <el-card>
          <el-form label-position="left" label-width="20%">
            <el-form-item label="主题">{{ item["topic"] }}</el-form-item>
            <el-form-item label="内容">{{ item["content"] }}</el-form-item>
            <el-form-item label="是否公开">{{ item["isPublic"] | displayIsPublic }}</el-form-item>
            <el-form-item label="更新时间">{{ item["updateTime"] | displayTime }}</el-form-item>
            <el-form-item label="上传时间">{{ item["createdTime"] | displayTime }}</el-form-item>
          </el-form>
        </el-card>
      </el-timeline-item>
    </el-timeline>

    <div class="ltnbw-bottom">
      <div class="ltnbwb-data">
        共&nbsp;{{ total }}&nbsp;条<br />
        每次加载&nbsp;{{ pageSize }}&nbsp;条<br />
      </div>

      <div class="ltnbwb-more">
        <el-button v-show="(pageNum-1) * pageSize < total" type="text" plain @click="list">
          加载更多
        </el-button>

        <span v-show="(pageNum-1) * pageSize >= total" class="ltnbwbm-no-more">没有更多了</span>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { ThinkingNote } from "shared_ui/ts/data";
import thinkingNoteAxios from "shared_ui/ts/axios_wrapper/thinking_note";

@Component
export default class ListThinkingNoteByWriter extends Vue {
  private notes: Array<ThinkingNote> = new Array<ThinkingNote>();

  private total = 0;
  private pageSize = 5;
  private pageNum = 1;

  private mounted() {
    this.list();
  }

  private list(): void {
    thinkingNoteAxios.listByWriter(this.$store.state.userID, this.pageSize, this.pageNum)
      .then(response => {
        if (response.data["hasError"]) {
          throw response.data["data"];
        }

        this.pageNum++;

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
          })
        }
      })
      .catch(err => {
        this.$message.error("获取当前用户记录的随想列表失败，错误：" + err);
      })
  }
}
</script>

<style lang="scss">
.list-thinking-note-by-writer {
  display: flex;
  height: calc(100vh - 20rem);
  overflow-y: auto;

  .ltnbw-content {
    width: 60vw;
    padding-left: 20vw;
    padding-top: 3vh;
  }

  .ltnbw-bottom {
    position: absolute;
    right: 0;
    bottom: 10rem;

    width: 20vw;
    height: calc(8vh + 10rem);
    margin-top: auto;

    font-size: 2rem;

    .ltnbwb-data {
      line-height: 3rem;
      padding-bottom: 5vh;
    }

    .ltnbwb-more {
      .ltnbwbm-no-more {
        color: darkgray;
      }

      .el-button {
        border: 0;
      }
    }
  }

  .el-form-item__label {
    font-size: 2rem;
    font-weight: 600;
  }

  .el-form-item__content {
    font-size: 1.8rem;
    text-align: left;
    white-space: pre-wrap;
  }
}
</style>
