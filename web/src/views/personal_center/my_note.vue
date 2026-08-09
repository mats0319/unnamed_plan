<template>
  <div class="pnote-tools">
    <elevated-button class="pt-item" :disabled="loading" :onClick="beforeCreate">写小纸条</elevated-button>
  </div>

  <el-table
    v-loading="loading"
    :data="noteStore.noteType == NoteType.MyNotes ? noteStore.notes : new Array<Note>()"
    height="60%"
  >
    <el-table-column type="expand">
      <template #default="scope">
        <el-descriptions class="pnote-table-description" title="小纸条详情" border column="1" label-width="15%">
          <el-descriptions-item label="小纸条ID">{{ scope.row.note_id }}</el-descriptions-item>
          <el-descriptions-item label="作者昵称">{{ scope.row.writer }}</el-descriptions-item>
          <el-descriptions-item label="主题">{{ scope.row.title }}</el-descriptions-item>
          <el-descriptions-item label="内容">{{ scope.row.content }}</el-descriptions-item>
          <el-descriptions-item label="写作时间">{{ displayTimestamp(scope.row.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="修改时间">{{ displayTimestamp(scope.row.updated_at) }}</el-descriptions-item>
        </el-descriptions>
      </template>
    </el-table-column>

    <el-table-column label="作者昵称" prop="writer" :min-width="2" />

    <el-table-column label="主题" prop="title" :min-width="5" show-overflow-tooltip />

    <el-table-column label="操作" :min-width="3">
      <template #default="scope">
        <div class="buttons-box">
          <outlined-button class="button-item" :onClick="()=>{beforeModify(scope.row)}">编辑</outlined-button>
          <outlined-button class="button-item" :onClick="()=>{beforeDelete(scope.row)}">删除</outlined-button>
        </div>
      </template>
    </el-table-column>
  </el-table>

  <el-pagination
    layout="prev,pager,next,->,total"
    :total="noteStore.count"
    :page-size="pageSize"
    :disabled="loading"
    background
    @current-change="listMyNote"
  />

  <el-dialog v-model="showCreateDialog" title="写小纸条">
    <el-form v-model="createNoteReq" label-width="20%">
      <el-form-item label="主题">
        <el-input v-model="createNoteReq.title" />
      </el-form-item>

      <el-form-item label="内容">
        <el-input v-model="createNoteReq.content" type="textarea" :rows="4" />
      </el-form-item>

      <el-form-item label="是否匿名发布">
        <el-switch v-model="createNoteReq.is_anonymous" />
        &emsp;{{ createNoteReq.is_anonymous ? "匿名" : "不匿名" }}
      </el-form-item>
    </el-form>

    <template #footer>
      <elevated-button bg="white" :onClick="()=>{showCreateDialog=false}">取消</elevated-button>
      <elevated-button bg="lightgray" :disabled="!canCreateFlag" :onClick="create">提交</elevated-button>
    </template>
  </el-dialog>

  <el-dialog v-model="showModifyDialog" title="编辑小纸条">
    <el-form v-model="modifyNoteReq" label-width="20%">
      <el-form-item label="小纸条ID">{{ originData.note_id }}</el-form-item>

      <el-form-item label="主题">
        <el-input v-model="modifyNoteReq.title" />
      </el-form-item>

      <el-form-item label="内容">
        <el-input v-model="modifyNoteReq.content" type="textarea" :rows="4" />
      </el-form-item>

      <el-form-item label="是否匿名发布">
        <el-switch v-model="modifyNoteReq.is_anonymous" />
        &emsp;{{ modifyNoteReq.is_anonymous ? "匿名" : "不匿名" }}
      </el-form-item>
    </el-form>

    <template #footer>
      <elevated-button bg="white" :onClick="()=>{showModifyDialog=false}">取消</elevated-button>
      <elevated-button bg="lightgray" :disabled="!canModifyFlag" :onClick="modify">提交</elevated-button>
    </template>
  </el-dialog>

  <el-dialog v-model="showDeleteDialog" title="删除小纸条">
    <el-form label-width="20%">
      <el-form-item><b>是否确认删除以下小纸条?</b></el-form-item>
      <el-form-item label="小纸条ID">{{ originData.note_id }}</el-form-item>
      <el-form-item label="主题">{{ originData.title }}</el-form-item>
      <el-form-item label="内容">{{ originData.content }}</el-form-item>
    </el-form>

    <template #footer>
      <elevated-button bg="white" :onClick="()=>{showDeleteDialog=false}">取消</elevated-button>
      <elevated-button bg="lightgray" :onClick="del">删除</elevated-button>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { onMounted, ref, watch } from "vue"
import { CreateNoteReq, DeleteNoteReq, ModifyNoteReq, Note } from "@/axios/ts/note.go.ts"
import { deepCopy, displayTimestamp } from "@/ts/util.ts"
import { NoteType, useNoteStore } from "@/pinia/note.ts"
import OutlinedButton from "@/components/outlined_button.vue"
import ElevatedButton from "@/components/elevated_button.vue"
import { pageSize } from "@/ts/data.ts"

const noteStore = useNoteStore()
const loading = ref<boolean>(false)

const showCreateDialog = ref<boolean>(false)
const canCreateFlag = ref<boolean>(false)
const createNoteReq = ref<CreateNoteReq>(new CreateNoteReq())

const showModifyDialog = ref<boolean>(false)
const canModifyFlag = ref<boolean>(false)
const originData = ref<Note>(new Note()) // use for modify and delete
const modifyNoteReq = ref<ModifyNoteReq>(new ModifyNoteReq())

const showDeleteDialog = ref<boolean>(false)
const deleteNoteReq = ref<DeleteNoteReq>(new DeleteNoteReq())

onMounted(() => {
    listMyNote()
})

function beforeCreate(): void {
    createNoteReq.value = new CreateNoteReq()

    showCreateDialog.value = true
}

async function create(): Promise<void> {
    loading.value = true
    await noteStore.create(createNoteReq.value.is_anonymous, createNoteReq.value.title, createNoteReq.value.content)

    showCreateDialog.value = false
    loading.value = false
}

function beforeModify(note: Note): void {
    originData.value = deepCopy(note)
    modifyNoteReq.value = {
        note_id: note.note_id,
        is_anonymous: note.is_anonymous,
        title: note.title,
        content: note.content,
    }

    showModifyDialog.value = true
}

async function modify(): Promise<void> {
    loading.value = true
    await noteStore.modify(
        modifyNoteReq.value.note_id,
        modifyNoteReq.value.is_anonymous,
        modifyNoteReq.value.title,
        modifyNoteReq.value.content,
    )

    showModifyDialog.value = false
    loading.value = false
}

function beforeDelete(note: Note): void {
    originData.value = deepCopy(note)
    deleteNoteReq.value.note_id = note.note_id

    showDeleteDialog.value = true
}

async function del(): Promise<void> {
    loading.value = true
    await noteStore.del(deleteNoteReq.value.note_id)

    showDeleteDialog.value = false
    loading.value = false
}

function listMyNote(pageNum: number = 1): void {
    noteStore.list(true, pageNum)
}

watch(createNoteReq, (newValue) => {
    canCreateFlag.value = newValue.content.length > 0
}, { deep: true })

watch(modifyNoteReq, (newValue) => {
    canModifyFlag.value =
        newValue.is_anonymous != originData.value.is_anonymous ||
        newValue.title != originData.value.title ||
        newValue.content != originData.value.content
}, { deep: true })
</script>

<style lang="less" scoped>
.pnote-tools {
	height: 4rem;

	.pt-item {
		width: 8rem;
	}
}

.el-pagination {
	height: calc(40% - 4rem);
}

.pnote-table-description {
	padding: 1rem 4rem;
	white-space: pre-wrap;
}

.buttons-box {
	display: flex;

	.button-item {
		margin-right: 0.5rem;
    height: 2rem;
	}
}
</style>
