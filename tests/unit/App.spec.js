import { mount, createLocalVue } from '@vue/test-utils';
import App from '@/App';
import Vuetify from 'vuetify';
import VueRouter from 'vue-router';
import * as mock from '../test-server/mock';
import router from '@/router';

const localVue = createLocalVue();
localVue.use(Vuetify);
localVue.use(VueRouter);

beforeAll(() => {
  mock.server.start();
});

afterAll(() => {
  mock.server.shutdown();
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
