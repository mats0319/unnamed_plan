// Generate File, should not edit.
// Author: mario.

import { Pagination, Error } from "./common.pb"

export namespace CloudFile {
  export class Data {
    file_id: "";
    file_name: "";
    last_modified_time: 0;
    file_url: "";
    is_public: false;
    update_time: 0;
    created_time: 0;
  }

  export enum ListRule  {
    UNSPECIFIED,
    UPLOADER,
    PUBLIC,
  }

  export class ListReq {
    rule: ListRule;
    operator_id: "";
    page: Pagination;
  }

  export class ListRes {
    total: 0;
    files: Array<Data>;
    err: Error;
  }

  export class UploadReq {
    operator_id: "";
    file: Blob;
    file_name: "";
    extension_name: "";
    file_size: 0;
    last_modified_time: 0;
    is_public: false;
  }

  export class UploadRes {
    err: Error;
  }

  export class ModifyReq {
    operator_id: "";
    file_id: "";
    password: "";
    file_name: "";
    extension_name: "";
    is_public: false;
    file: Blob;
    file_size: 0;
    last_modified_time: 0;
  }

  export class ModifyRes {
    err: Error;
  }

  export class DeleteReq {
    operator_id: "";
    password: "";
    file_id: "";
  }

  export class DeleteRes {
    err: Error;
  }

}
