<template>
  <v-card class="pa-2">
    <v-card-title>Add a new payment method</v-card-title>
    <v-divider />
    <v-card-text>
      <v-row>
        <v-col cols="12" md="6">
          <label>Card Number</label>
          <div id="card-number" class="card-input mt-1"></div>
        </v-col>

        <v-col cols="12" md="3">
          <label>Expiration Date</label>
          <div id="card-expiry" class="card-input mt-1"></div>
        </v-col>

        <v-col cols="12" md="3">
          <label>Security Code</label>
          <div id="card-cvc" class="card-input mt-1"></div>
        </v-col>
      </v-row>
      <div id="card-error"></div>
      <p class="text-caption my-2 text-medium-emphasis">
        By providing your card information, you allow Cloud4All, LLC to charge
        your card for future payments in accordance with their terms.
      </p>
    </v-card-text>

    <v-divider></v-divider>

    <v-card-actions class="justify-end">
      <BaseButton text="Cancel" variant="outlined" @click="$emit('onClose')" />
      <BaseButton
        :loading="loading"
        text="Create"
        color="secondary"
        @click="generateToken"
      />
    </v-card-actions>
  </v-card>
  <Toast ref="toast" />
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from "vue";
import BaseButton from "@/components/Form/BaseButton.vue";
import { loadStripe } from "@stripe/stripe-js";
import userService from "@/services/userService";
import Toast from "./Toast.vue";

const emit = defineEmits(["onClose", "updateData"]);

const stripe = ref();
const toast = ref(null);
const loading = ref(false);

let cardNumber, cardExpiry, cardCvc;

async function generateToken() {
  loading.value = true;
  const { token, error } = await stripe.value.createToken(cardNumber);
  if (error) {
    loading.value = false;
    document.getElementById("card-error").innerHTML = error.message;
    return;
  }
  addCard(token.id, token.type);
}

async function addCard(id, type) {
  userService
    .addCard(id, type)
    .then((response) => {
      toast.value.toast(response.data.msg, "#4caf50");
    })
    .catch((response) => {
      const { err } = response.response.data;
      toast.value.toast(err, "#FF5252");
    })
    .finally(() => {
      loading.value = false;
      emit("onClose");
      userService
        .getCards()
        .then((response) => {
          const { data } = response.data;
          emit("updateData", data);
        })
        .catch((response) => {
          const { err } = response.response.data;
          toast.value.toast(err, "#FF5252");
        });
    });
}

onMounted(async () => {
  stripe.value = await loadStripe(`${process.env.PUBLISHABLE_KEY}`);
  const elements = stripe.value.elements();
  const style = {
    base: {
      color: "#fff",
      fontSmoothing: "antialiased",
      "::placeholder": {
        color: "#aab7c4",
      },
    },
    invalid: {
      color: "#fa755a",
      iconColor: "#fa755a",
    },
  };

  cardNumber = elements.create("cardNumber", { style });
  cardNumber.mount("#card-number");

  cardExpiry = elements.create("cardExpiry", { style });
  cardExpiry.mount("#card-expiry");

  cardCvc = elements.create("cardCvc", { style });
  cardCvc.mount("#card-cvc");
});

onBeforeUnmount(() => {
  if (cardNumber) cardNumber.unmount();
  if (cardExpiry) cardExpiry.unmount();
  if (cardCvc) cardCvc.unmount();
});
</script>

<style scoped>
.card-input {
  position: relative;
  display: inline-block;
  width: 100%;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  padding: 8px 12px;
  background-color: #474747;
  transition: border-color 0.3s;
}

.card-input input {
  border: none;
  outline: none;
  width: 100%;
  font-size: 16px;
  color: #212121;
}

.card-input input::placeholder {
  color: #9e9e9e;
}

.card-number:focus-within {
  border-color: #fff;
}

#card-error {
  color: red;
}
</style>
