<template>
  <div class="cloud-file-list-by-uploader">
    <el-table :data="cloudFiles" height="calc(80vh - 32px)" stripe highlight-current-row>
      <el-table-column label="文件名" min-width="2">
        <template slot-scope="scope">
          <div class="cflbu-file-name">
            <a :href="scope.row.fileURL" target="_blank">{{ scope.row.fileName }}</a>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="是否公开" prop="isPublicDisplay" min-width="1" show-overflow-tooltip />
      <el-table-column label="上传时间" prop="createdTimeDisplay" min-width="3" show-overflow-tooltip />
      <el-table-column label="操作" min-width="1">
        <template slot-scope="scope">
          <el-button
            type="info"
            size="mini"
            plain
            @click="beforeDeleteFile(scope.row.fileID, scope.row.fileName)"
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
      @current-change="listCloudFileByUploader"
    />

    <el-dialog
      class="cflbu-delete-dialog"
      :visible.sync="deleteDialogController"
      title="删除文件"
      :before-close="resetDialogData"
    >
      <div class="cflbudd-content">
        <el-form label-position="left" label-width="20%">
          <el-form-item label="文件ID">{{ fileID }}</el-form-item>
          <el-form-item label="文件名">{{ fileName }}</el-form-item>
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
        <el-button type="info" plain @click="deleteFile">删除</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import {
  CloudFile,
  displayCloudFileIsPublic,
  displayCloudFileTime,
  generateCloudFileURL
} from "@/ts/data";
import axios from "axios";
import { calcSHA256 } from "@/ts/sha256";

@Component
export default class ListCloudFileByUploader extends Vue {
  private cloudFiles: Array<CloudFile> = new Array<CloudFile>();

  private total = 0;
  private pageSize = 10;
  private pageNum = 1;

  private deleteDialogController = false;
  private password = "";
  private fileID = "";
  private fileName = "";

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

  private deleteFile(): void {
    const pwd = calcSHA256(this.password);
    this.password = "";

    let data: FormData = new FormData();
    data.append("operatorID", this.$store.state.userID);
    data.append("password", pwd);
    data.append("fileID", this.fileID);

    axios.post(process.env.VUE_APP_cloud_file_delete_url, data)
      .then(response => {
        if (response.data.hasError) {
          throw response.data.data;
        }

        const payload = JSON.parse(response.data.data as string);
        if (payload.isSuccess) {
          this.$message.success("删除文件成功");

          this.listCloudFileByUploader(this.pageNum);
        } else {
          this.$message.error("删除文件失败");
        }
      })
      .catch(err => {
        console.log("delete cloud file failed, error:", err);
      })
  }

  private beforeDeleteFile(fileID: string, fileName: string): void {
    this.fileID = fileID;
    this.fileName = fileName;

    this.deleteDialogController = true;
  }

  private resetDialogData(): void {
    this.password = "";
    this.fileID = "";
    this.fileName = "";
  }
}
</script>

<style lang="scss">
.cloud-file-list-by-uploader {
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

  .cflbu-delete-dialog {
    text-align: left;

    .cflbudd-content {
      padding: 2vh 15%;
    }
  }
}
</style>
