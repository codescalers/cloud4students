/**
 * main.js
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

// Components
import App from "./App.vue";

// Composables
import { createApp } from "vue";
import { createPinia } from "pinia";
import moshaToast from "mosha-vue-toastify";
import Default from "./layouts/default/Default.vue";
import NoNavbar from "./layouts/NoNavbar.vue";
import "mosha-vue-toastify/dist/style.css";

const pinia = createPinia();
import { useUserStore } from "./store/UserStore";
// Plugins
import { registerPlugins } from "@/plugins";

const app = createApp(App);

registerPlugins(app);

app.component("Default-Layout", Default);
app.component("No-Navbar-Layout", NoNavbar);

app.use(pinia);
const store = useUserStore();
store.startSSE();

app.use(moshaToast).mount("#app");
