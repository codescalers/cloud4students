<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <v-data-table
          :loading="loading"
          :headers="headers"
          :items="events"
          class="d-flex justify-center elevation-1"
          :hide-default-footer="events == 0"
        >
          <template #[`item.action`]="{ item }">
            {{ item.action.replace("_", ".") }}
          </template>

          <template #[`item.timestamp`]="{ item }">
            {{ formatDate(item.timestamp) }}
          </template>
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

const events = ref([]);
const toast = ref(null);
const loading = ref(false);
const headers = ref([
  {
    title: "Action",
    key: "action",
  },
  {
    title: "Timestamp",
    key: "timestamp",
  },
]);

async function getAuditEvents() {
  loading.value = true;
  await userService
    .getAuditEvents()
    .then((response) => {
      const { data } = response.data;
      events.value = data;
    })
    .catch((response) => {
      const { err } = response.response.data;
      toast.value.toast(err, "#FF5252");
    })
    .finally(() => {
      loading.value = false;
    });
}

function formatDate(date) {
  var d = new Date(date),
    month = "" + (d.getMonth() + 1),
    day = "" + d.getDate(),
    year = d.getFullYear();

  if (month.length < 2) month = "0" + month;
  if (day.length < 2) day = "0" + day;

  return [day, month, year].join("-");
}

onMounted(async () => {
  await getAuditEvents();
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
