<template>
  <div class="user">
    list&#58;<br />
    1. func bar<br />
    1.1. modify user name or password, need re-login<br />
    1.2. modify permission, need permission >= 6, can re-set permission range: (0, myP)<br />
    1.3. lock user, need permission >= 8, need target user permission less than my<br />
    1.4. unlock user, need permission >= 8, need target user permission less than my<br />
    2. user table<br />
    tips: cols reference User table in db, only admin( p >= 6 ) can see permission<br />
    tips: can only list permission less than my
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class User extends Vue {
  private mounted() {
    // placeholder
  }
}
</script>

<style lang="scss">
.user {
  text-align: left;
  font-size: 1.5rem;
}
</style>
