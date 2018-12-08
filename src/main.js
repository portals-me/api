import Vue from 'vue'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import router from './router'
import vueConfig from 'vue-config'
import App from './App.vue'
import axios from 'axios'
import firebase from 'firebase'
import app from '@/instance/firebase'
import store from '@/store'

const isDev = process.env.NODE_ENV === 'development';

Vue.use(Vuetify);
Vue.use(vueConfig, {
  // firebase serve
  API: 'http://localhost:5000',
  axios,
  firebase,
  isDev,
});

new Vue({
  router,
  store,
  render: function (h) { return h(App) }
}).$mount('#app');
