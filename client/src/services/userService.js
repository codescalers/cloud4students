import axios from "axios";
import router from "@/router";

const baseClient = () =>
  axios.create({
    baseURL: window.configs.vite_app_endpoint,
  });

const authClient = () =>
  axios.create({
    baseURL: window.configs.vite_app_endpoint,
    headers: {
      Authorization: "Bearer " + localStorage.getItem("token"),
    },
  });

authClient().interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      console.error("Unauthorized access - redirecting to home");
      localStorage.removeItem("token");
      router.push({ name: "Home" });
    }
    return Promise.reject(error);
  }
);

let refreshInterval;

const startTokenRefreshInterval = function (timeout) {
  clearInterval(refreshInterval);

  refreshInterval = setInterval(() => {
    this.refresh_token();
  }, (timeout - 5) * 1000);
};

export default {
  async refresh_token() {
    return await authClient()
      .post("/user/refresh_token")
      .then((response) => {
        const { refresh_token } = response.data.data;
        localStorage.setItem("token", refresh_token);
      })
      .catch((error) => {
        console.error("Failed to refresh token", error);
        localStorage.removeItem("token");
        clearInterval(refreshInterval);
      });
  },

  // user
  async getUser() {
    return await authClient().get("/user");
  },

  async signUp(first_name, last_name, email, password, confirm_password) {
    return await baseClient().post("/user/signup", {
      first_name,
      last_name,
      email,
      password,
      confirm_password,
    });
  },

  async signIn(email, password) {
    return await baseClient()
      .post("/user/signin", {
        email,
        password,
      })
      .then((res) => {
        const { timeout } = res.data.data;
        startTokenRefreshInterval.call(this, timeout);
        return res;
      });
  },

  logout() {
    clearInterval(refreshInterval);
    localStorage.removeItem("token");
  },

  async forgotPassword(email) {
    return await baseClient().post("/user/forgot_password", { email });
  },

  async signUpVerification(email, code) {
    return await baseClient().post("/user/signup/verify_email", {
      email,
      code,
    });
  },

  async applyVoucher(balance, reason) {
    return await authClient().post("/user/apply_voucher", { balance, reason });
  },

  async activateVoucher(voucher) {
    return await authClient().put("/user/activate_voucher", { voucher });
  },

  async forgetPasswordVerification(email, code) {
    return await baseClient().post("/user/forget_password/verify_email", {
      email,
      code,
    });
  },

  async updateUser(first_name, ssh_key) {
    return await authClient().put("/user", {
      first_name,
      ssh_key,
    });
  },

  async addCard(token_id, token_type) {
    return await authClient().post("/user/card", {
      token_id,
      token_type,
    });
  },

  async getCards() {
    return await authClient().get("/user/card");
  },

  async setDefaultCard(payment_method_id) {
    return await authClient().put("/user/card/default", { payment_method_id });
  },

  async deleteCard(id) {
    return await authClient().delete(`/user/card/${id}`);
  },

  async changePassword(email, password, confirm_password) {
    return await authClient().put("/user/change_password", {
      email,
      password,
      confirm_password,
    });
  },

  async newVoucher(balance, reason) {
    return await authClient().post("/user/apply_voucher", {
      balance,
      reason,
    });
  },

  async chargeBalance(amount, payment_method_id) {
    return await authClient().put("/user/charge_balance", {
      amount,
      payment_method_id,
    });
  },

  async getQuota() {
    return await authClient().get("/quota");
  },

  async deleteAccount() {
    return await authClient().delete("/user");
  },

  // Invoices
  async getInvoices() {
    return await authClient().get("/invoice");
  },

  async payInvoice(id) {
    return await authClient().put("/invoice/pay", { id });
  },

  async getInvoice(id) {
    return await authClient().get("/invoice", { id });
  },

  async downloadInvoice(id) {
    return await authClient().get(`/invoice/download/${id}`, {
      responseType: "blob",
    });
  },

  // VM
  async getVms() {
    return await authClient().get("/vm");
  },

  async validateVMName(name) {
    return await authClient().get(`/vm/validate/${name}`);
  },

  async getRegions() {
    return await authClient().get("/region");
  },

  async deployVm(name, region, resources, isPublic) {
    return await authClient().post("/vm", {
      name,
      region,
      resources,
      public: isPublic,
    });
  },

  async deleteVm(id) {
    return await authClient().delete(`/vm/${id}`);
  },

  async deleteAllVms() {
    return await authClient().delete("/vm");
  },

  // K8s
  async getK8s() {
    return await authClient().get("/k8s");
  },

  async validateK8sName(name) {
    return await authClient().get(`/k8s/validate/${name}`);
  },

  async deployK8s(master_name, resources, workers, checked) {
    return await authClient().post("/k8s", {
      master_name,
      resources,
      workers,
      public: checked,
    });
  },

  async deleteK8s(id) {
    return await authClient().delete(`/k8s/${id}`);
  },

  async deleteAllK8s() {
    return await authClient().delete("/k8s");
  },

  // Users
  async getUsers() {
    return await authClient().get("/user/all");
  },

  // Deployments
  async getDeploymentsCount() {
    return await authClient().get("/deployment/count");
  },

  // Vouchers
  async getVouchers() {
    return await authClient().get("/voucher");
  },

  async approveVoucher(id, approved) {
    return await authClient().put(`/voucher/${id}`, { approved });
  },

  async approveAllVouchers() {
    return await authClient().put("/voucher");
  },

  async generateVoucher(length, vms, public_ips) {
    return await authClient().post("/voucher", { length, vms, public_ips });
  },

  async getAuditEvents() {
    return await authClient().get("/user/event");
  },

  async getAuditLogs() {
    return await authClient().get("/user/log");
  },

  // balance
  async getBalance() {
    return await authClient().get("/balance");
  },

  // announcement
  async sendAnnouncement(subject, announcement) {
    return await authClient().post("/announcement", { subject, announcement });
  },

  // email
  async sendEmail(subject, body, email) {
    return await authClient().post("/email", { subject, body, email });
  },

  async setAdmin(email, admin) {
    return await authClient().put("/set_admin", { email, admin });
  },

  // notifications
  async getNotifications() {
    return await authClient().get("/notification");
  },

  async seenNotification(id) {
    return await authClient().put(`/notification/${id}`);
  },

  // maintenance
  async maintenance() {
    return await baseClient().get("/maintenance");
  },

  // getting nextlaunch value
  async nextLaunch() {
    return await baseClient().get("/nextlaunch");
  },

  // setting next launch value
  async setNextLaunch(value) {
    return await authClient().put("/nextlaunch", {
      launched: value,
    });
  },
};
