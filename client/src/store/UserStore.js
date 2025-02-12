import { defineStore } from "pinia";
import router from "@/router";
import userService from "@/services/userService";

export const useUserStore = defineStore("userStore", {
  state: () => ({
    user: null,
    newUser: null,
    isLoaded: false,
    notifications: [],
    maintenance: false,
    next_launch_admin: true,
    isAuthenticated: localStorage.getItem("token"),
  }),
  actions: {
    async login(email, password) {
      try {
        const res = await userService.signIn(email, password);
        const { access_token } = res.data.data;
        localStorage.setItem("token", access_token);
        this.isLoaded = true;
        return res;
      } catch (error) {
        if (error) throw error;
      }
    },
    async getUserInfo() {
      try {
        const response = await userService.getUser();
        const { user } = response.data.data;
        this.user = user;
      } catch (error) {
        if (error.response.status == 401) {
          localStorage.removeItem("token");
          router.push("/");
          return error;
        }
        return error;
      } finally {
        this.isLoaded = true;
      }
    },

    async startSSE() {
      try {
        await userService.SSE((notification) => {
          this.notifications.push(notification);
        });
      } catch (error) {
        if (error.response.status == 401) {
          localStorage.removeItem("token");
          router.push("/");
          return error;
        }
      }
    },

    async getNextLaunch() {
      try {
        const response = await userService.nextLaunch();
        const { launched } = response.data.data;
        this.next_launch_admin = launched;
      } catch (error) {
        return error;
      }
    },

    async setNextLaunch(value) {
      try {
        const response = await userService.setNextLaunch(value);
        return response;
      } catch (error) {
        return error;
      }
    },

    async checkMaintenance() {
      try {
        const response = await userService.maintenance();
        const { active } = response.data.data;
        this.maintenance = { active };
      } catch (error) {
        return error;
      }
    },
  },
  getters: {
    isUserLoaded: (state) => state.isLoaded,
    isNextLaunchEnabled: (state) => state.next_launch_admin,
    getTotalBalance: (state) => state.user? state.user.balance + state.user.voucher_balance : 0,
    isAdmin: (state) => state.user.admin,
  },
});
