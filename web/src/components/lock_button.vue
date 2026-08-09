<template>
  <button
    class="lock-button center-hv"
    :disabled="loading || props.disabled"
    @click="handleClick"
  >
    <span v-show="!loading"><slot /></span>
    <span v-show="loading">loading...</span>
  </button>
</template>

<script lang="ts" setup>
import { ref } from "vue"

const props = defineProps({
    disabled: {
        type: Boolean,
        required: false,
    },
    onClick: {
        type: Function,
        required: true,
    },
})

const loading = ref<boolean>(false)

async function handleClick(): Promise<void> {
    if (loading.value) {
        return
    }

    loading.value = true

    try {
        // 表面上对函数没有要求，常规函数、Promise、async函数都行，
        // 实际上对于Promise函数，应返回该Promise（形如：`return api.login()[.then()]`），
        // 如果没有return，程序会立刻执行完成并解除按钮锁定，而非预期中等待Promise执行完成才解除锁定。
        // 所以对Promise不熟悉的，可以考虑使用async/await
        await props.onClick()
    } finally {
        loading.value = false
    }
}
</script>

<style lang="less" scoped>
.lock-button {
	padding: 0.5rem 1rem;
  line-height: 1.5; // 防止中-英文切换导致的高度抖动
}
.lock-button:hover {
  cursor: pointer;
	text-decoration-line: underline;
}
.lock-button:disabled {
  pointer-events: none;
}
</style>
