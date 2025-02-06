import { defineStore } from "pinia";
import router from "@/router";
import userService from "@/services/userService";

export const useUserStore = defineStore("userStore", {
  state: () => ({
    user: null,
    newUser: null,
    isLoaded: false,
    maintenance: false,
    next_launch: false,
    next_launch_admin: false,
  }),
  actions: {
    async login(email, password) {
      try {
        const res = await userService.signIn(email, password);
        const { access_token } = res.data.data;
        localStorage.setItem("token", access_token);
        return res;
      } catch (error) {
        return error;
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
          router.push("/login");
          return error
        }
        return error
      } finally {
        this.isLoaded = true;
      }
    },

    async getNextLaunch() {
      try {
        const response = await userService.nextLaunch();
        const { launched } = response.data.data;
        this.next_launch = launched;
        this.next_launch_admin = launched;
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
  },
});
