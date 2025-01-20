<template>
  <v-container>
    <AddPayment @updateData="updateCards" />

    <v-divider />

    <SavedCards :cards="cards" @updateData="updateCards" />

    <Alerts
      v-if="cards && cards.length == 1"
      type="warning"
      text="You cannot remove your current card without first adding another valid payment method."
    />

    <TopupPayment :cards="cards.length > 0" />
  </v-container>
  <Toast ref="toast" />
</template>
<script setup>
import { ref, onMounted } from "vue";
import Toast from "@/components/Toast.vue";
import AddPayment from "@/components/AddPayment.vue";
import SavedCards from "@/components/SavedCards.vue";
import Alerts from "@/components/Alerts.vue";
import TopupPayment from "@/components/TopupPayment.vue";
import userService from "@/services/userService";

const cards = ref([]);
const toast = ref();

function getCards() {
  userService
    .getCards()
    .then((response) => {
      const { data } = response.data;
      cards.value = data;
    })
    .catch((response) => {
      const { err } = response.response.data;
      toast.value.toast(err, "#FF5252");
    });
}

function updateCards(data) {
  if (!data) return;
  getCards();
}

onMounted(() => {
  getCards();
});
</script>
