// Generate File, should not edit.
// Author: mario.

import { Error, Pagination } from "./common.pb"

export namespace User {
  export class Data {
    user_id: "";
    user_name: "";
    nickname: "";
    is_locked: false;
    permission: 0;
  }

  export class LoginReq {
    user_name: "";
    password: "";
  }

  export class LoginRes {
    user_id: "";
    nickname: "";
    permission: 0;
    err: Error;
  }

  export class ListReq {
    operator_id: "";
    page: Pagination;
  }

  export class ListRes {
    total: 0;
    users: Array<Data>;
    err: Error;
  }

  export class CreateReq {
    operator_id: "";
    user_name: "";
    password: "";
    permission: 0;
  }

  export class CreateRes {
    err: Error;
  }

  export class LockReq {
    operator_id: "";
    user_id: "";
  }

  export class LockRes {
    err: Error;
  }

  export class UnlockReq {
    operator_id: "";
    user_id: "";
  }

  export class UnlockRes {
    err: Error;
  }

  export class ModifyInfoReq {
    operator_id: "";
    user_id: "";
    curr_pwd: "";
    nickname: "";
    password: "";
  }

  export class ModifyInfoRes {
    err: Error;
  }

  export class ModifyPermissionReq {
    operator_id: "";
    user_id: "";
    permission: 0;
  }

  export class ModifyPermissionRes {
    err: Error;
  }

  export class AuthenticateReq {
    user_id: "";
    password: "";
  }

  export class AuthenticateRes {
    err: Error;
  }

}
