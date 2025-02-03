<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <v-data-table
          :headers="headers"
          :items="logs"
          class="d-flex justify-center elevation-1"
          :hide-default-footer="logs == 0"
        >
        </v-data-table>
      </v-col>
    </v-row>
  </v-container>
  <Toast ref="toast" />
</template>
<script setup>
import { ref, onMounted } from "vue";
import userService from "@/services/userService";
import Toast from "@/components/Toast.vue";

const logs = ref([]);
const events = ref([]);
const toast = ref(null);
const headers = ref([
  {
    title: "Event",
    key: "event",
  },
  {
    title: "Date",
    key: "date",
  },
]);

async function getAuditEvents() {
  await userService
    .getAuditEvents()
    .then((response) => {
      const { data } = response.data;
      events.value = data;
    })
    .catch((response) => {
      const { err } = response.response.data;
      toast.value.toast(err, "#FF5252");
    });
}

async function getAuditLogs() {
  await userService
    .getAuditEvents()
    .then((response) => {
      const { data } = response.data;
      logs.value = data;
      console.log(logs.value)
    })
    .catch((response) => {
      const { err } = response.response.data;
      toast.value.toast(err, "#FF5252");
    });
}

onMounted(async () => {
  await getAuditEvents();
  await getAuditLogs();
});
</script>
<style>
thead th {
  background-color: #19647e !important;
}
tbody tr {
  background-color: #474747;
}
</style>
