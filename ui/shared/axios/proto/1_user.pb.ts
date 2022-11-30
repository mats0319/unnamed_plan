// Generate File, Should not Edit.
// Author: mario.

import { Error, Pagination } from "./common.pb"

export namespace User {
  export class Data {
    user_id: string = "";
    user_name: string = "";
    nickname: string = "";
    is_locked: boolean = false;
    permission: number = 0;
  }

  export class LoginReq {
    user_name: string = "";
    password: string = "";
  }

  export class LoginRes {
    user_id: string = "";
    nickname: string = "";
    permission: number = 0;
    err: Error = new Error();
  }

  export class ListReq {
    operator_id: string = "";
    page: Pagination = new Pagination();
  }

  export class ListRes {
    total: number = 0;
    users: Array<Data> = new Array<Data>();
    err: Error = new Error();
  }

  export class CreateReq {
    operator_id: string = "";
    user_name: string = "";
    password: string = "";
    permission: number = 0;
  }

  export class CreateRes {
    err: Error = new Error();
  }

  export class LockReq {
    operator_id: string = "";
    user_id: string = "";
  }

  export class LockRes {
    err: Error = new Error();
  }

  export class UnlockReq {
    operator_id: string = "";
    user_id: string = "";
  }

  export class UnlockRes {
    err: Error = new Error();
  }

  export class ModifyInfoReq {
    operator_id: string = "";
    user_id: string = "";
    curr_pwd: string = "";
    nickname: string = "";
    password: string = "";
  }

  export class ModifyInfoRes {
    err: Error = new Error();
  }

  export class ModifyPermissionReq {
    operator_id: string = "";
    user_id: string = "";
    permission: number = 0;
  }

  export class ModifyPermissionRes {
    err: Error = new Error();
  }

  export class AuthenticateReq {
    user_id: string = "";
    password: string = "";
  }

  export class AuthenticateRes {
    err: Error = new Error();
  }

}
