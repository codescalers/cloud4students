<template>
  <component :is="layout">
    <router-view />
  </component>
</template>

<script setup>
import { computed } from "vue";
import { useRoute } from "vue-router";

const route = useRoute();

const layout = computed(() => {
  const routeLayout = route.meta.layout;
  return routeLayout === "Default" ? "Default-Layout" : "No-Navbar-Layout";
});
</script>
<style>
body {
  background-color: #212121;
}
.v-container--fluid {
  max-width: 100% !important;
}
</style>
