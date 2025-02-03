// Composables
import { createRouter, createWebHistory } from "vue-router";
import { useUserStore } from "@/store/UserStore";
import Account from "@/views/Account.vue";
import VM from "@/views/VM.vue";
import Admin from "@/views/Admin.vue";
import NewPassword from "@/views/Newpassword.vue";
import ProfileTab from "@/views/tabs/Profile.vue";
import PaymentsTab from "@/views/tabs/Payments.vue";
import Invoices from "@/views/tabs/Invoices.vue";
import ChangePassword from "@/views/tabs/ChangePassword.vue";
import AuditLogs from "@/views/tabs/AuditLogs.vue";
import DeleteAccount from "@/views/tabs/DeleteAccount.vue";
import Deploy from "@/views/Deploy.vue";
import Home from "@/views/Home.vue";

const routes = [
  {
    path: "/home",
    name: "Landing",
    component: () => import("@/views/LandingPage.vue"),
    meta: {
      layout: "NoNavbar",
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
    path: "/",
    name: "Home",
    component: Home,
    meta: {
      layout: "Default",
      requiresAuth: true,
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
    meta: {
      layout: "Default",
      requiresAuth: true,
    },
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
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth);
  const isAuthenticated = localStorage.getItem("token");
  const store = useUserStore();

  if (!store.isStoreLoaded) await store.getUserInfo();
  if (store.maintenance) return "/maintenance";
  if (store.next_launch) return "/nextlaunch";

  if (requiresAuth && !isAuthenticated) {
    next("/home");
  } else if (to.path === "/login" && isAuthenticated) {
    next({ name: "Home" });
  }
  next();
});

export default router;
