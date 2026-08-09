<template>
  <div class="color-bg-1 center-hv">
    <div class="top center-hv">
      <div class="t-title" @click="routerLink('home')">Unnamed Plan Web</div>

      <div v-if="flags.wildScreenFlag" class="t-content">
        <div class="tc-item center-hv" @click="routerLink('gDefault')">小游戏</div>

        <div class="tc-item center-hv" @click="routerLink('note')">小纸条</div>

        <div class="tc-item center-hv">
          <a href="https://github.com/mats0319/unnamed_plan" target="_blank">本站代码</a>
        </div>

        <div class="tc-item center-hv">
          <outlined-button v-show="!userStore.isLogin()" :onClick="beforeOpenLoginDialog">
            登录
          </outlined-button>

          <div v-show="userStore.isLogin()">
            <el-dropdown placement="bottom-end">
              <template #default>
                <span> {{ userStore.user.nickname }}&nbsp;<span class="icon">&or;</span> </span>
              </template>

              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="routerLink('pDefault')">个人中心</el-dropdown-item>
                  <el-dropdown-item divided @click="exitLogin()">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </div>

      <!--	移动端只允许查看小纸条	-->
      <div v-else class="t-content" @click="routerLink('note')">小纸条</div>
    </div>
  </div>

  <el-dialog v-model="showLoginDialog" title="登录">
    <el-form v-model="loginReq" label-width="20%">
      <el-form-item label="用户名"><el-input v-model="loginReq.user_name" /></el-form-item>

      <el-form-item label="密码">
        <el-input v-model="loginReq.password" show-password />
      </el-form-item>

      <el-form-item>
        <outlined-button :disabled="!canLoginFlag" :onClick="register">注册新用户</outlined-button>
      </el-form-item>
    </el-form>

    <el-dialog v-model="showLoginMFADialog" title="MFA" append-to-body>
      <el-form v-model="loginMFAReq" labelWidth="20%">
        <el-form-item label="TOTP Code">
          <el-input-otp v-model="loginMFAReq.totp_code" type="filled" />
        </el-form-item>
      </el-form>

      <template #footer>
        <elevated-button bg="white" :onClick="()=>{showLoginMFADialog=false}">取消</elevated-button>
        <elevated-button bg="lightgray" :disabled="!canLoginMFAFlag" :onClick="loginMFA">确认</elevated-button>
      </template>
    </el-dialog>

    <template #footer>
      <elevated-button bg="white" :onClick="()=>{showLoginDialog=false}">取消</elevated-button>
      <elevated-button bg="lightgray" :disabled="!canLoginFlag" :onClick="login">登录</elevated-button>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { useUserStore } from "@/pinia/user.js"
import { ref, watch } from "vue"
import { useFlagStore } from "@/pinia/flag.js"
import OutlinedButton from "@/components/outlined_button.vue"
import { routerLink } from "@/ts/util.js"
import { LoginReq, LoginMFAReq } from "@/axios/ts/user.go.js"
import ElevatedButton from "@/components/elevated_button.vue"

const flags = useFlagStore()
const userStore = useUserStore()

const showLoginDialog = ref<boolean>(false)
const canLoginFlag = ref<boolean>(false) // also 'can register flag'
const loginReq = ref<LoginReq>(new LoginReq())

const showLoginMFADialog = ref<boolean>(false)
const canLoginMFAFlag = ref<boolean>(false)
const loginMFAReq = ref<LoginMFAReq>(new LoginMFAReq())

function beforeOpenLoginDialog(): void {
    loginReq.value = new LoginReq()

    showLoginDialog.value = true
}

async function login(): Promise<void> {
    const res = await userStore.login(loginReq.value.user_name, loginReq.value.password)

    res.enable_mfa ? beforeOpenLoginMFADialog(res.mfa_token) : showLoginDialog.value = false
}

function beforeOpenLoginMFADialog(mfaToken: string): void {
    loginMFAReq.value = new LoginMFAReq()
    loginMFAReq.value.mfa_token = mfaToken

    showLoginMFADialog.value = true
}

async function loginMFA(): Promise<void> {
    await userStore.loginMFA(loginMFAReq.value.mfa_token, loginMFAReq.value.totp_code)

    showLoginMFADialog.value = false
    showLoginDialog.value = false
}

async function register(): Promise<void> {
    await userStore.register(loginReq.value.user_name, loginReq.value.password)

    showLoginDialog.value = false
}

function exitLogin(): void {
    userStore.exitLogin()
    routerLink("home")
}

watch(loginReq, (newValue) => {
    canLoginFlag.value = newValue.user_name.length > 0 && newValue.password.length > 0
}, { deep: true })

watch(loginMFAReq, (newValue) => {
    canLoginMFAFlag.value = newValue.mfa_token.length > 0 && newValue.totp_code.length > 0
}, { deep: true })
</script>

<style lang="less" scoped>
.top {
	@media (min-width: 1280px) {
		width: 80rem;
	}
	width: 40rem;
	height: 6.25rem;

	.t-title {
		@media (min-width: 1280px) {
			width: 24rem;
			padding: 0 3rem;
			font-size: 2.2rem;
		}
		width: 50%;
		font-size: 1.6rem;
	}

	.t-title:hover {
		cursor: pointer;
	}

	.t-content {
		display: flex;
		justify-content: flex-end;
		@media (min-width: 1280px) {
			width: 50rem;
		}
		width: 50%;

		.tc-item {
			width: 8rem;
			padding: 0 1rem;
			font-size: 1.2rem;

			a {
				color: black;
				text-decoration-line: none;
			}

			.el-dropdown {
				font-size: 1.6rem;

				.icon {
					font-size: 70%;
				}
			}
		}

		.tc-item:hover {
			cursor: pointer;
		}
	}
}
</style>
