<template>
  <div class="cloud-file-list">
    <el-form class="cfl-tools" label-position="left" label-width="fit-content">
      <el-form-item label="查询条件">
        <el-select v-model="selectedRuleIndex" @change="listFiles()">
          <el-option v-for="(item, index) in listRule" :key="index" :label="item.label" :value="index" />
        </el-select>
      </el-form-item>
    </el-form>

    <el-table class="cfl-table" :data="files" stripe highlight-current-row>
      <el-table-column label="文件名" min-width="3">
        <template #default="scope">
          <a :href="generateCloudFileURL(scope.row.file_url)" target="_blank">{{ scope.row.file_name }}</a>
        </template>
      </el-table-column>

      <el-table-column label="是否公开" min-width="1">
        <template #default="scope">{{ $filters.displayIsPublic(scope.row.is_public) }}</template>
      </el-table-column>

      <el-table-column label="上传时间" min-width="3">
        <template #default="scope">{{ $filters.displayTime(scope.row.created_time) }}</template>
      </el-table-column>

      <el-table-column label="修改时间" min-width="3">
        <template #default="scope">{{ $filters.displayTime(scope.row.update_time) }}</template>
      </el-table-column>

      <el-table-column v-if="selectedRuleIndex === 0" label="操作" min-width="3">
        <template #default="scope">
          <el-button
            type="info"
            plain
            @click="beforeModifyFile(scope.$index)"
          >
            修改
          </el-button>

          <el-button type="info" plain @click="beforeDeleteFile(scope.$index)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      :current-page="pageNum"
      :page-size="pageSize"
      :total="total"
      :layout="layout"
      @update:current-page="onPageChange"
    />

    <!--  modify  -->
    <el-dialog v-model="showModifyDialog" class="cfl-dialog">
      <template #header>
        <span class="cfld-title">修改文件</span>
      </template>

      <el-form label-position="left" label-width="30%">
        <el-form-item label="文件ID">{{ this.files[this.fileIndex].file_id }}</el-form-item>

        <el-form-item label="文件名">{{ this.files[this.fileIndex].file_name }}</el-form-item>

        <el-divider />

        <el-form-item label="选择文件">
          <input type="file" id="cfu-file" accept="application/pdf" @change="onFileChanged" />
        </el-form-item>

        <el-form-item label="文件名">
          <el-input v-model="fileName" placeholder="文件名" />
        </el-form-item>

        <el-form-item label="是否公开">
          <el-checkbox v-model="isPublic">公开</el-checkbox>
        </el-form-item>

        <el-form-item label="密码">
          <el-input v-model="password" type="password" placeholder="请输入密码" show-password clearable />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button plain @click="showModifyDialog = false">取消</el-button>
        <el-button type="info" plain @click="modifyFile">修改</el-button>
      </template>
    </el-dialog>

    <!--  delete  -->
    <el-dialog v-model="showDeleteDialog" class="cfl-dialog">
      <template #header>
        <span class="cfld-title">删除文件</span>
      </template>

      <el-form label-position="left" label-width="30%">
        <el-form-item label="文件ID">{{ this.files[this.fileIndex].file_id }}</el-form-item>

        <el-form-item label="文件名">{{ this.files[this.fileIndex].file_name }}</el-form-item>

        <el-divider />

        <el-form-item label="密码">
          <el-input v-model="password" type="password" placeholder="请输入密码" show-password clearable />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button plain @click="showDeleteDialog = false">取消</el-button>
        <el-button type="info" plain @click="deleteFile">删除</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue"
import { CloudFile } from "@/ts/axios/proto/2_cloud_file.pb"
import cloudFileAxios from "@/ts/axios/cloud_file"
import { timestampMSToS } from "@/ts/utils";

export default defineComponent({
  name: "CloudFileList",
  data() {
    return {
      files: [] as CloudFile.Data[],
      pageSize: 10,
      pageNum: 1,
      total: 0,
      layout: "prev,pager,next,->,total",

      selectedRuleIndex: 0,
      listRule: [
        { value: CloudFile.ListRule.UPLOADER, label: "查看我上传的文件" },
        { value: CloudFile.ListRule.PUBLIC, label: "查看公开的文件" },
      ],

      showModifyDialog: false,
      fileIndex: 0,
      fileName: "",
      isPublic: false,
      password: "",

      fileList: FileList,
      extensionName: "",
      lastModifiedTime: 0,

      showDeleteDialog: false,
    }
  },
  mounted() {
    this.listFiles()
  },
  methods: {
    listFiles(): void {
      cloudFileAxios.list(
        this.$store.state.userID,
        this.listRule[this.selectedRuleIndex].value,
        this.pageSize,
        this.pageNum,
      )
        .then(response => {
          if (response.err) {
            throw response.err
          }

          this.files = response.files ? response.files : new Array<CloudFile.Data>()
          this.total = response.total ? response.total : 0
        })
        .catch(err => {
          this.files = new Array<CloudFile.Data>()
          this.total = 0

          this.$message.error("获取云文件列表失败")
          console.log("> list cloud file failed, error: ", err)
        })
    },

    onPageChange(newValue: number): void {
      this.pageNum = newValue
      this.listFiles()
    },

    modifyFile(): void {
      if (!this.isValidModifyParams()) {
        this.$message.info("当前未执行任何修改")
        return
      } else if (this.password.length < 1) {
        this.$message.info("请输入密码")
        return
      }

      let res: Promise<CloudFile.ModifyRes>
      if (this.fileList && this.fileList.item(0)) {
        res = cloudFileAxios.modify(this.$store.state.userID, this.fileID, this.password, this.fileName,
          this.isPublic, this.lastModifiedTime, this.fileList.item(0) as File, this.extensionName)
      } else {
        res = cloudFileAxios.modify(this.$store.state.userID, this.fileID, this.password, this.fileName,
          this.isPublic, this.lastModifiedTime)
      }

      res
        .then(response => {
          if (response.err) {
            throw response.err
          }

          this.$message.success("修改文件成功")
          this.showModifyDialog = false

          this.listFiles()
        })
        .catch(err => {
          this.$message.error("修改文件失败")
          console.log("> modify file failed.", err)
        })
        .finally(() => {
          this.password = ""
        })
    },

    deleteFile(): void {
      if (this.password.length < 1) {
        this.$message.info("请输入密码")
        return
      }

      cloudFileAxios.delete(this.$store.state.userID, this.fileID, this.password)
        .then(response => {
          if (response.err) {
            throw response.err
          }

          this.$message.success("删除文件成功")
          this.showDeleteDialog = false

          this.onPageChange(1)
        })
        .catch(err => {
          this.$message.error("删除文件失败")
          console.log("> delete cloud file failed.", err)
        })
        .finally(() => {
          this.password = ""
        })
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

    beforeModifyFile(index: number): void {
      this.fileIndex = index

      this.showModifyDialog = true
    },

    beforeDeleteFile(index: number): void {
      this.fileIndex = index

      this.showDeleteDialog = true
    },

    // isValidModifyParams require change one or more of 'file' / 'file name' / 'is public'
    isValidModifyParams(): boolean {
      return (this.fileList && this.fileList.length > 0) ||
        this.files[this.fileIndex].file_name != this.fileName ||
        this.files[this.fileIndex].is_public != this.isPublic
    },

    generateCloudFileURL(url: string): string {
      return import.meta.env.Vite_axios_base_url + "/cloud-file/" + url
    },
  }
})
</script>

<style lang="less">
.cloud-file-list {
  height: inherit;
  padding-right: 2rem;

  .cfl-tools {
    height: 5rem;
    padding: 1rem 2rem;

    .el-form-item__label {
      font-size: 2rem;
    }
  }

  .cfl-table {
    margin-bottom: 2rem;
    height: calc(100% - 13rem);

    a {
      color: darkgray;

      &:hover {
        color: lightgray;
      }
    }
  }

  .cfl-dialog {
    text-align: left;

    .cfld-title {
      font-size: 3rem;
      font-weight: 600;
    }

    .el-form-item, .el-divider {
      margin-left: 20%;
    }

    .el-divider, .el-input {
      width: 60%;
    }

    .el-form-item__label {
      font-size: 2rem;
    }
  }
}
</style>
