// Generate File, should not edit.
// Author: mario.

import { Pagination, Error } from "./common.pb"

export namespace Note {
  export class Data {
    note_id: "";
    write_by: "";
    topic: "";
    content: "";
    is_public: false;
    update_time: 0;
    created_time: 0;
  }

  export enum ListRule  {
    UNSPECIFIED,
    WRITER,
    PUBLIC,
  }

  export class ListReq {
    rule: ListRule;
    operator_id: "";
    page: Pagination;
  }

  export class ListRes {
    total: 0;
    notes: Array<Data>;
    err: Error;
  }

  export class CreateReq {
    operator_id: "";
    topic: "";
    content: "";
    is_public: false;
  }

  export class CreateRes {
    err: Error;
  }

  export class ModifyReq {
    operator_id: "";
    note_id: "";
    password: "";
    topic: "";
    content: "";
    is_public: false;
  }

  export class ModifyRes {
    err: Error;
  }

  export class DeleteReq {
    operator_id: "";
    password: "";
    note_id: "";
  }

  export class DeleteRes {
    err: Error;
  }

}
