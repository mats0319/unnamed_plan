// Generate File, should not edit.
// Author: mario.

import { Error } from "./common.pb"

export namespace Task {
  export class Data {
    task_id: "";
    task_name: "";
    description: "";
    pre_task_ids: Array<"">;
    status: 0;
    update_time: 0;
    created_time: 0;
  }

  export class ListReq {
    operator_id: "";
  }

  export class ListRes {
    total: 0;
    tasks: Array<Data>;
    err: Error;
  }

  export class CreateReq {
    operator_id: "";
    task_name: "";
    description: "";
    pre_task_ids: Array<"">;
  }

  export class CreateRes {
    err: Error;
  }

  export class ModifyReq {
    operator_id: "";
    task_id: "";
    password: "";
    task_name: "";
    description: "";
    pre_task_ids: Array<"">;
    status: 0;
  }

  export class ModifyRes {
    err: Error;
  }

}
