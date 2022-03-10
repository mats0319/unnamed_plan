<template>
  <div class="cloud-file-upload">
    <el-form label-position="left" label-width="15%">
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

      <el-form-item>
        <el-button type="info" :disabled="isLoadingFile" plain @click="beforeUpload">
          {{ isLoadingFile ? "读取文件……" : "上传" }}
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { tips_IsPublic, tips_CloudFile_FileName } from "shared/ts/const";
import cloudFileAxios from "shared/ts/axios_wrapper/cloud_file";

@Component
export default class UploadCloudFile extends Vue {
  private fileName = "";
  private extensionName = "";
  private lastModifiedTime = 0;
  private isPublic = false;

  private isLoadingFile = false;
  private fileList: FileList;

  // const
  private tips_IsPublic = tips_IsPublic;
  private tips_CloudFile_FileName = tips_CloudFile_FileName;

  private mounted() {
    // placeholder
  }

  private uploadFile(): void {
    cloudFileAxios.upload(this.$store.state.userID, this.fileName, this.extensionName, this.lastModifiedTime,
      this.isPublic, this.fileList.item(0) as File)
      .then(response => {
        if (response.data["hasError"]) {
          throw response.data["data"];
        }

        this.$message.success("上传文件成功");

        this.fileName = "";
        this.isPublic = false;
      })
      .catch(err => {
        this.$message.error("上传文件失败");
        console.log("> upload file failed.", err);
      });
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

    this.lastModifiedTime = parseInt((this.fileList.item(0).lastModified / 1000).toString());
  }

  private beforeUpload(): void {
    if (!this.fileList || !this.fileList.item(0)) {
      this.$message.info("请选择想要上传的文件");
      return;
    }

    this.uploadFile();
  }
}
</script>

<style lang="scss">
.cloud-file-upload {
  padding: 7vh 15vw;
  text-align: left;

  .el-form-item {
    margin: 5vh 0;
  }

  .el-form-item__label {
    font-size: 2rem;
  }

  .el-input, .el-checkbox {
    width: 60%;
  }

  .el-popover__reference-wrapper {
    margin-left: 5vh;
  }
}
</style>
