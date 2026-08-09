import { defineStore } from "pinia"
import { noteAxios } from "@/axios/ts/note.http.ts"
import { ListNoteRes, Note } from "@/axios/ts/note.go.ts"
import { log } from "@/ts/log.ts"
import { Pagination } from "@/axios/ts/common.go.ts"
import { ref } from "vue"

export enum NoteType {
    placeholder = -1,
    AllNotes = 1,
    MyNotes = 2,
}

const pageSize = 10

export const useNoteStore = defineStore("note", () => {
    const noteType = ref<NoteType>(NoteType.placeholder) // 标识count/notes的含义是全部/我的小纸条
    const count = ref<number>(0)
    const notes = ref<Array<Note>>(new Array<Note>())

    const currentNum = ref<number>(0) // 当前小纸条在全部小纸条中的序号，[0,count]
    const currentNote = ref<Note>(new Note())
    const nextPage = ref<number>(1)

    async function next(): Promise<void> { // 默认在全部小纸条上执行该函数
        let currentIndex = -1
        for (let i = 0; i < notes.value.length; i++) {
            if (currentNote.value.note_id == notes.value[i].note_id) {
                currentIndex = i
                break
            }
        }

        if (currentIndex >= notes.value.length - 1) { // 已经是当前数组的最后一个了
            if (currentNum.value >= count.value) { // 没有更多了，从头开始
                currentNum.value = 0
                await list(false)
            } else { // 下一页
                await list(false, nextPage.value)
            }
        } else { // 正常获取数组的下一个元素，同时可以兼容当前小纸条不存在的情况(`currentIndex=-1`)
            currentNum.value++
            currentNote.value = notes.value.length > 0 ? notes.value[currentIndex + 1] : new Note()
        }
    }

    async function list(onlyOperator: boolean, pageNum: number = 1): Promise<void> {
        const pagination: Pagination = { size: pageSize, num: pageNum }
        const { data }: { data: ListNoteRes } = await noteAxios.listNote(onlyOperator, pagination)

        noteType.value = onlyOperator ? NoteType.MyNotes : NoteType.AllNotes
        count.value = data.count
        notes.value = data.notes

        currentNum.value = pageNum * pageSize + 1
        currentNote.value = notes.value.length > 0 ? notes.value[0] : new Note()
        nextPage.value = pageNum + 1

        log.success("List Note")
    }

    async function create(isAnonymous: boolean, title: string, content: string): Promise<void> {
        await noteAxios.createNote(isAnonymous, title, content)

        log.success("Create Note")

        await list(true) // 默认写小纸条后查询自己的小纸条
    }

    async function modify(noteID: string, isAnonymous: boolean, title: string, content: string): Promise<void> {
        await noteAxios.modifyNote(noteID, isAnonymous, title, content)

        log.success("Modify Note")

        await list(true) // 默认修改小纸条后查询自己的小纸条
    }

    // delete是关键字，所以这里命名为del
    async function del(noteID: string): Promise<void> {
        await noteAxios.deleteNote(noteID)

        log.success("Delete Note")

        await list(true) // 默认删除小纸条后查询自己的小纸条
    }

    return { noteType, count, notes, currentNum, currentNote, next, list, create, modify, del }
})
