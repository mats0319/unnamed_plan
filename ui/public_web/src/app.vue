<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class App extends Vue {
  private created() {
    if (sessionStorage.getItem(process.env.VUE_APP_axios_source_sign)) {
      this.$store.replaceState(
        Object.assign(
          {},
          this.$store.state,
          JSON.parse(sessionStorage.getItem(process.env.VUE_APP_axios_source_sign))
        )
      );

      sessionStorage.removeItem(process.env.VUE_APP_axios_source_sign);
    }
  }

  private mounted() {
    window.addEventListener("beforeunload", () => {
      sessionStorage.setItem(process.env.VUE_APP_axios_source_sign, JSON.stringify(this.$store.state));
    });
  }
}
</script>

<style lang="less">
html {
  font-size: 62.5%; // now, 1 rem = 10 px
}

body {
  margin: 0;
  padding: 0;
  overflow-x: hidden;
}

#app {
  font-family: KaiTi, STKaiti, Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
</style>
