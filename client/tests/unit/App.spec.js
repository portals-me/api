import { shallowMount } from '@vue/test-utils';
import Vue from 'vue';
import Vuetify, { Menu } from 'vuetify';
import router from '@/router';
import store from '@/store';
import Vuex from 'vuex';
import App from '@/App';

Vue.use(Vuex);
Vue.use(Vuetify);
Vue.use(router);

router.push('/');

describe('App view', () => {
  describe('With Test User', () => {
    const testUser = {
      uid: 'testUid',
      display_name: 'testUserName',
      photoURL: '',
    };

    it('should commit to store', async () => {
      store.commit('setUser', testUser);
      store.commit('setInitialized');

      expect(store.state.user).toBe(testUser);
      expect(store.state.initialized).toBe(true);
    });

    describe('Mount', () => {
      const wrapper = shallowMount(App, { store, router });

      it('should render signOut button', async () => {
        expect(wrapper.find({ name: 'v-toolbar' }).find({ name: 'v-list-tile' })).toBeTruthy();
      });
      
      it('should display name', async () => {
        expect(wrapper.find({ name: 'v-toolbar' }).text().includes(testUser.display_name)).toBe(true);
      });

      describe('SignOut', () => {
        beforeEach(async () => {
          await wrapper.vm.signOut();
        });

        it('should erase display name', () => {
          expect(wrapper.find({ name: 'v-toolbar' }).text().includes(testUser.display_name)).toBe(false);
        });
  
        it('should erase user state', () => {
          expect(store.state.user).toBe(null);
        });

        it('should be redirected to /signin', () => {
          expect(wrapper.vm.$route.path).toBe('/signin');
        })
      });
    });
  });
});
