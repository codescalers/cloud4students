<template>
  <component :is="layout">
    <router-view />
  </component>
  <Toast ref="toast" />
</template>

<script setup>
import { computed, ref } from "vue";
import { useRoute } from "vue-router";
import Toast from "./components/Toast.vue";

const route = useRoute();
const toast = ref();

const layout = computed(() => {
  const NoNavbar_layout = "No-Navbar";
  return (route.meta.layout || NoNavbar_layout) + "-Layout";
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
