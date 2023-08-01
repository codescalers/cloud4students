// Composables
import { createRouter, createWebHistory } from "vue-router";
import Home from "@/views/Home.vue";
import About from "@/views/About.vue";
import VM from "@/views/VM.vue";
import K8s from "@/views/K8s.vue";
import Profile from "@/views/Profile.vue";
import Admin from "@/views/Admin.vue";
import Cancel from "@/views/Cancel.vue";
import Success from "@/views/Success.vue";
import NewPassword from "@/views/Newpassword.vue";
import userService from "@/services/userService.js";

const routes = [{
        path: "/",
        name: "Landing",
        component: () =>
            import ("@/views/LandingPage.vue"),
        meta: {
            requiredAuth: false,
            layout: "Default",
        },
    },
    {
        path: "/login",
        name: "Login",
        component: () =>
            import ("@/views/Login.vue"),
        meta: {
            requiredAuth: false,
            layout: "Default",
        },
    },
    {
        path: "/signup",
        name: "Signup",
        component: () =>
            import ("@/views/Signup.vue"),
        meta: {
            requiredAuth: false,
            layout: "Default",
        },
    },
    {
        path: "/forgetPassword",
        name: "ForgetPassword",
        component: () =>
            import ("@/views/Forgetpassword.vue"),
        meta: {
            requiredAuth: false,
            layout: "Default",
        },
    },
    {
        path: "/otp",
        name: "OTP",
        component: () =>
            import ("@/views/Otp.vue"),
        meta: {
            requiredAuth: false,
            layout: "Default",
        },
    },
    {
        path: "/newPassword",
        name: "NewPassword",
        component: () =>
            import ("@/views/Newpassword.vue"),
        meta: {
            requiredAuth: false,
            layout: "Default",
        },
    },
    {
        path: "/about",
        name: "About",
        component: About,
        meta: {
            requiredAuth: false,
            layout: "Default",
        },
    },
    {
        path: "/cancel",
        name: "Cancel",
        component: Cancel,
        meta: {
            requiredAuth: true,
            layout: "Default",
        },
    },
    {
        path: "/success",
        name: "Success",
        component: Success,
        meta: {
            requiredAuth: true,
            layout: "Default",
        },
    },
    {
        path: "/maintenance",
        name: "Maintenance",
        component: () =>
            import ("@/views/Maintenance.vue"),
        meta: {
            requiredAuth: false,
            layout: "Default",
        },
    },
    {
        path: "/",
        component: () =>
            import ("@/layouts/default/Default.vue"),
        meta: {
            requiredAuth: true,
        },
        children: [{
                path: "/home",
                name: "Home",
                component: Home,
                meta: {
                    requiredAuth: true,
                },
            },
            {
                path: "/profile",
                name: "Profile",
                component: Profile,
                meta: {
                    requiredAuth: true,
                },
            },
            {
                path: "/changePassword",
                name: "ChangePassword",
                component: NewPassword,
                meta: {
                    requiredAuth: true,
                },
            },
            {
                path: "/vm",
                name: "VM",
                component: VM,
                meta: {
                    requiredAuth: true,
                },
            },
            {
                path: "/k8s",
                name: "K8s",
                component: K8s,
                meta: {
                    requiredAuth: true,
                },
            },
            {
                path: "/about",
                name: "About",
                component: About,
                meta: {
                    requiredAuth: false,
                },
            },
            {
                path: "admin",
                name: "Admin",
                component: Admin,
                meta: {
                    requiredAuth: true,
                },
            },
            {
                path: "/logout",
                name: "Logout",
                redirect: "/login",
            },
        ],
    },
    {
        path: "/:pathMatch(.*)*",
        name: "PageNotFound",
        component: () =>
            import ("@/views/PageNotFound.vue"),
        meta: {
            requiredAuth: false,
            layout: "NoNavbar",
        },
    },
];

const router = createRouter({
    history: createWebHistory(process.env.BASE_URL),
    routes,
});

router.beforeEach(async(to, from, next) => {
    let token = localStorage.getItem("token");
    userService.maintenance();

    if (to.meta.requiredAuth && !token) {
        next("/login");
    } else if (to.path == "/" && token) {
        await userService.refresh_token();
        next("/home");
    } else if (to.meta.requiredAuth) {
        await userService.refresh_token();
        next();
    } else {
        next();
    }
});
export default router;