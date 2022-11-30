// Generate File, Should not Edit.
// Author: mario.

import { Pagination, Error } from "./common.pb"

export namespace CloudFile {
  export class Data {
    file_id: string = "";
    file_name: string = "";
    last_modified_time: number = 0;
    file_url: string = "";
    is_public: boolean = false;
    update_time: number = 0;
    created_time: number = 0;
  }

  export enum ListRule {
    UNSPECIFIED,
    UPLOADER,
    PUBLIC,
  }

  export class ListReq {
    rule: ListRule = ListRule.UNSPECIFIED;
    operator_id: string = "";
    page: Pagination = new Pagination();
  }

  export class ListRes {
    total: number = 0;
    files: Array<Data> = new Array<Data>();
    err: Error = new Error();
  }

  export class UploadReq {
    operator_id: string = "";
    file: Blob = new Blob();
    file_name: string = "";
    extension_name: string = "";
    last_modified_time: number = 0;
    is_public: boolean = false;
  }

  export class UploadRes {
    err: Error = new Error();
  }

  export class ModifyReq {
    operator_id: string = "";
    file_id: string = "";
    password: string = "";
    file_name: string = "";
    extension_name: string = "";
    is_public: boolean = false;
    file: Blob = new Blob();
    file_size: number = 0;
    last_modified_time: number = 0;
  }

  export class ModifyRes {
    err: Error = new Error();
  }

  export class DeleteReq {
    operator_id: string = "";
    password: string = "";
    file_id: string = "";
  }

  export class DeleteRes {
    err: Error = new Error();
  }

}
