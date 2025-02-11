<template>
  <v-container>
    <v-row class="d-flex justify-end my-5">
      <BaseButton
        @click="downloadAllInvoices"
        text="Download All"
        prepend-icon="mdi-download"
        :disabled="invoices == 0"
        color="secondary"
      />
    </v-row>
    <v-row>
      <v-col cols="12">
        <v-data-table
          :loading="loading"
          :headers="headers"
          :items="invoices"
          class="d-flex justify-center elevation-1"
          :hide-default-footer="invoices == 0"
        >
          <template #[`item.created_at`]="{ item }">
            {{ formatDate(item.created_at) }}
          </template>

          <template #[`item.invoice`]="{ item }">
            {{ item.userID }}
          </template>

          <template #[`item.download`]="{ item }">
            <BaseButton
              @click="downloadInvoice(item.id)"
              text="Download"
              prepend-icon="mdi-download"
              variant="text"
            />
          </template>
          <template #no-data>
            <p class="text-capitalize">{{ message }}</p>
          </template>
        </v-data-table>
      </v-col>
    </v-row>
    <Toast ref="toast" />
  </v-container>
</template>
<script setup>
import { ref, onMounted } from "vue";
import JSZip from "jszip";
import { saveAs } from "file-saver";
import BaseButton from "@/components/Form/BaseButton.vue";
import userService from "@/services/userService";
import Toast from "@/components/Toast.vue";
import { storeToRefs } from "pinia";
import { useUserStore } from "@/store/UserStore";

const invoices = ref();
const toast = ref(null);
const message = ref();
const store = useUserStore();
const { user } = storeToRefs(store);
const userID = ref(user.value.id);
const loading = ref(false);
const headers = ref([
  {
    title: "Date",
    key: "created_at",
  },
  {
    title: "Invoice Number",
    key: "invoice",
  },
  {
    title: "Download",
    key: "download",
  },
]);

function formatDate(date) {
  var d = new Date(date),
    month = "" + (d.getMonth() + 1),
    day = "" + d.getDate(),
    year = d.getFullYear();

  if (month.length < 2) month = "0" + month;
  if (day.length < 2) day = "0" + day;

  return [day, month, year].join("-");
}

// TODO handle invoices InvoiceID
function getInvoices() {
  loading.value = true;
  userService
    .getInvoice(userID.value)
    .then((response) => {
      const { data, msg } = response.data;
      invoices.value = data;
      message.value = msg;
    })
    .catch((response) => {
      const { err } = response.response.data;
      toast.value.toast(err, "#FF5252");
    })
    .finally(() => (loading.value = false));
}

function downloadInvoice(id) {
  userService
    .downloadInvoice(id)
    .then(async (response) => {
      const blob = new Blob([response.data], { type: "application/pdf" });
      const url = window.URL.createObjectURL(blob);

      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", "download.pdf");

      document.body.appendChild(link);
      link.click();

      link.parentNode.removeChild(link);
      window.URL.revokeObjectURL(url);
    })
    .catch((response) => {
      toast.value.toast(response, "#FF5252");
    });
}

async function downloadAllInvoices() {
  const zip = new JSZip();
  const promises = [];

  for (const invoice of invoices.value) {
    promises.push(
      userService
        .downloadInvoice(invoice.id)
        .then((response) => {
          const blob = new Blob([response.data], { type: "application/pdf" });
          zip.file(`invoice_${invoice.id}.pdf`, blob);
        })
        .catch((error) => {
          toast.value.toast(
            `Error downloading invoice ${invoice.id}: ${error}`,
            "#FF5252"
          );
        })
    );
  }

  await Promise.all(promises);

  zip.generateAsync({ type: "blob" }).then((content) => {
    saveAs(content, "invoices.zip");
  });
}

onMounted(() => getInvoices());
</script>
<style>
tbody tr {
  background-color: #474747;
}

.v-btn--disabled.bg-default {
  background-color: transparent !important;
}
</style>
