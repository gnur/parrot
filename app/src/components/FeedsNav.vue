<template>
  <nav class="panel">
    <p class="panel-heading">
      feeds
    </p>
    <a
      v-for="f in feeds"
      :key="f.source"
      class="panel-block"
      :class="{ 'is-active': activeFeed === f.source }"
      @click="setActive(f)"
    >
      {{ f.source }}
      <span class="tag is-primary has-auto-margin-left">{{ f.logs.length }}</span>
    </a>
  </nav>
</template>

<script>
export default {
  name: "feeds-nav",
  props: ["feeds"],
  data() {
    return {
      activeFeed: ""
    };
  },
  methods: {
    setActive(f) {
      this.activeFeed = f.source;

      this.$emit("select", f);
    }
  }
};
</script>

<style lang="scss">
.has-auto-margin-left {
  margin-left: auto;
}
</style>
