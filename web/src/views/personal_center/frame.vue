<template>
  <div class="personal-center color-bg-0">
    <div class="p-left">
      <elevated-button v-if="userStore.user.is_admin" class="pl-item" :onClick="routerToListUser">
        查看用户
      </elevated-button>

      <elevated-button class="pl-item" :onClick="routerToModifyUser">修改个人信息</elevated-button>

      <elevated-button class="pl-item" :onClick="routerToSetMFAStatus">设置MFA状态</elevated-button>

      <elevated-button class="pl-item" :onClick="routerToMyNote">我的小纸条</elevated-button>
    </div>

    <el-divider direction="vertical" />

    <div class="p-right">
      <router-view />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { routerLink } from "@/ts/util.ts"
import { useUserStore } from "@/pinia/user.ts"
import ElevatedButton from "@/components/elevated_button.vue"

const userStore = useUserStore()

function routerToListUser(): void { routerLink('pListUser') }

function routerToModifyUser(): void { routerLink('pModifyUser') }

function routerToSetMFAStatus(): void { routerLink('pSetMFAStatus') }

function routerToMyNote(): void { routerLink('pNote') }
</script>

<style lang="less" scoped>
.personal-center {
	display: flex;
	height: calc(100vh - 6.25rem);

	.p-left {
		width: calc(20vw - 1px);
		margin: 6rem 2vw;

		.pl-item {
			margin: 2rem auto;
      width: 80%;
      height: 4rem;
		}
    .pl-item:deep(button) {
      font-size: 1.4rem;
    }
	}

	.p-right {
		width: 60vw;
		margin: 8rem 8vw;
	}

	.el-divider--vertical {
		width: 0;
		height: 100%;
		margin: 0;
		border-left: 1px solid darkgray;
	}
}
</style>
