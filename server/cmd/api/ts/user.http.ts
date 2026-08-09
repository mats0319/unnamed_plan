// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v0.2.5

import { axiosWrapper } from "./config"
import { AxiosResponse } from "axios"
import { RegisterRes, RegisterReq, LoginRes, LoginReq, LoginMFARes, LoginMFAReq, ListUserRes, ListUserReq, ModifyUserRes, ModifyUserReq, NewTOTPKeyRes, SetMFAStatusRes, SetMFAStatusReq } from "./user.go"
import { Pagination } from "./common.go"

class UserAxios {
    public register(user_name: string, password: string): Promise<AxiosResponse<RegisterRes>> {
        let req: RegisterReq = {
            user_name: user_name,
            password: password,
        }

        return axiosWrapper.post("/register", req)
    }

    public login(user_name: string, password: string): Promise<AxiosResponse<LoginRes>> {
        let req: LoginReq = {
            user_name: user_name,
            password: password,
        }

        return axiosWrapper.post("/login", req)
    }

    public loginMFA(mfa_token: string, totp_code: string): Promise<AxiosResponse<LoginMFARes>> {
        let req: LoginMFAReq = {
            mfa_token: mfa_token,
            totp_code: totp_code,
        }

        return axiosWrapper.post("/login-mfa", req)
    }

    public listUser(page: Pagination): Promise<AxiosResponse<ListUserRes>> {
        let req: ListUserReq = {
            page: page,
        }

        return axiosWrapper.post("/user/list", req)
    }

    public modifyUser(nickname: string, password: string): Promise<AxiosResponse<ModifyUserRes>> {
        let req: ModifyUserReq = {
            nickname: nickname,
            password: password,
        }

        return axiosWrapper.post("/user/modify", req)
    }

    public newTOTPKey(): Promise<AxiosResponse<NewTOTPKeyRes>> {
        return axiosWrapper.post("/totp-key/new")
    }

    public setMFAStatus(enable_mfa: boolean, apply_new_key_flag: boolean, totp_code: string): Promise<AxiosResponse<SetMFAStatusRes>> {
        let req: SetMFAStatusReq = {
            enable_mfa: enable_mfa,
            apply_new_key_flag: apply_new_key_flag,
            totp_code: totp_code,
        }

        return axiosWrapper.post("/mfa/set-status", req)
    }
}

export const userAxios: UserAxios = new UserAxios()
