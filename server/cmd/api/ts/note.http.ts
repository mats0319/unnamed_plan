// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v0.2.5

import { axiosWrapper } from "./config"
import { AxiosResponse } from "axios"
import { Pagination } from "./common.go"
import { CreateNoteRes, CreateNoteReq, ListNoteRes, ListNoteReq, ModifyNoteRes, ModifyNoteReq, DeleteNoteRes, DeleteNoteReq } from "./note.go"

class NoteAxios {
    public createNote(is_anonymous: boolean, title: string, content: string): Promise<AxiosResponse<CreateNoteRes>> {
        let req: CreateNoteReq = {
            is_anonymous: is_anonymous,
            title: title,
            content: content,
        }

        return axiosWrapper.post("/note/create", req)
    }

    public listNote(only_operator: boolean, page: Pagination): Promise<AxiosResponse<ListNoteRes>> {
        let req: ListNoteReq = {
            only_operator: only_operator,
            page: page,
        }

        return axiosWrapper.post("/note/list", req)
    }

    public modifyNote(note_id: string, is_anonymous: boolean, title: string, content: string): Promise<AxiosResponse<ModifyNoteRes>> {
        let req: ModifyNoteReq = {
            note_id: note_id,
            is_anonymous: is_anonymous,
            title: title,
            content: content,
        }

        return axiosWrapper.post("/note/modify", req)
    }

    public deleteNote(note_id: string): Promise<AxiosResponse<DeleteNoteRes>> {
        let req: DeleteNoteReq = {
            note_id: note_id,
        }

        return axiosWrapper.post("/note/delete", req)
    }
}

export const noteAxios: NoteAxios = new NoteAxios()
