<template>
  <el-form class="cloud-file-upload" label-position="left" label-width="25%">
    <el-form-item label="选择文件">
      <input type="file" id="cfu-file" accept="application/pdf" @change="onFileChanged" />
    </el-form-item>

    <el-form-item label="文件名">
      <el-input v-model="fileName" placeholder="文件名" />

      <el-tooltip :content="tips.cloudFileFileName" effect="light" placement="top">
        <el-icon size="2rem"><InfoFilled /></el-icon>
      </el-tooltip>
    </el-form-item>

    <el-form-item label="是否公开">
      <el-checkbox v-model="isPublic">公开</el-checkbox>

      <el-tooltip :content="tips.isPublic" effect="light" placement="top">
        <el-icon size="2rem"><InfoFilled /></el-icon>
      </el-tooltip>
    </el-form-item>

    <el-divider />

    <el-form-item label="密码">
      <el-input v-model="password" type="password" placeholder="请输入密码" show-password clearable />
    </el-form-item>

    <el-form-item>
      <el-button type="info" :disabled="isLoadingFile" plain @click="beforeUpload">
        {{ isLoadingFile ? "读取文件……" : "上传" }}
      </el-button>
    </el-form-item>
  </el-form>
</template>

<script lang="ts">
import { defineComponent } from "vue"
import { tips } from "@/ts/const"
import cloudFileAxios from "@/ts/axios/cloud_file"
import { timestampMSToS } from "@/ts/utils"

export default defineComponent({
  name: "CloudFileUpload",
  data() {
    return {
      fileName: "",
      extensionName: "",
      lastModifiedTime: 0,
      isPublic: false,
      password: "",

      isLoadingFile: false,
      fileList: FileList,

      // const
      tips: tips,
    }
  },
  mounted() {
    // placeholder
  },
  methods: {
    uploadFile(): void {
      cloudFileAxios.upload(
        this.$store.state.userID,
        this.fileList.item(0) as File,
        this.fileName,
        this.extensionName,
        this.lastModifiedTime,
        this.isPublic,
        this.password,
      )
        .then(response => {
          if (response.err) {
            throw response.err
          }

          this.$message.success("上传文件成功")

          this.fileName = ""
          this.isPublic = false
        })
        .catch(err => {
          this.$message.error("上传文件失败")
          console.log("> upload file failed, error: ", err)
        })
    },

    beforeUpload(): void {
      if (!this.fileList || !this.fileList.item(0)) {
        this.$message.info("请选择想要上传的文件");
        return;
      }

      this.uploadFile();
    },

    onFileChanged(ev: Event): void {
      //@ts-ignore-next-line
      if (!ev.target || !ev.target.files || ev.target.files.length < 1) {
        return
      }

      //@ts-ignore-next-line
      this.fileList = ev.target.files

      const fileNameSplit = this.fileList.item(0)?.name.split(".")
      this.fileName = ""
      for (let i = 0; i < fileNameSplit.length - 1; i++) {
        this.fileName += fileNameSplit[i]
      }
      this.extensionName = fileNameSplit.pop()

      this.lastModifiedTime = timestampMSToS(this.fileList.item(0).lastModified) 
    },
  }
})
</script>

<style lang="less">
.cloud-file-upload {
  margin: 7rem 20vw 0;

  .el-form-item {
    margin: 5vh 0;
  }

  .el-form-item__label {
    font-size: 2rem;
  }

  .el-input, .el-checkbox {
    width: 60%;
  }

  .el-icon {
    margin-left: 5%;
  }
}

.el-popper {
  font-size: 1.8rem;
}
</style>
