import { axiosWrapper } from "./config/config"
import { calcSHA256, objectToFormData } from "./utils"
import { CloudFile } from "./proto/2_cloud_file.pb"

class CloudFileAxios {
  public upload(
    operatorID: string,
    file: File,
    fileName: string,
    extensionName: string,
    lastModifiedTime: number,
    isPublic: boolean,
    password: string,
  ): Promise<CloudFile.UploadRes> {
    let request: CloudFile.UploadReq = {
      operator_id: operatorID,
      file: file,
      file_name: fileName,
      extension_name: extensionName,
      file_size: 0,
      last_modified_time: lastModifiedTime,
      is_public: isPublic,
      password: calcSHA256(password),
    }
    
    return axiosWrapper.post("/api/cloudFile/upload", objectToFormData(request))
  }

  public list(operatorID: string, rule: CloudFile.ListRule, pageSize: number, pageNum: number): Promise<CloudFile.ListRes> {
    let request: CloudFile.ListReq = {
      operator_id: operatorID,
      rule: rule,
      page: {
        page_size: pageSize,
        page_num: pageNum,
      },
    }

    return axiosWrapper.post("/api/cloudFile/list", objectToFormData(request))
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
  ): Promise<CloudFile.ModifyRes> {
    let request: CloudFile.ModifyReq = {
      operator_id: operatorID,
      file_id: fileID,
      password: calcSHA256(password),
      file_name: fileName,
      is_public: isPublic,
      last_modified_time: lastModifiedTime,
      file: new Blob(),
      file_size: 0,
      extension_name: "",
    }

    if (file && extensionName) {
      request.file = file
      request.extension_name = extensionName
    }

    return axiosWrapper.post("/api/cloudFile/modify", objectToFormData(request))
  }

  public delete(operatorID: string, fileID: string, password: string): Promise<CloudFile.DeleteRes> {
    let request: CloudFile.DeleteReq = {
      operator_id: operatorID,
      file_id: fileID,
      password: calcSHA256(password)
    }

    return axiosWrapper.post("/api/cloudFile/delete", objectToFormData(request))
  }
}

const cloudFileAxios = new CloudFileAxios()
export default cloudFileAxios
