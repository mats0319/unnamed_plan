<template>
  <div class="content-template">
    <div class="ct-left">
      <div class="ctl-item-back" @click="backToIndex">&lt;&nbsp;返回</div>

      <template v-for="(item, index) in subPageNav">
        <div
          v-if="!item.permission || $store.state.permission >= item.permission"
          :key="index"
          class="ctl-item"
          @click="linkToSubPage(item.routerName)"
        >
          {{ item.name }}
        </div>
      </template>
    </div>

    <el-divider class="ct-divider" direction="vertical" />

    <div class="ct-right">
      <template v-if="currRouterName === ''">
        <div class="ctr-text-wrapper">
          <span class="ctr-text">{{ indexPageContent }}</span>
        </div>
      </template>

      <router-view v-else />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from "vue";
import { SubPageNavItem } from "@/ts/structure";

export default defineComponent({
  name: "ContentTemplate",
  props: {
    subPageNav: {
      type: Array as PropType<Array<SubPageNavItem>>,
      required: true,
    },
    indexPageRouterName: {
      type: String,
      required: true,
    },
    indexPageContent: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      currRouterName: "", // distinguish index page with sub-pages, empty value means on index page now
    }
  },
  created() {
    const value = sessionStorage.getItem("router_name")
    if (value) {
      this.currRouterName = value
      sessionStorage.removeItem("router_name")
    }
  },
  mounted() {
    window.addEventListener("beforeunload", () => {
      sessionStorage.setItem("router_name", this.currRouterName)
    })
  },
  methods: {
    backToIndex(): void {
      if (this.currRouterName != "") { // on a sub page
        this.currRouterName = ""
        this.$router.push({ name: this.indexPageRouterName })
      }
    },

    linkToSubPage(routerName: string): void {
      if (this.currRouterName != routerName) { // on index page or other sub-pages
        this.currRouterName = routerName
        this.$router.push({ name: routerName })
      }
    }
  }
})
</script>

<style lang="less">
.content-template {
  height: inherit;
  display: flex;

  min-height: 50rem;

  .ct-left {
    width: 20vw;
    padding-top: 4rem;

    .ctl-item-back {
      width: 7vw;
      height: 3rem;
      line-height: 3rem;
    }

    .ctl-item {
      height: 4rem;
      line-height: 4rem;
      margin: 1rem 0;
      font-size: 2rem;
    }

    .ctl-item-back:hover, .ctl-item:hover {
      cursor: pointer;
      background-color: lightgray;
    }
  }

  .ct-divider {
    height: inherit;
    min-height: 50rem;
  }

  .ct-right {
    height: 100%;
    width: calc(80vw - 17px);

    .ctr-text-wrapper {
      height: inherit;
      display: flex;
      justify-content: center;

      .ctr-text {
        height: inherit;
        font-size: 3.5rem;

        display: flex;
        align-items: center;
      }
    }
  }
}
</style>
