// Generate File, Should not Edit.
// Author: mario.

import { Error } from "./common.pb"

export namespace ConfigCenter {
  export class GetServiceConfigReq {
    service_id: string = "";
    level: string = "";
  }

  export class GetServiceConfigRes {
    config: string = "";
    err: Error = new Error();
  }

}

export namespace RegistrationCenterCore {
  export class RegisterReq {
    service_id: string = "";
    target: string = "";
  }

  export class RegisterRes {
    err: Error = new Error();
  }

  export class ListServiceTargetReq {
    service_id: string = "";
  }

  export class ListServiceTargetRes {
    targets: Array<string> = new Array<string>();
    err: Error = new Error();
  }

}

export namespace RegistrationCenterEmbedded {
  export class CheckHealthReq {
    data: string = "";
  }

  export class CheckHealthRes {
    err: Error = new Error();
  }

}
