// Generate File, Should not Edit.
// Author: mario.

import { Pagination, Error } from "./common.pb"

export namespace Note {
  export class Data {
    note_id: string = "";
    write_by: string = "";
    topic: string = "";
    content: string = "";
    is_public: boolean = false;
    update_time: number = 0;
    created_time: number = 0;
  }

  export enum ListRule {
    UNSPECIFIED,
    WRITER,
    PUBLIC,
  }

  export class ListReq {
    rule: ListRule = ListRule.UNSPECIFIED;
    operator_id: string = "";
    page: Pagination = new Pagination();
  }

  export class ListRes {
    total: number = 0;
    notes: Array<Data> = new Array<Data>();
    err: Error = new Error();
  }

  export class CreateReq {
    operator_id: string = "";
    topic: string = "";
    content: string = "";
    is_public: boolean = false;
  }

  export class CreateRes {
    err: Error = new Error();
  }

  export class ModifyReq {
    operator_id: string = "";
    note_id: string = "";
    password: string = "";
    topic: string = "";
    content: string = "";
    is_public: boolean = false;
  }

  export class ModifyRes {
    err: Error = new Error();
  }

  export class DeleteReq {
    operator_id: string = "";
    password: string = "";
    note_id: string = "";
  }

  export class DeleteRes {
    err: Error = new Error();
  }

}
