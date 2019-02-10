import Vue from 'vue';
import Vuex, { StoreOptions, MutationTree, GetterTree } from 'vuex';

Vue.use(Vuex);

interface State {
  'isDrawerOpened': boolean,
}

const state: State = {
  isDrawerOpened: true,
};

const getters: GetterTree<State, State> = {
  isDrawerOpened (state): boolean {
    return state.isDrawerOpened;
  }
};

const mutations: MutationTree<State> = {
  toggleDrawer (state) {
    state.isDrawerOpened = !state.isDrawerOpened
  }
}

export default new Vuex.Store<State>({
  state,
  getters,
  mutations,
});

