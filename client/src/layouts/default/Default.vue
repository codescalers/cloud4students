<template>
  <v-app class="overflow-hidden">
    <div v-if="user">
      <DefaultBar />
      <DefaultView />
      <FooterComponent />
    </div>
    <v-progress-circular indeterminate v-else></v-progress-circular>
  </v-app>
</template>

<script setup>
import DefaultBar from "./AppBar.vue";
import DefaultView from "./View.vue";
import FooterComponent from "@/components/Footer.vue";
import { storeToRefs } from "pinia";
import { useUserStore } from "@/store/UserStore";

const store = useUserStore();
const { user } = storeToRefs(store);
</script>
