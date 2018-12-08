import Vue from 'vue';
import Vuex from 'vuex';
import firebase from 'firebase';

Vue.use(Vuex);

export default new Vuex.Store({
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

