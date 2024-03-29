<template>
	<v-container>
		<Toast ref="toast" />
		<h5 class="text-h5 text-md-h4 font-weight-bold text-center my-10 secondary">
			Verification Code
		</h5>
		<v-row justify="center">
			<v-col cols="12" sm="6">
				<v-form @submit.prevent="onSubmit">
					<v-hover v-slot="{ isHovering, props }" open-delay="200">
						<v-img :style="isHovering
							? 'transform:scale(1.1);transition: transform .5s;'
							: 'transition: transform .5s;'
							" transition="transform .2s" src="@/assets/verification_code.png" class="w-50 mx-auto"
							:class="{ 'on-hover': isHovering }" v-bind="props" />
					</v-hover>

					<div>
						<v-otp-input v-model="otp" length="4"></v-otp-input>
						<div class="w-50 mx-auto text-center my-5">
							<v-btn block class="my-5" style="
                  background-color: transparent;
                  box-shadow: none;
                  text-transform: unset !important;
                " :disabled="countDown > 0" @click="resetHandler">
								Re-send confirmation code</v-btn>
							<span class="block">00:{{ countDown }}</span>
						</div>
					</div>

					<v-btn type="submit" block :disabled="otp.length != 4" :loading="loading" variant="flat" color="primary"
						class="text-capitalize mx-auto bg-primary">
						Confirm Code
					</v-btn>
				</v-form>
			</v-col>
		</v-row>
	</v-container>
</template>

<script>
import { ref, watchEffect } from "vue";
import { useRouter, useRoute } from "vue-router";
import axios from "axios";
import Toast from "@/components/Toast.vue";

export default {
	components: {
		Toast,
	},

	setup() {
		const route = useRoute();
		const router = useRouter();
		const countDown = ref(route.query.timeout)
		const loading = ref(false);
		const otp = ref("");
		const toast = ref(null);

		watchEffect(() => {
			if (countDown.value > 0) {
				setTimeout(() => {
					countDown.value--;
				}, 1000);
			}
		});

		const resetHandler = async () => {
			if (route.query.isForgetPassword) {
				await axios
					.post(window.configs.vite_app_endpoint + "/user/forgot_password", {
						email: route.query.email,
					})
					.then((response) => {
						toast.value.toast(response.data.msg);
						countDown.value = route.query.timeout;
					})
					.catch((error) => {
						toast.value.toast(error.response.data.err, "#FF5252");
					})
					.finally(() => {
						otp.value = "";
					});
			} else {
				await axios
					.post(window.configs.vite_app_endpoint + "/user/signup", {
						name: localStorage.getItem("fullName"),
						email: route.query.email,
						password: localStorage.getItem("password"),
						confirm_password: localStorage.getItem("confirm_password"),
						team_size: Number(localStorage.getItem("teamSize")),
						project_desc: localStorage.getItem("projectDescription"),
						college: localStorage.getItem("faculty"),
						ssh_key: localStorage.getItem("sshKey"),
					})
					.then((response) => {
						toast.value.toast(response.data.msg);
						countDown.value = route.query.timeout;
					})
					.catch((error) => {
						toast.value.toast(error.response.data.err, "#FF5252");
					})
					.finally(() => {
						otp.value = "";
					});
			}
		};

		const onSubmit = async () => {
			loading.value = true;

			if (route.query.isSignup) {
				await axios
					.post(
						window.configs.vite_app_endpoint + "/user/signup/verify_email",
						{
							email: route.query.email,
							code: Number(otp.value),
						}
					)
					.then(async (response) => {
						await axios.post(window.configs.vite_app_endpoint + "/user/apply_voucher", {
							vms: Number(localStorage.getItem("vms")),
							public_ips: Number(localStorage.getItem("ips")),
							reason: localStorage.getItem("projectDescription"),
						}, {
							headers: {
								Authorization: "Bearer " + response.data.data.access_token,
							}
						}).catch((error) => {
							toast.value.toast(error.response.data.err, "#FF5252");
						});
						toast.value.toast(response.data.msg);
						localStorage.removeItem("fullName");
						localStorage.removeItem("password");
						localStorage.removeItem("confirm_password");
						localStorage.removeItem("teamSize");
						localStorage.removeItem("projectDescription");
						localStorage.removeItem("faculty");
						localStorage.removeItem("sshKey");
						localStorage.removeItem("vms");
						localStorage.removeItem("ips");
						router.push({
							name: "Login",
						});
					})
					.catch((error) => {
						toast.value.toast(error.response.data.err, "#FF5252");
						loading.value = false;
					})
					.finally(() => {
						otp.value = "";
					});
			} else {
				await axios
					.post(
						window.configs.vite_app_endpoint +
						"/user/forget_password/verify_email",
						{
							email: route.query.email,
							code: Number(otp.value),
						}
					)
					.then((response) => {
						toast.value.toast(response.data.msg);
						localStorage.setItem(
							"password_token",
							response.data.data.access_token
						);

						router.push({
							name: "NewPassword",
							query: { email: route.query.email },
						});
					})
					.catch((error) => {
						toast.value.toast(error.response.data.err, "#FF5252");
						loading.value = false;
					})
					.finally(() => {
						otp.value = "";
					});
			}
		};

		return {
			onSubmit,
			otp,
			resetHandler,
			countDown,
			loading,
			toast,
		};
	},
};
</script>
<style>
.otp-input {
	max-width: 80px;
	max-height: 300px;
	padding: 5px;
	margin: 0 10px;
	font-size: 20px;
	border-radius: 4px;
	border: 1px solid rgba(0, 0, 0, 0.3);
	text-align: center;
}

/* Background color of an input field with value */
.otp-input.is-complete {
	background-color: #e4e4e4;
}

.otp-input::-webkit-inner-spin-button,
.otp-input::-webkit-outer-spin-button {
	-webkit-appearance: none;
	margin: 0;
}

input::placeholder {
	font-size: 15px;
	text-align: center;
	font-weight: 600;
}

.reset-btn {
	display: flex;
	justify-content: center;
	align-items: center;
}
</style>
