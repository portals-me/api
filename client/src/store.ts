import Vue from 'vue';
import Vuex, { MutationTree, GetterTree, ActionTree } from 'vuex';
import * as types from '@/types';
import sdk from '@/app/sdk';

Vue.use(Vuex);

interface State {
  isDrawerOpened: boolean,
  collections: Array<types.Collection> | null,
}

const state: State = {
  isDrawerOpened: true,
  collections: null,
};

const getters: GetterTree<State, State> = {
  isDrawerOpened (state): boolean {
    return state.isDrawerOpened;
  }
};

const mutations: MutationTree<State> = {
  toggleDrawer (state) {
    state.isDrawerOpened = !state.isDrawerOpened
  },
  setCollections (state, collections: Array<types.Collection>) {
    state.collections = collections;
  },
}

const actions: ActionTree<State, State> = {
  async loadCollections ({ commit, state }, payload: { force?: boolean }) {
    if (state.collections == null || payload.force === true) {
      commit('setCollections', (await sdk.collection.list()).data);
    }
  }
};

export default new Vuex.Store<State>({
  state,
  getters,
  mutations,
  actions,
});

