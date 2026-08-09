<template>
  <el-form v-model="setMFAStatusReq" class="set-mfa-status" labelWidth="20%" size="large">
    <el-form-item label="是否启用MFA">
      <el-switch v-model="setMFAStatusReq.enable_mfa"/>&emsp;
      {{ setMFAStatusReq.enable_mfa ? "是" : "否" }}
    </el-form-item>

    <template v-if="setMFAStatusReq.enable_mfa">
      <el-form-item label="TOTP密钥">
        <div class="item">
          <span v-if="setMFAStatusReq.apply_new_key_flag">
            <div>{{ totpKey }}</div>
            <img :src="qrCodeText" alt="loading..."/>
          </span>

          <span v-else-if="userStore.user.has_totp_key">
            <b><i>[&nbsp;继续使用历史密钥&nbsp;]&emsp;</i></b>
          </span>

          <elevatedButton :onClick="applyNewTOTPKey">申请新的密钥</elevatedButton>
        </div>
      </el-form-item>

      <el-form-item label="TOTP Code">
        <el-input-otp v-model="setMFAStatusReq.totp_code" type="filled"/>
      </el-form-item>
    </template>

    <el-form-item>
      <elevated-button :onClick="setMFAStatus">设置MFA状态</elevated-button>
    </el-form-item>
  </el-form>
</template>

<script lang="ts" setup>
import { useUserStore } from "@/pinia/user.ts"
import { SetMFAStatusReq } from "@/axios/ts/user.go.ts"
import { onMounted, ref } from "vue"
import ElevatedButton from "@/components/elevated_button.vue"
import { generateQRCode } from "@/ts/util.ts"

const userStore = useUserStore()

const qrCodeText = ref<string>("")

const totpKey = ref<string>("")
const setMFAStatusReq = ref<SetMFAStatusReq>(new SetMFAStatusReq())

onMounted(() => {
    setMFAStatusReq.value = new SetMFAStatusReq()
    setMFAStatusReq.value.enable_mfa = userStore.user.enable_mfa
})

async function applyNewTOTPKey() {
    const res = await userStore.applyTOTPKey()

    totpKey.value = res.totp_key
    setMFAStatusReq.value.apply_new_key_flag = true

    qrCodeText.value = await generateQRCode(totpKey.value)
}

async function setMFAStatus() {
    await userStore.setMFAStatus(
        setMFAStatusReq.value.enable_mfa,
        setMFAStatusReq.value.apply_new_key_flag,
        setMFAStatusReq.value.totp_code,
    )
}
</script>

<style lang="less" scoped>
.set-mfa-status {
  .item {
    font-size: 1.2rem;
    width: 100%;
  }

  .tips {
    opacity: 0.6;
  }
}
</style>
