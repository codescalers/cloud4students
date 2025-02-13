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
          <template #[`item.metadata`]="{ item }">
            {{ item.metadata }}
          </template>

          <template #[`item.timestamp`]="{ item }">
            {{ timeAgo.format(convertDate(item.timestamp)) }}
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
import TimeAgo from "javascript-time-ago";
import en from "javascript-time-ago/locale/en";
TimeAgo.addLocale(en)

const timeAgo = ref(new TimeAgo("en-US"));
const events = ref([]);
const toast = ref(null);
const loading = ref(false);
const headers = ref([
  {
    title: "Action",
    key: "metadata",
  },
  {
    title: "Timestamp",
    key: "timestamp",
  },
]);

function convertDate(date) {
  return new Date(date);
}

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
