// Composables
import { createRouter, createWebHistory } from "vue-router";
import { useUserStore } from "@/store/UserStore";
import Account from "@/views/Account.vue";
import VM from "@/views/VM.vue";
import Admin from "@/views/Admin.vue";
import NewPassword from "@/views/Newpassword.vue";
import ProfileTab from "@/views/accManagemenTabs/Profile.vue";
import PaymentsTab from "@/views/accManagemenTabs/Payments.vue";
import Invoices from "@/views/accManagemenTabs/Invoices.vue";
import ChangePassword from "@/views/accManagemenTabs/ChangePassword.vue";
import AuditLogs from "@/views/accManagemenTabs/AuditLogs.vue";
import DeleteAccount from "@/views/accManagemenTabs/DeleteAccount.vue";
import Deploy from "@/views/Deploy.vue";
import Home from "@/views/HomeWrapper.vue";
import Requests from "@/views/adminTabs/Requests.vue";
import History from "@/views/adminTabs/History.vue";

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
    meta: {
      layout: "Default",
    },
  },
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/Login.vue"),
    meta: {
      layout: "NoNavbar",
    },
  },
  {
    path: "/signup",
    name: "Signup",
    component: () => import("@/views/Signup.vue"),
    meta: {
      layout: "NoNavbar",
    },
  },
  {
    path: "/forgetPassword",
    name: "ForgetPassword",
    component: () => import("@/views/Forgetpassword.vue"),
    meta: {
      layout: "NoNavbar",
    },
  },
  {
    path: "/otp",
    name: "OTP",
    component: () => import("@/views/Otp.vue"),
    meta: {
      layout: "NoNavbar",
    },
  },
  {
    path: "/newPassword",
    name: "NewPassword",
    component: NewPassword,
    meta: {
      layout: "NoNavbar",
    },
  },
  {
    path: "/maintenance",
    name: "Maintenance",
    component: () => import("@/views/Maintenance.vue"),
    meta: {
      layout: "NoNavbar",
    },
  },
  {
    path: "/nextlaunch",
    name: "NextLaunch",
    component: () => import("@/views/NextLaunch.vue"),
    meta: {
      requiresAuth: true,
      layout: "NoNavbar",
    },
  },
  {
    path: "/changePassword",
    name: "ChangePassword",
    component: NewPassword,
    meta: {
      layout: "Default",
      requiresAuth: true,
    },
  },
  {
    path: "/vm",
    name: "VM",
    component: VM,
    meta: {
      requiresAuth: true,
      layout: "Default",
    },
  },
  {
    path: "/account",
    component: Account,
    meta: {
      layout: "Default",
      requiresAuth: true,
    },
    children: [
      {
        path: "",
        component: ProfileTab,
      },
      {
        path: "payments",
        component: PaymentsTab,
      },
      {
        path: "change-password",
        component: ChangePassword,
      },
      {
        path: "delete-account",
        component: DeleteAccount,
      },
      {
        path: "invoices",
        component: Invoices,
      },
      {
        path: "audit-logs",
        component: AuditLogs,
      },
    ],
  },
  {
    path: "/deploy",
    name: "Deploy",
    component: Deploy,
    meta: {
      layout: "Default",
      requiresAuth: true,
    },
  },
  {
    path: "/admin",
    name: "Admin",
    component: Admin,
    beforeEnter(to, from, next) {
      const store = useUserStore();

      if (!store.isAdmin) {
        next("/");
      }
      next();
    },
    meta: {
      layout: "Default",
      requiresAuth: true,
    },
    children: [
      {
        path: "",
        component: Requests,
      },
      {
        path: "history",
        component: History,
      },
    ],
  },
  {
    path: "/logout",
    name: "Logout",
    redirect: "/login",
  },
  {
    path: "/:pathMatch(.*)*",
    name: "PageNotFound",
    component: () => import("@/views/PageNotFound.vue"),
    meta: {
      layout: "NoNavbar",
    },
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

router.beforeEach(async (to, from, next) => {
  const isAuthenticated = localStorage.getItem("token") !== null;
  const store = useUserStore();

  if (to.meta.requiresAuth && !isAuthenticated) {
    next("/");
  } else {
    if (!store.user && isAuthenticated) await store.getUserInfo();
    next();
  }
});
export default router;
