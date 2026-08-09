<template>
  <div class="elevated-button">
    <lock-button
      :style="{ 'background-color': props.bg }"
      :disabled="props.disabled"
      :on-click="props.onClick"
    >
      <slot />
    </lock-button>
  </div>
</template>

<script lang="ts" setup>
import LockButton from "@/components/lock_button.vue"

const props = defineProps({
    bg: {
        type: String,
        default: "rgb(240, 239, 226)",
    },
    disabled: {
        type: Boolean,
        default: false,
    },
    onClick: {
        type: Function,
        required: true,
    },
})
</script>

<style lang="less" scoped>
.elevated-button{
  width: fit-content;
  height: fit-content;

  button {
    width: 100%;
    height: 100%;
    border: 1px solid transparent; // 为了不让hover/loading时出现的border造成布局抖动，这里不使用border:none
  }
  button:hover,
  button.loading {
    border-color: darkgray;
    opacity: 0.8;
  }
}
</style>
