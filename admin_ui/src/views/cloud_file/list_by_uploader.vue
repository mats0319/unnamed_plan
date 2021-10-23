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
      <el-table-column label="修改时间" prop="updateTimeDisplay" min-width="3" show-overflow-tooltip />
      <el-table-column label="上传时间" prop="createdTimeDisplay" min-width="3" show-overflow-tooltip />
      <el-table-column label="操作" min-width="2">
        <template slot-scope="scope">
          <el-button
            type="info"
            size="mini"
            plain
            @click="beforeModifyFile(scope.row.fileID, scope.row.fileName, scope.row.isPublic)"
          >
            修改
          </el-button>

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
      @current-change="listByUploader"
    />

    <el-dialog
      class="cflbu-dialog"
      :visible.sync="modifyDialogController"
      title="修改文件"
      :before-close="resetModifyDialogData"
    >
      <div class="cflbud-content">
        <el-form label-position="left" label-width="20%">

          <el-form-item label="文件ID">{{ fileID }}</el-form-item>
          <el-form-item label="原文件名">{{ oldFileName }}</el-form-item>

          <hr />

          <el-form-item label="选择文件">
            <input
              type="file"
              id="cfu-file"
              accept="application/pdf"
              @change="onFileChanged"
            />
          </el-form-item>

          <el-form-item label="文件名">
            <el-input v-model="fileName" placeholder="文件名" />

            <el-popover trigger="hover" placement="top" :content="tips_CloudFile_FileName">
              <i slot="reference" class="el-icon-warning-outline" />
            </el-popover>
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
        <el-button type="info" plain @click="modifyFile">修改</el-button>
      </div>
    </el-dialog>

    <el-dialog
      class="cflbu-dialog"
      :visible.sync="deleteDialogController"
      title="删除文件"
      :before-close="resetDeleteDialogData"
    >
      <div class="cflbud-content">
        <el-form label-position="left" label-width="20%">
          <el-form-item label="文件ID">{{ fileID }}</el-form-item>
          <el-form-item label="文件名">{{ fileName }}</el-form-item>

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
        <el-button type="info" plain @click="deleteFile">删除</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import {
  CloudFile,
  displayIsPublic,
  displayTime,
  generateCloudFileURL
} from "@/ts/data";
import axios from "axios";
import { calcSHA256 } from "@/ts/utils";
import { tips_CloudFile_FileName, tips_IsPublic } from "@/ts/const";

@Component
export default class listByUploader extends Vue {
  private cloudFiles: Array<CloudFile> = new Array<CloudFile>();

  private total = 0;
  private pageSize = 10;
  private pageNum = 1;
  
  private modifyDialogController = false;
  private oldFileName = "";
  private extensionName = "";
  private lastModifiedTime = 0;
  private oldIsPublic = false;
  private isPublic = false;
  private fileList: FileList;

  private deleteDialogController = false;
  private password = "";
  private fileID = "";
  private fileName = "";

  // const
  private tips_IsPublic = tips_IsPublic;
  private tips_CloudFile_FileName = tips_CloudFile_FileName;

  private mounted() {
    this.listByUploader();
  }

  private listByUploader(currPage?: number): void {
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
            isPublicDisplay: displayIsPublic(item.isPublic),
            updateTime: item.updateTime,
            updateTimeDisplay: displayTime(item.updateTime),
            createdTime: item.createdTime,
            createdTimeDisplay: displayTime(item.createdTime)
          });
        }
      })
      .catch(err => {
        this.$message.error("获取当前用户上传的文件列表失败，错误：" + err);
      })
  }

  private modifyFile(): void {
    if (!this.isValidModifyParams()) {
      this.$message.info("当前未执行任何修改");
      return;
    }

    const pwd = calcSHA256(this.password);
    this.password = "";

    let data: FormData = new FormData();
    data.append("operatorID", this.$store.state.userID);
    data.append("fileID", this.fileID);
    data.append("password", pwd);
    data.append("fileName", this.fileName);
    data.append("isPublic", this.isPublic.toString());
    data.append("lastModifiedTime", this.lastModifiedTime.toString());
    if (this.fileList && this.fileList.item(0)) {
      data.append("file", this.fileList.item(0) as File);
      data.append("extensionName", this.extensionName);
    }

    axios.post(process.env.VUE_APP_cloud_file_modify_url, data)
      .then(response => {
        if (response.data.hasError) {
          throw response.data.data;
        }

        const payload = JSON.parse(response.data.data as string);
        if (payload.isSuccess) {
          this.$message.success("修改文件成功");

          this.listByUploader(this.pageNum);
        } else {
          this.$message.error("修改文件失败");
        }
      })
      .catch(err => {
        this.$message.error("修改文件失败，错误：" + err);
      })
      .finally(() => {
        this.modifyDialogController = false;
      })
  }

  private deleteFile(): void {
    const pwd = calcSHA256(this.password);
    this.password = "";

    let data: FormData = new FormData();
    data.append("operatorID", this.$store.state.userID);
    data.append("fileID", this.fileID);
    data.append("password", pwd);

    axios.post(process.env.VUE_APP_cloud_file_delete_url, data)
      .then(response => {
        if (response.data.hasError) {
          throw response.data.data;
        }

        const payload = JSON.parse(response.data.data as string);
        if (payload.isSuccess) {
          this.$message.success("删除文件成功");

          this.listByUploader(this.pageNum);
        } else {
          this.$message.error("删除文件失败");
        }
      })
      .catch(err => {
        this.$message.error("删除文件失败，错误：" + err);
      })
      .finally(() => {
        this.deleteDialogController = false;
      })
  }

  private onFileChanged(ev: Event): void {
    //@ts-ignore-next-line
    if (!ev.target || !ev.target.files || ev.target.files.length < 1) {
      return;
    }

    //@ts-ignore-next-line
    this.fileList = ev.target.files;

    const fileNameSplit = this.fileList.item(0)?.name.split(".");
    this.fileName = "";
    for (let i = 0; i < fileNameSplit.length - 1; i++) {
      this.fileName += fileNameSplit[i];
    }
    this.extensionName = fileNameSplit.pop();

    this.lastModifiedTime = this.fileList.item(0).lastModified;
  }

  private isValidModifyParams(): boolean {
    return (this.fileList && this.fileList.length > 0) || this.oldFileName != this.fileName || this.isPublic != this.oldIsPublic;
  }

  private beforeModifyFile(fileID: string, oldFileName: string, isPublic: boolean): void {
    this.fileID = fileID;
    this.oldFileName = oldFileName;
    this.oldIsPublic = isPublic;
    this.isPublic = isPublic;

    this.modifyDialogController = true;
  }

  private beforeDeleteFile(fileID: string, fileName: string): void {
    this.fileID = fileID;
    this.fileName = fileName;

    this.deleteDialogController = true;
  }

  private resetModifyDialogData(): void {
    this.password = "";
    this.fileID = "";
    this.fileName = "";
    this.extensionName = "";
    this.isPublic = false;

    this.modifyDialogController = false;
  }

  private resetDeleteDialogData(): void {
    this.password = "";
    this.fileID = "";
    this.fileName = "";

    this.deleteDialogController = false;
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

  .cflbu-dialog {
    text-align: left;

    .cflbud-content {
      padding: 2vh 15%;
    }

    .el-input, .el-checkbox {
      width: 60%;
    }

    .el-popover__reference-wrapper {
      margin-left: 5vh;
    }
  }
}
</style>
