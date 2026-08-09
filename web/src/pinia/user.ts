import { defineStore } from "pinia"
import { ListUserRes, LoginRes, LoginMFARes, User, NewTOTPKeyRes } from "@/axios/ts/user.go.ts"
import { ref } from "vue"
import CryptoJs from "crypto-js"
import { userAxios } from "@/axios/ts/user.http.ts"
import { log } from "@/ts/log.ts"

export enum UserStatus {
    placeholder = -1,
    exit = 1,
}

export const useUserStore = defineStore("user", () => {
    const user = ref<User>(new User()) // current user
    const userStatus = ref<UserStatus>(UserStatus.placeholder)

    const count = ref<number>(0) // list users
    const users = ref<Array<User>>(new Array<User>())

    async function register(userName: string, password: string): Promise<void> {
        await userAxios.register(userName, calcSHA256(password))

        log.success("Register")

        await login(userName, password) // 新注册用户，此时的登录不涉及MFA
    }

    async function login(userName: string, password: string): Promise<LoginRes> {
        const { data }: { data: LoginRes } = await userAxios.login(userName, calcSHA256(password))

        // 开启MFA时，因为还需要用户输入totp code，所以不能在这里直接调`loginMFA`
        if (!data.enable_mfa) { // disable MFA, login done
            Object.assign(user.value, {
                user_name: data.user_name,
                nickname: data.nickname,
                is_admin: data.is_admin,
                enable_mfa: false,
                has_totp_key: data.has_totp_key,
            })

            log.success("Login")
        }

        return data
    }

    async function loginMFA(mfaToken: string, totpCode: string): Promise<void> {
        const { data }: { data: LoginMFARes } = await userAxios.loginMFA(mfaToken, totpCode)

        Object.assign(user.value, {
            user_name: data.user_name,
            nickname: data.nickname,
            is_admin: data.is_admin,
            enable_mfa: true,
            has_totp_key: true,
        })

        log.success("Login (Enable MFA)")
    }

    async function modify(nickname: string, password: string): Promise<void> {
        await userAxios.modifyUser(nickname, calcSHA256(password))

        if (nickname != "") {
            user.value.nickname = nickname
        }

        log.success("Modify User")
    }

    async function list(pageSize: number, pageNum: number): Promise<void> {
        const { data }: { data: ListUserRes } = await userAxios.listUser({ size: pageSize, num: pageNum })

        count.value = data.count
        users.value = data.users

        log.success("List User")
    }

    async function applyTOTPKey(): Promise<NewTOTPKeyRes> {
        const { data }:{ data: NewTOTPKeyRes } = await userAxios.newTOTPKey()

        log.success("Apply New TOTP Key")

        return data
    }

    async function setMFAStatus(enableMFA: boolean, applyNewKeyFlag: boolean, totpCode: string): Promise<void> {
        await userAxios.setMFAStatus(enableMFA, applyNewKeyFlag, totpCode)

        log.success("Set MFA Status")

        user.value.enable_mfa = true
        user.value.has_totp_key = true
    }

    function isLogin(): boolean { return user.value.user_name.length > 0 }

    function setUserStatus(status: UserStatus): void {
        userStatus.value = status
    }

    function exitLogin(): void {
        user.value = new User()
        localStorage.removeItem("access_token")
        localStorage.removeItem("login_data")
    }

    return { user, userStatus, count, users, isLogin, setUserStatus, exitLogin,
        register, login, loginMFA, modify, list, applyTOTPKey, setMFAStatus }
})

function calcSHA256(password: string): string { // internal func
    return password.length > 0 ? CryptoJs.SHA256(password).toString(CryptoJs.enc.Hex) : ""
}
