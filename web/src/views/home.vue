<template>
  <div class="home center-hv color-bg-0">
    <div class="h-content center-vertical color-bg-1">
      <div class="hc-item hc-item-1 center-vertical">
        你不是你以为的你自己，也不是别人眼中的你自己，而是你以为的别人眼中的你自己。<br>
        <span class="hci-en">
          You are not who you think you are, nor who others see you as, but rather who you believe others see you as.
        </span>
      </div>

      <div class="hc-item hc-item-2 center-vertical">
        不服高人有罪<br>
        <span class="hci-en">It is a sin not to admire talented people.</span>
      </div>

      <div class="hc-item hc-item-3 center-vertical">
        让每一行代码都有意义<br>
        <span class="hci-en">Make every line of code indispensable.</span>
      </div>
    </div>
  </div>

  <bottom />
</template>

<script lang="ts" setup>
import Bottom from "@/views/components/bottom.vue"
import { onMounted } from "vue"
import { UserStatus, useUserStore } from "@/pinia/user.ts"

const userStore = useUserStore()

onMounted(() => {
    if (userStore.userStatus == UserStatus.exit) {
        userStore.exitLogin() // 配合http拦截器，如果是401错误，则退出登录
    }
})
</script>

<style scoped lang="less">
.home {
	height: calc(100vh - 6.25rem - 6.25rem);

	.h-content {
		width: fit-content;
		height: 80%;
		border-radius: 1rem;

		.hc-item {
			flex-direction: column;
			height: 28%;
			font-style: italic;

			@media (min-width: 1280px) {
				font-size: 1.6rem;
			}
			font-size: 1.4rem;

			.hci-en {
				@media (min-width: 1280px) {
					font-size: 1.2rem;
				}
				font-size: 1rem;
			}
		}

		.hc-item-1 {
			@media (min-width: 1280px) {
				width: 76rem;
				padding-left: 4rem;
			}
			width: 32rem;
			padding-left: 4rem;
		}

		.hc-item-2 {
			@media (min-width: 1280px) {
				width: 28rem;
				padding-left: 52rem;
			}
			width: 36rem;
			padding-left: 4rem;
		}

		.hc-item-3 {
			@media (min-width: 1280px) {
				width: 64rem;
				padding-left: 16rem;
			}
			width: 36rem;
			padding-left: 4rem;
		}
	}
}
</style>
