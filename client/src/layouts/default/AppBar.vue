<template>
  <v-app-bar>
    <v-container>
      <v-navigation-drawer v-model="drawer" app>
        <v-list>
          <v-list-tile
            v-for="item in menuItems"
            :key="item.title"
            :to="item.path"
          >
            <v-list-tile-content>{{ item.title }}</v-list-tile-content>
          </v-list-tile>
        </v-list>
      </v-navigation-drawer>

      <v-toolbar app>
        <span class="hidden-sm-and-up">
          <v-toolbar-side-icon @click="drawer = !drawer"> </v-toolbar-side-icon>
        </span>
        <v-toolbar-title>
          <router-link to="/" style="cursor: pointer">
            <v-img src="@/assets/logo_c4all.png" width="70" />
          </router-link>
        </v-toolbar-title>
        <v-spacer></v-spacer>
        <v-toolbar-items class="hidden-xs-only">
          <v-btn
            class="text-capitalize"
            flat
            v-for="item in navItems"
            :key="item.title"
            :to="item.path"
          >
            {{ item.title }}
          </v-btn>
          <v-btn v-if="user.admin" class="text-capitalize" flat to="/admin">
            Admin
          </v-btn>
        </v-toolbar-items>
        <v-btn class="text-capitalize" flat>
          Balance: ${{ user.balance }}
        </v-btn>
        <v-menu id="notifications" location="bottom">
          <template v-slot:activator="{ props }">
            <v-btn class="mr-1 text-capitalize" v-bind="props">
              <v-badge
                v-if="notifications.length > 0"
                color="secondary"
                :content="notifications.length"
              >
                <v-icon size="25">mdi-bell</v-icon>
              </v-badge>
              <v-icon v-else size="25">mdi-bell</v-icon>
            </v-btn>
          </template>

          <v-list
            max-height="400px"
            v-if="notifications.length > 0"
            density="compact"
          >
            <v-list-subheader>Unseen</v-list-subheader>
            <v-list-item
              v-for="item in notifications"
              :key="item.id"
              class="tile text-white"
            >
              <v-list-item-title>
                <router-link
                  style="padding: 15px"
                  :to="item.type == 'vms' ? '/vm' : '/'"
                  class="d-flex text-white text-decoration-none"
                  @click="
                    seen(item.id);
                    setActive(item.type == 'vms' ? 2 : 3, item.type);
                  "
                >
                  <span @click="seen(item.id)">{{ item.msg }}</span>
                </router-link>
              </v-list-item-title>
            </v-list-item>
          </v-list>

          <v-list v-if="notifications.length == 0">
            <v-list-item>
              <v-list-item-title>
                <span>You don't have any notifications yet</span>
              </v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
        <v-menu>
          <template v-slot:activator="{ props }">
            <v-btn class="text-capitalize" v-bind="props">
              <v-icon size="25" class="mr-2">mdi-account-circle-outline</v-icon>
              {{ user.first_name }}
            </v-btn>
          </template>
          <v-list>
            <v-list-item>
              <v-list-item-title>
                <router-link
                  v-for="item in menuItems"
                  :key="item.title"
                  :to="item.path"
                  class="d-flex my-3 text-white text-decoration-none"
                >
                  <span @click="checkTitle(item.title)">{{ item.title }}</span>
                </router-link>
              </v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </v-toolbar>
    </v-container>
  </v-app-bar>
  <Toast ref="toast" />
</template>

<script setup>
import { ref, onMounted } from "vue";
import userService from "@/services/userService";
import { useRoute } from "vue-router";
import Toast from "@/components/Toast.vue";
import { storeToRefs } from "pinia";
import { useUserStore } from "@/store/UserStore";

const route = useRoute();
const drawer = ref(false);
const isActive = ref(0);
const notifications = ref([]);
const toast = ref(null);
const { user } = storeToRefs(useUserStore());
const navItems = ref([
  { title: "Home", path: "/" },
  { title: "Virtual Machines", path: "/vm" },
]);

const menuItems = ref([
  {
    title: "Account Management",
    path: "/account",
  },
  {
    title: "Sign Out",
    path: "/logout",
  },
]);

const setActive = (index, item) => {
  if (item == null) {
    isActive.value = null;
  } else {
    isActive.value = index;
  }
};

const checkTitle = (title) => {
  if (title == "Sign Out") {
    userService.logout();
  }
};

const getNotifications = () => {
  userService
    .getNotifications()
    .then((response) => {
      const { data } = response.data;
      notifications.value = data.filter((item) => !item.seen);
    })
    .catch((response) => {
      const { err } = response.response.data;
      toast.value.toast(err, "#FF5252");
    });
};

const seen = (id) => {
  userService.seenNotification(id).catch((response) => {
    const { err } = response.response.data;
    toast.value.toast(err, "#FF5252");
  });
};

onMounted(async () => {
  await getNotifications();
  const token = localStorage.getItem("token");
  const response = await fetch("http://localhost:3000/v1/notification/stream", {
    method: "GET",
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (!response.ok) {
    console.error("Failed to connect to SSE endpoint");
    return;
  }

  const reader = response.body.pipeThrough(new TextDecoderStream()).getReader();
  let finished = false;

  while (!finished) {
    const { value, done } = await reader.read();
    if (done) break;
    notifications.value.push(value);
  }

  if (route.redirectedFrom) checkTitle(route.redirectedFrom.name);
});
</script>

<style>
.v-toolbar {
  background: #212121 !important;
}
</style>
