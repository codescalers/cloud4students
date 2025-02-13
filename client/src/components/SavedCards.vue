<template>
  <h2 class="font-weight-bold my-5">Your Saved Cards</h2>

  <v-card
    border="opacity-25 sm"
    class="my-5"
    v-for="card in cards"
    :key="card.id"
  >
    <v-card-actions>
      <v-list-item class="w-100">
        <template v-slot:prepend>
          <v-avatar rounded="0">
            <v-img :src="`/src/assets/cards_logos/${card.brand}.png`" />
          </v-avatar>
        </template>

        <v-list-item-title class="text-capitalize">
          {{ card.brand }} Card ending in {{ card.last_4 }}
        </v-list-item-title>

        <v-list-item-subtitle class="text-medium-emphasis">
          Expires
          {{ card.exp_month < 10 ? `0${card.exp_month}` : card.exp_month }}/{{
            card.exp_year.toString().slice(-2)
          }}
        </v-list-item-subtitle>

        <template v-slot:append>
          <BaseButton
            class="mx-5 text-normal"
            text="Set as default"
            size="small"
            rounded
            color="info"
            :disabled="getDefaultCard(card.payment_method_id)"
            @click="setDefaultCard(card.payment_method_id)"
          />

          <div class="justify-self-end">
            <v-dialog v-model="dialog" max-width="500">
              <template v-slot:activator="{ props: activatorProps }">
                <v-icon
                  v-bind="activatorProps"
                  class="me-1"
                  icon="mdi-trash-can-outline"
                  :disabled="cards.length == 1"
                  @click="setItemToDelete(card)"
                />
              </template>

              <template v-slot:default="{ isActive }">
                <v-card class="pa-2">
                  <v-card-title>Delete</v-card-title>
                  <v-divider />
                  <v-card-text>
                    Are you sure you need to delete card?
                  </v-card-text>

                  <v-card-actions>
                    <v-spacer></v-spacer>

                    <BaseButton
                      text="Cancel"
                      @click="isActive.value = false"
                      variant="outlined"
                    />
                    <BaseButton
                      text="Yes Delete"
                      color="error"
                      @click="deleteCard(itemToDelete)"
                    />
                  </v-card-actions>
                </v-card>
              </template>
            </v-dialog>
          </div>
        </template>
      </v-list-item>
    </v-card-actions>
  </v-card>
  <Toast ref="toast" />
</template>
<script setup>
import { ref, watch } from "vue";
import BaseButton from "./Form/BaseButton.vue";
import userService from "@/services/userService";
import Toast from "./Toast.vue";
import { storeToRefs } from "pinia";
import { useUserStore } from "@/store/UserStore";

const props = defineProps({
  cards: {
    type: Array,
  },
});

const emit = defineEmits("updateData");

const toast = ref();
const dialog = ref(false);
const { user } = storeToRefs(useUserStore());
const itemToDelete = ref();

function setDefaultCard(id) {
  userService
    .setDefaultCard(id)
    .then((response) => {
      user.value.stripe_default_payment_id = id;
      toast.value.toast(response.data.msg, "#4caf50");
    })
    .catch((error) => {
      toast.value.toast(error.message, "#FF5252");
    });
}

function getDefaultCard(id) {
  return user.value.stripe_default_payment_id == id;
}

function setItemToDelete(card) {
  itemToDelete.value = card.id;
  dialog.value = true;
}

function deleteCard(id) {
  userService
    .deleteCard(id)
    .then((response) => {
      emit("updateData", response.data);
      toast.value.toast(response.data.msg, "#4caf50");
    })
    .catch((response) => {
      const { err } = response.response.data;
      toast.value.toast(err, "#FF5252");
    })
    .finally(() => {
      dialog.value = false;
    });
}

watch(
  props.cards,
  () => {
    if (props.cards.length == 1) {
      setDefaultCard(props.cards[0].payment_method_id);
    }
  },
  { immediate: true }
);
</script>
