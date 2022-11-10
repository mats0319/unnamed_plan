// Generate File, should not edit.
// Author: mario.

import { Error } from "./common.pb"

export namespace ConfigCenter {
  export class GetServiceConfigReq {
    service_id: "";
    level: "";
  }

  export class GetServiceConfigRes {
    config: "";
    err: Error;
  }

}

export namespace RegistrationCenterCore {
  export class RegisterReq {
    service_id: "";
    target: "";
  }

  export class RegisterRes {
    err: Error;
  }

  export class ListServiceTargetReq {
    service_id: "";
  }

  export class ListServiceTargetRes {
    targets: Array<"">;
    err: Error;
  }

}

export namespace RegistrationCenterEmbedded {
  export class CheckHealthReq {
    data: "";
  }

  export class CheckHealthRes {
    err: Error;
  }

}
