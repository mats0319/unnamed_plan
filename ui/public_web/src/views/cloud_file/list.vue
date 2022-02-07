<template>
  <div class="list-cloud-file">
    <el-table :data="cloudFiles" height="calc(94vh - 20rem - 32px)" stripe highlight-current-row>
      <el-table-column label="文件名" min-width="2">
        <template slot-scope="scope">
          <div class="lcf-file-name">
            <a :href="scope.row.fileURL" target="_blank">{{ scope.row.fileName }}</a>
          </div>
        </template>
      </el-table-column>

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
import { Component, Vue, Watch } from "vue-property-decorator";
import { CloudFile } from "shared/ts/data";
import { generateCloudFileURL } from "shared/ts/utils";
import cloudFileAxios from "shared/ts/axios_wrapper/cloud_file";

@Component
export default class ListCloudFile extends Vue {
  private cloudFiles: Array<CloudFile> = new Array<CloudFile>();

  private total = 0;
  private pageSize = 10;
  private pageNum = 1;

  private mounted() {
    this.list();
  }

  private list(currPage?: number): void {
    this.total = 0;
    this.cloudFiles = [];

    let axiosResPromise;

    switch (this.$store.state.cloudFilePageType) {
      case "0":
        axiosResPromise = cloudFileAxios.listByUploader(this.$store.state.userID, this.pageSize, currPage ? currPage : 1)
        break;
      case "1":
        axiosResPromise = cloudFileAxios.listPublic(this.$store.state.userID, this.pageSize, currPage ? currPage : 1)
        break;
    }

    axiosResPromise
      .then(response => {
        if (response.data["hasError"]) {
          throw response.data["data"];
        }

        if (currPage) {
          this.pageNum = currPage;
        }

        const payload = JSON.parse(response.data["data"] as string);

        this.total = payload.total;
        for (let i = 0; i < payload.files.length; i++) {
          const item = payload.files[i];

          this.cloudFiles.push({
            fileID: item.fileID,
            fileName: item.fileName,
            fileURL: generateCloudFileURL(item.fileURL),
            isPublic: item.isPublic,
            updateTime: item.updateTime,
            createdTime: item.createdTime
          });
        }
      })
      .catch(err => {
        let errMsg = "";
        switch (this.$store.state.cloudFilePageType) {
          case "0":
            errMsg = "获取当前用户上传的文件列表失败"
            break;
          case "1":
            errMsg = "获取公开的文件列表失败"
            break;
        }

        this.$message.error(errMsg);
        console.log("> get cloud file failed.", err)
      })
  }

  @Watch("$store.state.cloudFilePageType")
  private watchPageType(): void {
    this.list();
  }
}
</script>

<style lang="scss">
.list-cloud-file {
  padding: 3vh 10vw;

  .lcf-file-name {
    a {
      color: darkgray;
    }
  }

  .lcf-file-name:hover {
    a {
      color: lightgray;
    }
  }
}
</style>
