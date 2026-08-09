import { ref } from "vue"
import { defineStore } from "pinia"

export const useFlagStore = defineStore("flag", () => {
    const wildScreenFlag = ref<boolean>(true)
    const loading = ref<boolean>(false)

    function onScreenWidthChanged(width: number): void {
        wildScreenFlag.value = width > 1280
    }

    function setLoading(flag: boolean): void {
        loading.value = flag
    }

    return { wildScreenFlag, onScreenWidthChanged, loading, setLoading }
})
