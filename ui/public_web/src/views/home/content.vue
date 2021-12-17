<template>
  <div class="default-home-content">
    <div class="sentence-one">
      <i>
        我一直在假装努力，你却在真的成长。<br />
        <label class="sentence-en">I've been pretending to try while you're really growing.</label>
      </i>
    </div>

    <div class="sentence-two">
      <i>
        不服高人有罪。<br />
        <label class="sentence-en">Guilty of disobeying a superior man.</label>
      </i>
    </div>

    <div class="sentence-three">
      <i>
        让每一行代码都有意义。<br />
        <label class="sentence-en">Making sense of every line of code.</label>
      </i>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class Content extends Vue {
  private mounted() {
    // 多登检查未通过（异地登录或超时），退出登录并返回到首页（此处仅执行退出登录逻辑）
    if (this.$route.params.exit && this.$route.params.exit.length > 0) {
      this.$store.state.isLogin = false;

      sessionStorage.removeItem("auth");
    }
  }
}
</script>

<style lang="scss">
.default-home-content {
  height: calc(100vh - 20rem);
  font-size: 3rem;
  text-align: left;
  line-height: 5vh;

  .sentence-one {
    width: fit-content;
    padding-top: 15vh;
    padding-left: 8vw;
  }

  .sentence-two {
    width: fit-content;
    padding-top: 10vh;
    padding-left: 65vw;
  }

  .sentence-three {
    width: fit-content;
    padding-top: 10vh;
    padding-left: 40vw;
  }

  .sentence-en {
    font-size: 2rem;
  }
}
</style>
