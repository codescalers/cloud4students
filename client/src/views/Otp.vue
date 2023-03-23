<template>
  <v-container>
    <Toast ref="toast" />

    <h5 class="text-h5 text-md-h4 text-center mt-10 mb-0 secondary">
      Verfication Code
    </h5>
    <v-row justify="center">
      <v-col cols="12" sm="6">
        <v-form v-model="verify" @submit.prevent="onSubmit">
          <v-hover v-slot="{ isHovering, props }" open-delay="200">
            <v-img :style="isHovering ? 'transform:scale(1.1);transition: transform .5s;' : 'transition: transform .5s;'"
              transition="transform .2s" contain height="500" src="@/assets/otp.png" :class="{ 'on-hover': isHovering }"
              v-bind="props" />
          </v-hover>

          <div>
            <v-otp-input ref="otpInput" class="justify-center" input-classes="otp-input" separator="-" :num-inputs="4"
              style="grid-area: unset;" :should-auto-focus="true" :is-input-num="true"
              :conditionalClass="['one', 'two', 'three', 'four']" :placeholder="['', '', '', '']"
              @on-change="handleOnChange" @on-complete="handleOnComplete" />
            <div class="float-sm-center my-5 reset-btn">
              00:{{ countDown }}
              <v-btn style="background-color:transparent;box-shadow:none;text-transform: unset !important;"
                :disabled="countDown > 0" @click="resetHandler"> Re-send confirmation code</v-btn>
            </div>
          </div>

          <v-btn min-width="228" size="x-large" type="submit" block :disabled="!verify" :loading="loading" variant="flat"
            color="primary" class="text-capitalize mx-auto bg-primary">
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
import VOtpInput from 'vue3-otp-input';
import axios from "axios";
import Toast from "@/components/Toast.vue";

export default {
  components: {
    VOtpInput,
    Toast,
  },

  setup() {
    const route = useRoute();
    const router = useRouter();
    const verify = ref(false);
    const countDown = ref(30)
    const otpInput = ref(null);
    const loading = ref(false);
    const otp = ref(null)
    const toast = ref(null);


    const handleOnComplete = (value) => {
      verify.value = true;
      otp.value = value;

    };

    const handleOnChange = (value) => {
      console.log('OTP changed: ', value);
    };

    const clearInput = () => {
      otpInput.value.clearInput()
    };
    watchEffect(() => {
      if (countDown.value > 0) {
        setTimeout(() => {
          countDown.value--;
        }, 1000);
      }
    });

    const resetHandler = () => {

      if (route.query.isForgetpassword) {
        axios
          .post(import.meta.env.VITE_API_ENDPOINT+"/user/forgot_password", {
            email: route.query.email,
          })
          .then((response) => {

            toast.value.toast(response.data.msg);
            countDown.value = 30;


          })
          .catch((error) => {
            toast.value.toast(error.response.data.err, "#FF5252");

          });
      } else {

        axios
          .post(import.meta.env.VITE_API_ENDPOINT+"/user/signup", {
            name: localStorage.getItem('fullname'),
            email: route.query.email,
            password: localStorage.getItem('password'),
            confirm_password: localStorage.getItem('confirm_password'),
          })
          .then((response) => {
            toast.value.toast(response.data.msg);
            countDown.value = 30;

          })
          .catch((error) => {
            toast.value.toast(error.response.data.err, "#FF5252");

          });

      }

    }

    const onSubmit = () => {
      if (!verify.value) return;
      loading.value = true;


      if (route.query.isSignup) {


        axios
          .post(import.meta.env.VITE_API_ENDPOINT+"/user/signup/verify_email", {
            email: route.query.email,
            code: Number(otp.value),
          })
          .then((response) => {
            toast.value.toast(response.data.msg);
            localStorage.removeItem('fullname');
            localStorage.removeItem('password');
            localStorage.removeItem('confirm_password');
            router.push({
              name: 'Login',
            });

          })
          .catch((error) => {
            toast.value.toast(error.response.data.err, "#FF5252");
            loading.value = false;

          });
      } else {
        axios
          .post(import.meta.env.VITE_API_ENDPOINT+"/user/forget_password/verify_email", {
            email: route.query.email,
            code: Number(otp.value),
          })
          .then((response) => {
            toast.value.toast(response.data.msg);

            router.push({
              name: 'NewPassword',
              query: { "email": route.query.email, }
            });

          })
          .catch((error) => {
            toast.value.toast(error.response.data.err, "#FF5252");

            loading.value = false;

          });

      }
    };



    return {
      onSubmit,
      handleOnComplete,
      handleOnChange,
      clearInput,
      resetHandler,
      otpInput,
      countDown,
      verify,
      loading,
      toast
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

/* Background colour of an input field with value */
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