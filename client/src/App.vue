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
const isAuthenticated = localStorage.getItem('token');

const layout = computed(() => {
  const routeLayout = route.meta.layout || 'defaultLayout';

  if (route.path === "/") {
    return isAuthenticated ? 'Default-Layout' : 'No-Navbar-Layout';
  }
  return routeLayout === "Default" ? 'Default-Layout' : 'No-Navbar-Layout';
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
