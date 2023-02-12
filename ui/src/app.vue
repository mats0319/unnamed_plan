<template>
  <div id="app">
    <top class="a-top" />

    <div class="a-content">
      <router-view />
    </div>

    <bottom class="a-bottom" />
  </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Top from "@/components/top.vue";
import Bottom from "@/components/bottom.vue";

export default defineComponent({
  name: "App",
  components: {
    Top,
    Bottom,
  },
  created() {
    if (sessionStorage.getItem(import.meta.env.Vite_axios_source_sign)) {
      this.$store.replaceState(
        Object.assign(
          {},
          this.$store.state,
          JSON.parse(sessionStorage.getItem(import.meta.env.Vite_axios_source_sign))
        )
      )

      sessionStorage.removeItem(import.meta.env.Vite_axios_source_sign)
    }
  },
  mounted() {
    window.addEventListener("beforeunload", () => {
      sessionStorage.setItem(import.meta.env.Vite_axios_source_sign, JSON.stringify(this.$store.state))
    })
  }
})
</script>

<style lang="less">
#app {
  width: 100vw;
  height: 100vh;

  .a-top {
    height: 10rem;
    background-color: lightgray; // for test
  }

  .a-content {
    height: calc(100vh - 20rem);
    min-height: 50rem;
  }

  .a-bottom {
    height: 10rem;
    background-color: lightgray;
  }
}
</style>
