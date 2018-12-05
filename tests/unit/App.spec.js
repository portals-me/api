import { mount, createLocalVue } from '@vue/test-utils';
import App from '@/App';
import Vuetify from 'vuetify';
import VueRouter from 'vue-router';
import router from '@/router';
import adapter from 'axios-mock-adapter';
import * as axios from 'axios';
import vueConfig from 'vue-config'

const stub = new adapter(axios);
const API = 'http://localhost:5000';

stub.onGet(`${API}/users/me`).reply(200, {
  id: '1',
  user_name: 'me',
});

const localVue = createLocalVue();
localVue.use(Vuetify);
localVue.use(VueRouter);
localVue.use(vueConfig, {
  API: API,
  axios: axios,
});

describe('App view', () => {
  const wrapper = mount(App, { localVue, router });

  describe('User', () => {
    it('should load test user', async () => {
      // make sure that mounted promise resolved!
      await wrapper.vm.loadUser();
      expect(wrapper.vm.user.id).toEqual('1');
      expect(wrapper.vm.user.user_name).toEqual('me');
    });
  });
});
