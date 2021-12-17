<template>
  <div id="app">
    <router-view/>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class App extends Vue {
  private created() {
    if (sessionStorage.getItem("vuex")) {
      this.$store.replaceState(
        Object.assign(
          {},
          this.$store.state,
          JSON.parse(sessionStorage.getItem("vuex") as string)
        )
      );

      sessionStorage.removeItem("vuex");
    }
  }

  private mounted() {
    window.addEventListener("beforeunload", () => {
      sessionStorage.setItem("vuex", JSON.stringify(this.$store.state));
    });
  }
}
</script>

<style lang="scss">
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
