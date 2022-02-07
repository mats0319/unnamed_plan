import { axiosWrapper } from "./config";
import { calcSHA256 } from "../utils";

class NoteAxios {
  public create(operatorID: string, topic: string, content: string, isPublic: boolean) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("topic", topic);
    data.append("content", content);
    data.append("isPublic", isPublic.toString());
    return axiosWrapper.post("/api/note/create", data);
  }

  public listByWriter(operatorID: string, pageSize: number, pageNum: number) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("pageSize", pageSize.toString());
    data.append("pageNum", pageNum.toString());
    return axiosWrapper.post("/api/note/listByWriter", data);
  }

  public listPublic(operatorID: string, pageSize: number, pageNum: number) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("pageSize", pageSize.toString());
    data.append("pageNum", pageNum.toString());
    return axiosWrapper.post("/api/note/listPublic", data);
  }

  public modify(operatorID: string, noteID: string, password: string, topic: string, content: string, isPublic: boolean) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("noteID", noteID);
    data.append("password", calcSHA256(password));
    data.append("topic", topic);
    data.append("content", content);
    data.append("isPublic", isPublic.toString());
    return axiosWrapper.post("/api/note/modify", data);
  }

  public delete(operatorID: string, noteID: string, password: string) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("noteID", noteID);
    data.append("password", calcSHA256(password));
    return axiosWrapper.post("/api/note/delete", data);
  }
}

const noteAxios = new NoteAxios();
export default noteAxios;
