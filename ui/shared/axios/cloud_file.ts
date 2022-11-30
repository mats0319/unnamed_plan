import { axiosWrapper } from "./config/config";
import { calcSHA256 } from "./utils/utils";
import { CloudFile } from "./proto/2_cloud_file.pb";

class CloudFileAxios {
  public upload(
    operatorID: string,
    file: File,
    fileName: string,
    extensionName: string,
    lastModifiedTime: number,
    isPublic: boolean,
  ) {
    let request: CloudFile.UploadReq = {
      operator_id: operatorID,
      file: file,
      file_name: fileName,
      extension_name: extensionName,
      last_modified_time: lastModifiedTime,
      is_public: isPublic,
    }
    return axiosWrapper.post("/api/cloudFile/upload", request);
  }

  public listByUploader(operatorID: string, pageSize: number, pageNum: number) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("pageSize", pageSize.toString());
    data.append("pageNum", pageNum.toString());
    return axiosWrapper.post("/api/cloudFile/listByUploader", data);
  }

  public listPublic(operatorID: string, pageSize: number, pageNum: number) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("pageSize", pageSize.toString());
    data.append("pageNum", pageNum.toString());
    return axiosWrapper.post("/api/cloudFile/listPublic", data);
  }

  public modify(
    operatorID: string,
    fileID: string,
    password: string,
    fileName: string,
    isPublic: boolean,
    lastModifiedTime: number,
    file?: File,
    extensionName?: string
  ) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("fileID", fileID);
    data.append("password", calcSHA256(password));
    data.append("fileName", fileName);
    data.append("isPublic", isPublic.toString());
    data.append("lastModifiedTime", lastModifiedTime.toString());
    if (file && extensionName) {
      data.append("file", file);
      data.append("extensionName", extensionName);
    }
    return axiosWrapper.post("/api/cloudFile/modify", data);
  }

  public delete(operatorID: string, fileID: string, password: string) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("fileID", fileID);
    data.append("password", calcSHA256(password));
    return axiosWrapper.post("/api/cloudFile/delete", data);
  }
}

const cloudFileAxios = new CloudFileAxios();
export default cloudFileAxios;
