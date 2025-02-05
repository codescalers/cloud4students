<template>
  <component :is="isAuthenticated ? Home : LandingPage" />
</template>

<script setup>
import { computed } from "vue";
import Home from "@/views/Home.vue";
import LandingPage from "@/views/LandingPage.vue";

const isAuthenticated = computed(() => localStorage.getItem("token"));
</script>
