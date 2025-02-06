<template>
  <v-row>
    <v-col cols="12" v-if="cards">
      <v-card variant="flat">
        <div class="d-flex flex-no-wrap justify-space-between">
          <div>
            <v-card-title class="text-h5">
              Top-up your payment account
            </v-card-title>

            <v-card-subtitle
              >Top-up your payment account yo charge your balance
            </v-card-subtitle>
          </div>

          <v-sheet class="d-flex flex-column pa-5 align-center" color="primary">
            <p class="font-weight-bold">${{ balance }}</p>
            <span class="text-disabled text-capitalize mt-1"
              >current balance</span
            >
          </v-sheet>
        </div>
      </v-card>
    </v-col>

    <v-col cols="12" md="6">
      <v-card variant="flat">
        <v-card-title class="text-h6 text-capitalize"
          >use your voucher
        </v-card-title>

        <v-form v-model="verifyVoucher" @submit.prevent="activateVoucher">
          <div class="d-flex w-75">
            <BaseInput
              v-model="voucher"
              placeholder="Voucher code"
              class="ml-4"
              :rules="[required]"
              :loading="loading"
            />
            <BaseButton
              type="submit"
              :disabled="!verifyVoucher"
              color="secondary"
              text="Verify"
              class="mx-2"
            />
          </div>
        </v-form>
      </v-card>
    </v-col>

    <v-col cols="12" v-if="cards">
      <v-card variant="flat">
        <v-card-title class="text-h6 text-capitalize"
          >Choose your Top-up amount:</v-card-title
        >
        <v-form v-model="verify" @submit.prevent="chargeBalance">
          <div class="d-flex align-center w-50 mb-5">
            <v-icon size="25"> mdi-currency-usd </v-icon>
            <BaseInput
              placeholder="Enter custom amount"
              class="mr-5"
              v-model="amount"
              hide-details
              :disabled="selection"
            />

            <v-chip-group
              v-model="selection"
              :disabled="amount > 0"
              selected-class="bg-secondary"
            >
              <v-chip
                v-for="amount in amounts"
                :key="amount"
                class="px-5"
                :value="amount"
                label
                variant="outlined"
                >${{ amount }}</v-chip
              >
            </v-chip-group>
          </div>
          <v-divider></v-divider>
          <p class="text-caption my-5 text-medium-emphasis">
            *By clicking "Charge your balance" I give permission for Cloud4All
            to initiate this payment transaction from credit card.
          </p>

          <div class="d-flex justify-space-between">
            <p class="text-caption my-5 text-medium-emphasis">
              I also agree and consent to
              <span class="font-weight-bold"
                >Cloud4All Terms and Condition</span
              >
            </p>
            <BaseButton
              type="submit"
              color="secondary"
              text="Charge your Balance"
              :disabled="amount == null && !selection"
            />
          </div>
        </v-form>
      </v-card>
    </v-col>
    <Toast ref="toast" />
  </v-row>
</template>
<script setup>
import { ref, watch } from "vue";
import BaseInput from "./Form/BaseInput.vue";
import BaseButton from "./Form/BaseButton.vue";
import userService from "@/services/userService";
import Toast from "./Toast.vue";
import { storeToRefs } from "pinia";
import { useUserStore } from "@/store/UserStore";

defineProps({
  cards: {
    type: Boolean,
  },
});
const selection = ref(null);
const verify = ref(false);
const loading = ref(false);
const verifyVoucher = ref(false);
const toast = ref(null);
const store = useUserStore();
const { user } = storeToRefs(store);
const balance = ref(user.value.balance);
const defaultCard = ref(user.value.stripe_default_payment_id);
const voucher = ref();
const amount = ref(null);
const amounts = ref(["50", "100"]);

function required(v) {
  return !!v || "Field is required";
}
// FIXME after voucher is activated
function activateVoucher() {
  loading.value = true;
  userService
    .activateVoucher(voucher.value)
    .then((response) => {
      toast.value.toast(response.data.msg, "#4caf50");
    })
    .catch((response) => {
      const { err } = response.response.data;
      toast.value.toast(err, "#FF5252");
    })
    .finally(() => {
      loading.value = false;
      voucher.value = "";
      verifyVoucher.value = false;
    });
}
async function chargeBalance() {
  userService
    .chargeBalance(
      amount.value ? +amount.value : +selection.value,
      defaultCard.value
    )
    .then((response) => {
      toast.value.toast(response.data.msg, "#4caf50");
    })
    .catch((response) => {
      const { err } = response.response.data;
      toast.value.toast(err, "#FF5252");
    })
    .finally(async () => {
      amount.value = null;
      selection.value = null;
      await store.getUserInfo();
    });
}

watch(
  user,
  (newVal) => {
    if (newVal) {
      balance.value = newVal.balance;
    }
  },
  { immediate: true }
);
</script>
