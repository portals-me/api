import Vue from 'vue'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import router from './router'
import vueConfig from 'vue-config'
import App from './App.vue'
import axios from 'axios'
import firebase from 'firebase'
import Vuex from 'vuex'

const isDev = process.env.NODE_ENV === 'development';

firebase.initializeApp({
  apiKey: "AIzaSyAL-NuKxMhShZVARoxzNMvXrGN3A65OEps",
  authDomain: "portals-me.firebaseapp.com",
  databaseURL: "https://portals-me.firebaseio.com",
  projectId: "portals-me",
  storageBucket: "portals-me.appspot.com",
  messagingSenderId: "670077302427"
});

Vue.use(Vuetify);
Vue.use(Vuex);
Vue.use(vueConfig, {
  // firebase serve
  API: 'http://localhost:5000',
  axios,
  firebase,
  isDev,
});

const store = new Vuex.Store({
  state: {
    user: null,
    initialized: false,
  },
  actions: {
    initialize ({ state, commit }) {
      if (!state.initialized) {
        firebase.auth().onAuthStateChanged((user) => {
          if (user) {
            commit('setUser', user);
          }
          commit('setInitialized');
        });
      }
    },
    signOut ({ commit }) {
      firebase.auth().signOut();
      commit('setUser', null);
    }
  },
  mutations: {
    setUser (state, user) {
      state.user = user;
    },
    setInitialized (state) {
      state.initialized = true;
    },
  },
});

new Vue({
  router,
  store,
  render: function (h) { return h(App) }
}).$mount('#app');
