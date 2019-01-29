import Vue from 'vue'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import router from './router'
import vueConfig from 'vue-config'
import App from './App.vue'
import axios from 'axios'
import store from '@/store'
import GAuth from 'vue-google-oauth2'
import VueAxios from 'vue-axios'
import VueAuthenticate from 'vue-authenticate'

const isDev = process.env.NODE_ENV === 'development';

Vue.use(Vuetify);
Vue.use(vueConfig, {
  API: 'https://v6bnqbi2hf.execute-api.ap-northeast-1.amazonaws.com/prod',
  isDev,
  providers: require('../../token/auth.json'),
});

Vue.use(GAuth, {
  clientId: '670077302427-0r21asrffhmuhkvfq10qa8kj86cslojn.apps.googleusercontent.com',
  scope: 'profile email https://www.googleapis.com/auth/plus.login'
});

Vue.use(VueAxios, axios);
Vue.use(VueAuthenticate, {
  baseUrl: 'http://localhost:8080',
  providers: require('../../token/auth.json'),
});

new Vue({
  router,
  store,
  render: function (h) { return h(App) }
}).$mount('#app');
