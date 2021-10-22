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
      <el-table-column label="是否公开" prop="isPublicDisplay" min-width="2" show-overflow-tooltip />
      <el-table-column label="上传时间" prop="createdTimeDisplay" min-width="3" show-overflow-tooltip />
    </el-table>

    <el-pagination
      :total="total"
      :page-size="pageSize"
      :current-page="pageNum"
      layout="prev, pager, next, ->, total"
      @current-change="listCloudFileByUploader"
    />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { CloudFile, displayCloudFileIsPublic, displayCloudFileTime, generateCloudFileURL } from "@/ts/data";
import axios from "axios";

@Component
export default class ListCloudFileByUploader extends Vue {
  private cloudFiles: Array<CloudFile> = new Array<CloudFile>();

  private total = 0;
  private pageSize = 10;
  private pageNum = 1;

  private mounted() {
    this.listCloudFileByUploader();
  }

  private listCloudFileByUploader(currPage?: number): void {
    this.total = 0;
    this.cloudFiles = [];

    let data: FormData = new FormData();
    data.append("operatorID", this.$store.state.userID);
    data.append("pageSize", this.pageSize.toString());
    data.append("pageNum", currPage ? currPage.toString() : "1");

    axios.post(process.env.VUE_APP_cloud_file_list_by_uploader_url, data)
      .then(response => {
        if (response.data.hasError) {
          throw response.data.data;
        }

        if (currPage) {
          this.pageNum = currPage;
        }

        const payload = JSON.parse(response.data.data as string);

        this.total = payload.total;
        for (let i = 0; i < payload.files.length; i++) {
          const item = payload.files[i];

          this.cloudFiles.push({
            fileID: item.fileID,
            fileName: item.fileName,
            fileURL: generateCloudFileURL(item.fileURL),
            isPublic: item.isPublic,
            isPublicDisplay: displayCloudFileIsPublic(item.isPublic),
            updateTime: item.updateTime,
            updateTimeDisplay: displayCloudFileTime(item.updateTime),
            createdTime: item.createdTime,
            createdTimeDisplay: displayCloudFileTime(item.createdTime)
          });
        }
      })
      .catch(err => {
        console.log("list cloud file by uploader failed, error:", err);
      })
  }
}
</script>

<style lang="scss">
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
