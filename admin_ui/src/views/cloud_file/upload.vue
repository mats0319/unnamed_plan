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

        <el-popover trigger="hover" placement="top" :content="fileNameTips">
          <i slot="reference" class="el-icon-warning-outline" />
        </el-popover>
      </el-form-item>

      <el-form-item label="是否公开">
        <el-checkbox v-model="isPublic">公开</el-checkbox>

        <el-popover trigger="hover" placement="top" :content="isPublicTips">
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
import axios from "axios";

@Component
export default class UploadCloudFile extends Vue {
  private fileName = "";
  private extensionName = "";
  private isPublic = false;

  private isLoadingFile = false;
  private fileList: FileList;

  // const
  private fileNameTips = "设置文件显示名称，默认使用已选择文件的文件名";
  private isPublicTips = "是否公开该文件，公开的文件可能被其他人查看";

  private mounted() {
    // placeholder
  }

  private upload(): void {
    let data: FormData = new FormData();
    data.append("operatorID", this.$store.state.userID);
    data.append("fileName", this.fileName);
    data.append("extensionName", this.extensionName);
    data.append("isPublic", this.isPublic.toString());
    data.append("file", this.fileList.item(0) as File);

    axios.post(process.env.VUE_APP_cloud_file_upload_url, data)
      .then(response => {
        if (response.data.hasError) {
          throw response.data.data;
        }

        const payload = JSON.parse(response.data.data as string);
        if (payload.isSuccess) {
          this.$message.success("上传文件成功");
        } else {
          this.$message.error("上传文件失败");
        }
      })
      .catch(err => {
        console.log("upload cloud file failed, error:", err);
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
  }

  private beforeUpload(): void {
    if (this.fileList && this.fileList.item(0).name.length < 1) {
      this.$message.info("请选择想要上传的文件");
      return;
    }

    this.upload();
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
