<template>
  <v-app v-if="isUserLoaded" class="overflow-hidden">
    <DefaultBar />
    <DefaultView />
    <FooterComponent />
  </v-app>
</template>

<script setup>
import DefaultBar from "./AppBar.vue";
import DefaultView from "./View.vue";
import FooterComponent from "@/components/Footer.vue";
import { storeToRefs } from "pinia";
import { useUserStore } from "@/store/UserStore";

const store = useUserStore();
const { isUserLoaded } = storeToRefs(store);
</script>
