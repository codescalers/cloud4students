<template>
  <v-dialog v-model="dialog" max-width="800">
    <template v-slot:activator="{ props: activatorProps }">
      <BaseButton
        v-bind="activatorProps"
        class="my-5"
        text="+ Add Card"
        variant="outlined"
      />
    </template>

    <PaymentCard @updateData="updateCards" @onClose="dialog = false" />
  </v-dialog>
</template>
<script setup>
import { ref } from "vue";
import BaseButton from "@/components/Form/BaseButton.vue";
import PaymentCard from "./PaymentCard.vue";

const emit = defineEmits("updateData");
const dialog = ref(false);

function updateCards(data) {
  if (!data) return;
  emit("updateData", data);
}
</script>
