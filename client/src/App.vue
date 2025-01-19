<template>
  <component :is="layout">
    <router-view />
  </component>
  <Toast ref="toast" />
</template>

<script setup>
import { computed, onMounted, provide, ref } from "vue";
import { useRoute } from "vue-router";
import userService from "./services/userService";
import Toast from "./components/Toast.vue";
import router from "./router";

const route = useRoute();
const fetchedUser = ref({});
const toast = ref();

const layout = computed(() => {
  const NoNavbar_layout = "No-Navbar";
  return (route.meta.layout || NoNavbar_layout) + "-Layout";
});

const getUser = () => {
  userService
    .getUser()
    .then((response) => {
      const { user } = response.data.data;
      fetchedUser.value = user;
    })
    .catch(() => {
      const token = localStorage.getItem("token");
      token ? userService.refresh_token() : router.push("/");
    });
};
provide("user", fetchedUser);

onMounted(() => {
  getUser();
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
