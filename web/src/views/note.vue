<template>
  <div class="note color-bg-0 center-hv">
    <stacked-cards class="n-card color-bg-1">
      <div class="nc-content">
        <p class="ncc-title">
          小纸条<span class="ncct-right">{{ noteStore.currentNum }} /
            {{ noteStore.noteType == NoteType.AllNotes ? noteStore.count : 0 }}</span>
        </p>

        <el-form v-model="noteStore.currentNote" class="ncc-form" label-width="20%">
          <el-form-item label="ID">{{ noteStore.currentNote.note_id }}</el-form-item>
          <el-form-item label="作者">
            {{ noteStore.currentNote.is_anonymous ? "-" : noteStore.currentNote.writer }}
          </el-form-item>
          <el-form-item label="主题">{{ noteStore.currentNote.title }}</el-form-item>
          <el-form-item label="内容">{{ noteStore.currentNote.content }}</el-form-item>
          <el-form-item label="写作时间">{{ displayTimestamp(noteStore.currentNote.created_at) }}</el-form-item>
          <el-form-item label="修改时间">{{ displayTimestamp(noteStore.currentNote.updated_at) }}</el-form-item>
        </el-form>
      </div>
    </stacked-cards>

    <elevated-button class="n-options" :onClick="noteStore.next">下一张</elevated-button>
  </div>

  <bottom />
</template>

<script lang="ts" setup>
import Bottom from "@/views/components/bottom.vue"
import StackedCards from "@/components/stacked_cards.vue"
import { onMounted } from "vue"
import { NoteType, useNoteStore } from "@/pinia/note.ts"
import { displayTimestamp } from "@/ts/util.ts"
import ElevatedButton from "@/components/elevated_button.vue"

const noteStore = useNoteStore()

onMounted(() => {
    noteStore.list(false)
})
</script>

<style lang="less" scoped>
.note {
	height: calc(100vh - 6.25rem - 6.25rem);

	.n-card {
		@media (min-width: 1280px) {
			width: 80rem;
		}
		width: 40rem;
		height: calc(100vh - 6.25rem - 6.25rem - 10rem);

		position: relative;
		z-index: 3;

		.nc-content {
			padding: 2rem 4rem;

			.ncc-title {
				font-style: italic;
				font-weight: 600;
				height: 2rem;

				.ncct-right {
					float: right;
				}
			}

			.ncc-form {
				max-height: calc(100vh - 6.25rem - 6.25rem - 10rem - 8rem - 4rem);
				overflow-y: auto;
				white-space: pre-wrap;

				&::-webkit-scrollbar {
					display: none;
				}
			}
		}
	}

	.n-options {
		@media (min-width: 1280px) {
			bottom: 8rem;
		}

		width: 10rem;
		position: absolute;
		right: 1rem;
		bottom: 3rem;
	}
}
</style>
