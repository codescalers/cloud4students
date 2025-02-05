<template>
  <v-app class="overflow-hidden">
    <DefaultBar />
    <DefaultView />
    <FooterComponent />
  </v-app>
</template>

<script setup>
import DefaultBar from "./AppBar.vue";
import DefaultView from "./View.vue";
import FooterComponent from "@/components/Footer.vue";

</script>
