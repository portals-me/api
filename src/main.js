import Vue from 'vue'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import router from './router'
import vueConfig from 'vue-config'
import App from './App.vue'
import axios from 'axios'

const isDev = process.env.NODE_ENV === 'development';

(async () => {
  /*
  const firebase = await (isDev ? import('@firebase/testing') : import('firebase'));
  
  if (isDev) {
    firebase.initializeTestApp({
      projectId: `project-${Date.now()}`,
      auth: { uid: 'alice', email: 'alice@example.com' }
    });
  } else {
    firebase.initializeApp({
      apiKey: "AIzaSyAL-NuKxMhShZVARoxzNMvXrGN3A65OEps",
      authDomain: "portals-me.firebaseapp.com",
      databaseURL: "https://portals-me.firebaseio.com",
      projectId: "portals-me",
      storageBucket: "portals-me.appspot.com",
      messagingSenderId: "670077302427"
    });
  }
  */

  Vue.use(Vuetify);
  Vue.use(vueConfig, {
    // firebase serve
    API: 'http://localhost:5000',
    axios,
//    firebase,
    isDev,
  });

  new Vue({
    router,
    render: function (h) { return h(App) }
  }).$mount('#app');
})();
