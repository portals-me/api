import Vue from 'vue'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import router from './router'
import vueConfig from 'vue-config'
import App from './App.vue'
import axios from 'axios'
import firebase from 'firebase'
import Vuex from 'vuex'
import app from '@/instance/firebase'

const isDev = process.env.NODE_ENV === 'development';

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
      return new Promise((resolve, reject) => {
        if (!state.initialized) {
          firebase.auth().onAuthStateChanged((user) => {
            if (user) {
              commit('setUser', user);
            }
            commit('setInitialized');

            resolve();
          });
        } else {
          resolve();
        }
      });
    },
    signOut ({ commit }) {
      firebase.auth().signOut();
      commit('setUser', null);
    },
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
