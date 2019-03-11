import Vue from 'vue'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import router from './router'
import App from './App.vue'
import store from '@/store'
const vueConfig = require('vue-config');
const GAuth = require('vue-google-oauth2').default;

let keys;
let key = process.env.VUE_APP_TWITTER_KEY;
try {
  keys = require('../../token/auth.json');
  key = keys.twitter;
} catch (e) {
  console.log(e);
}
key = key.split('.');

Vue.use(Vuetify);
Vue.use(vueConfig, {
  API: process.env.API_ENDPOINT,
  isDev: process.env.NODE_ENV === 'development',
  providers: {
    auth: {
      twitter: {
        clientId: key[0],
        clientSecret: key[1],
      }
    }
  },
});

Vue.use(GAuth, {
  clientId: keys.google || process.env.VUE_APP_GOOGLE_KEY,
})

new Vue({
  router,
  store,
  render: function (h) { return h(App) }
}).$mount('#app');
