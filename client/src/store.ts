import Vue from 'vue';
import Vuex, { StoreOptions, MutationTree, GetterTree, ActionTree } from 'vuex';
import * as types from '@/types';

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
  async loadCollections ({ commit, state }, loader: () => Promise<Array<types.Collection>>) {
    if (state.collections == null) {
      commit('setCollections', await loader());
    }
  }
};

export default new Vuex.Store<State>({
  state,
  getters,
  mutations,
  actions,
});

