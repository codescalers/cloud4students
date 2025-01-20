<template>
  <v-container>
    <v-row class="d-flex align-center justify-center">
      <v-col cols="12" md="6">
        <v-img :src="nextLaunchImg" class="mx-auto my-12" />
        <h5 class="text-h5 text-md-h4 font-weight-bold text-center my-5">
          We're working hard to bring you the next big thing. Stay tuned!
        </h5>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import nextLaunchImg from "@/assets/next_launch.png";
const router = useRouter();
const nextLaunch = ref(localStorage.getItem("nextlaunch") == "true");
if (!nextLaunch.value) {
  router.push({ name: "Home" });
}
</script>
