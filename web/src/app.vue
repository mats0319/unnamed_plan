<template>
  <lock-screen />
  <top v-if="!$route.meta.hideTop" />
  <router-view />
</template>

<script setup lang="ts">
import { useFlagStore } from "@/pinia/flag.ts"
import { onMounted } from "vue"
import { useUserStore } from "@/pinia/user.ts"
import { deepCopy } from "@/ts/util.ts"
import Top from "@/views/components/top.vue"
import LockScreen from "@/components/lock_screen.vue"

const flags = useFlagStore()
const userStore = useUserStore()

onMounted(() => {
    console.log("version: " + __PackageVersion__)

    // on re-size screen width
    flags.onScreenWidthChanged(screen.width)

    window.addEventListener("resize", () => {
        flags.onScreenWidthChanged(screen.width)
    })

    // keep 'login info' during refresh
    window.addEventListener("beforeunload", () => {
        localStorage.setItem("login_data", JSON.stringify(userStore.user))
    })

    const loginData = localStorage.getItem("login_data")
    if (loginData) {
        userStore.user = deepCopy(JSON.parse(loginData))
    }
})
</script>

<style lang="less">
body {
	padding: 0;
	margin: 0;

	overflow-x: hidden;
}
</style>
