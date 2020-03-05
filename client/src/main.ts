import Vue from "vue";
import App from "./App.vue";
import "./registerServiceWorker";
import router from "./router";
import store from "./store";
import axios from "axios";
import VueAxios from "vue-axios";

import "@/style.scss";

Vue.config.productionTip = false;

axios.defaults.headers.post["Access-Control-Allow-Origin"] = "*";

Vue.use(VueAxios, axios);

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
