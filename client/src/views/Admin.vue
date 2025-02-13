<template>
  <v-container fluid>
    <h5 class="text-h5 text-md-h4 font-weight-bold my-5">Admin Panel</h5>
    <v-divider />

    <v-row>
      <v-col cols="12">
        <v-card flat class="my-5">
          <v-list-item>
            <template v-slot:prepend>
              <v-sheet class="px-5 py-3" border rounded
                >Balance: {{ balance }} TFT</v-sheet
              >

              <v-sheet class="px-5 py-3 d-flex align-center">
                <v-icon class="mr-2" size="35">mdi-server</v-icon>
                <span>Deployed VMs: {{ deployedResources }} </span>
              </v-sheet>

              <v-sheet class="px-5 py-3 d-flex align-center">
                <v-icon class="mr-2" size="35">mdi-ip</v-icon>
                <span>Reserved IPs: {{ reservedIPs }} </span>
              </v-sheet>
            </template>
            <template v-slot:append>
              <BaseButton
                v-if="isNextLaunchEnabled"
                @click="setNextLaunch"
                text="Disable Launch"
                color="error"
                class="mr-2"
              />
              <BaseButton
                v-else
                @click="setNextLaunch"
                text="Enable Launch"
                color="success"
                class="mr-2"
              />

              <BaseButton text="Generate a Voucher" color="success" />
            </template>
          </v-list-item>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card class="my-5">
          <v-tabs v-model="activeTab" class="tabs">
            <v-tab v-for="tab in tabs" :key="tab" :to="tab.route" exact>{{
              tab.name
            }}</v-tab>
          </v-tabs>

          <v-card-text>
            <v-tabs-window v-model="activeTab">
              <router-view
                :vouchers="vouchers"
                :pendingVouchers="pendingVouchers"
                @update-vouchers="getVouchers"
              />
            </v-tabs-window>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <Toast ref="toast" />
  </v-container>
</template>

<script setup>
import { ref, onMounted, watch } from "vue";
import BaseButton from "@/components/Form/BaseButton.vue";
import userService from "@/services/userService.js";
import Toast from "@/components/Toast.vue";
import { storeToRefs } from "pinia";
import { useUserStore } from "@/store/UserStore";
import router from "@/router";

const vouchers = ref([]);
const pendingVouchers = ref([]);
const approveAllCount = ref(0);
const userInfo = ref(null);
const users = ref([]);
const toast = ref();
const tabs = ref([
  { name: "Users Requests", route: `/admin` },
  { name: "Users History", route: `/admin/history` },
]);
const activeTab = ref(tabs.value[0]);
const loading = ref(false);
const usedResources = ref(0);
const deployedResources = ref(0);
const balance = ref(0);
const usedIPs = ref(0);
const reservedIPs = ref(0);
const announcementDialog = ref(false);
// const showUserInfo = ref(false);
const subject = ref(null);
const announcement = ref(null);
const nextLaunchDialog = ref(false);
const store = useUserStore();
const { isNextLaunchEnabled } = storeToRefs(store);
nextLaunchDialog.value = localStorage.getItem("nextlaunchadmin") == "true";

// const openUserInfo = (user) => {
//   showUserInfo.value = true;
//   userInfo.value = user;
// };

const getUsers = () => {
  userService
    .getUsers()
    .then((response) => {
      const { data } = response.data;
      users.value = data;
      users.value.map((usedData) => {
        // FIXME
        usedResources.value += usedData.used_vms;
        usedIPs.value += usedData.used_public_ips;
      });
    })
    .catch((response) => {
      const { err } = response.response.data;
      toast.value.toast(err, "#FF5252");
    });
};

const getDeploymentsCount = () => {
  userService
    .getDeploymentsCount()
    .then((response) => {
      const { data } = response.data;
      deployedResources.value += data.vms;
      reservedIPs.value += data.ips;
    })
    .catch((response) => {
      const msg = response.message;
      toast.value.toast(msg, "#FF5252");
    });
};

const getBalance = () => {
  userService
    .getBalance()
    .then((response) => {
      const { data } = response.data;
      balance.value = data;
    })
    .catch((response) => {
      const { err } = response.response.data;
      toast.value.toast(err, "#FF5252");
    });
};

async function getVouchers() {
  loading.value = true;

  try {
    const response = await userService.getVouchers();
    const { data } = response.data;

    await updateVouchers(data);

    vouchers.value = data.filter(
      (voucher) => voucher.approved || voucher.rejected
    );
    pendingVouchers.value = data.filter(
      (voucher) => !voucher.approved && !voucher.rejected
    );
  } catch (error) {
    const { err } = error.response.data;
    toast.value.toast(err, "#FF5252");
  } finally {
    loading.value = false;
  }
}

async function updateVouchers(data) {
  const updatePromises = data.map(async (voucher) => {
    await new Promise((resolve) => setTimeout(resolve, 10));

    if (voucher.approved && voucher.rejected) {
      approveAllCount.value++;
    }

    if (users.value && voucher.user_id) {
      userInfo.value = users.value.find((user) => user.ID === voucher.user_id);

      if (userInfo.value && voucher.user_id === userInfo.value.ID) {
        Object.assign(voucher, {
          email: userInfo.value.email,
          name: userInfo.value.first_name,
        });
      }
    }
  });

  await Promise.all(updatePromises);
}

watch(announcementDialog, (val) => {
  if (val) {
    announcement.value = "";
    subject.value = "";
  }
});

// TODO generate voucher
// const generateVoucher = async () => {
//   var { valid } = await form.value.validate();
//   if (!valid) return;

//   userService
//     .generateVoucher(+length.value, +vms.value, +ips.value)
//     .then((response) => {
//       const { data, msg } = response.data;
//       message.value = msg;
//       voucher.value = data.voucher;
//     })
//     .catch((response) => {
//       toast.value.toast(response.response.data.err, "#FF5252");
//     })
//     .finally(() => {
//       getVouchers();
//       dialog.value = false;
//     });
// };

// const sendAnnouncement = async () => {
//   var { valid } = await form.value.validate();
//   if (!valid) return;

//   userService
//     .sendAnnouncement(subject.value, announcement.value)
//     .then((response) => {
//       const { msg } = response.data;
//       toast.value.toast(msg, "#388E3C");
//     })
//     .catch((response) => {
//       toast.value.toast(response.response.data.err, "#FF5252");
//     })
//     .finally(() => {
//       announcementDialog.value = false;
//     });
// };
function setNextLaunch() {
  store
    .setNextLaunch(!isNextLaunchEnabled.value)
    .then((response) => {
      const { msg } = response.data;
      toast.value.toast(msg, "#388E3C");
      store.next_launch_admin = false;
      router.push("/nextlaunch");
    })
    .catch((response) => {
      const { err } = response.response.data;
      toast.value.toast(err, "#FF5252");
    });
}

onMounted(() => {
  getUsers();
  getVouchers();
  getBalance();
  getDeploymentsCount();
});
</script>
