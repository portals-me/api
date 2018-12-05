import Vue from 'vue'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import router from './router'
import firebase from 'firebase'
import vueConfig from 'vue-config'
import App from './App.vue'
import axios from 'axios'

Vue.config.productionTip = false;

Vue.use(Vuetify);
Vue.use(vueConfig, {
  // firebase serve
  API: 'http://localhost:5000',
  axios: axios,
});

// Initialize Firebase
var config = {
  apiKey: "AIzaSyAL-NuKxMhShZVARoxzNMvXrGN3A65OEps",
  authDomain: "portals-me.firebaseapp.com",
  databaseURL: "https://portals-me.firebaseio.com",
  projectId: "portals-me",
  storageBucket: "portals-me.appspot.com",
  messagingSenderId: "670077302427"
};
firebase.initializeApp(config);

new Vue({
  router,
  render: function (h) { return h(App) }
}).$mount('#app')
