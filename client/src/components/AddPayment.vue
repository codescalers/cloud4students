<template>
  <v-dialog v-model="dialog" max-width="800">
    <template v-slot:activator="{ props: activatorProps }">
      <BaseButton
        v-bind="activatorProps"
        class="my-5"
        text="+ Add Payment Method"
        variant="outlined"
      />
    </template>

    <PaymentCard
      title="Add a new payment method"
      @on-update="updateData"
      @on-close="dialog = false"
    />
  </v-dialog>
</template>
<script setup>
import { ref } from "vue";
import BaseButton from "@/components/Form/BaseButton.vue";
import PaymentCard from "./PaymentCard.vue";

const emit = defineEmits("updateData");
const dialog = ref(false);

function updateData(data) {
  emit("updateData", data);
}
</script>
