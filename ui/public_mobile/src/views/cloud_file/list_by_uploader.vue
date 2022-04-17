<template>
  <div class="list-cloud-file-by-uploader">
    <el-table :data="cloudFiles" height="calc(94vh - 20rem - 32px)" stripe highlight-current-row>
      <el-table-column label="文件名" min-width="3">
        <template slot-scope="scope">
          <div class="cflbu-file-name">
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
import { Component, Vue } from "vue-property-decorator";
import { CloudFile } from "shared/ts/data";
import { generateCloudFileURL } from "shared/ts/utils";
import cloudFileAxios from "shared/ts/axios_wrapper/cloud_file";

@Component
export default class ListCloudFileByUploader extends Vue {
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

    cloudFileAxios.listByUploader(this.$store.state.userID, this.pageSize, currPage ? currPage : 1)
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
        this.$message.error("获取当前用户上传的文件列表失败，错误：" + err);
      })
  }
}
</script>

<style lang="less">
.list-cloud-file-by-uploader {
  padding: 3vh 0;

  .cflbu-file-name {
    a {
      color: darkgray;
    }
  }

  .cflbu-file-name:hover {
    a {
      color: lightgray;
    }
  }
}
</style>
