import Vue from 'vue';
import Vuex from 'vuex';
import sdk from '@/app/sdk';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    user: null,
    initialized: false,
  },
  actions: {
    initialize ({ state, commit }) {
      return new Promise((resolve, reject) => {
        if (!state.initialized && localStorage.getItem('id_token')) {
          sdk.user.me()
            .then(user => {
              commit('setUser', user.data);
              commit('setInitialized');

              resolve();
            })
            .catch(_ => {
              reject();
            });
        } else {
          resolve();
        }
      });
    },
    signOut ({ commit }) {
      commit('setUser', null);
    },
  },
  getters: {
    isAuthenticated (state) {
      return state.user != null;
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

