import Vue from 'vue';
import Vuex, { MutationTree, GetterTree, ActionTree } from 'vuex';
import * as types from '@/types';
import sdk from '@/app/sdk';

Vue.use(Vuex);

interface State {
  isDrawerOpened: boolean,
  collections: Array<types.Collection> | null,
  collectionsForceReload: boolean,
}

const state: State = {
  isDrawerOpened: true,
  collections: null,
  collectionsForceReload: false,
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
  markCollectionsReload (state) {
    state.collectionsForceReload = true;
  },
  unmarkCollectionsReload (state) {
    state.collectionsForceReload = false;
  },
}

const actions: ActionTree<State, State> = {
  async loadCollections ({ commit, state }, payload: { force?: boolean }) {
    if (state.collections == null || payload.force === true || state.collectionsForceReload) {
      commit('setCollections', (await sdk.collection.list()).data);
      commit('unmarkCollectionsReload');
    }
  }
};

export default new Vuex.Store<State>({
  state,
  getters,
  mutations,
  actions,
});

