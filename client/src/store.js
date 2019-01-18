import Vue from 'vue';
import Vuex from 'vuex';
import sdk from '@/app/sdk';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    isDrawerOpened: true
  },
  actions: {
  },
  mutations: {
    toggleDrawer (state) {
      state.isDrawerOpened = !state.isDrawerOpened
    }
  },
});

