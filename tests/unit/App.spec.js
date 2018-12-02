import { mount, createLocalVue } from '@vue/test-utils';
import App from '@/App';
import Vuetify from 'vuetify';

const localVue = createLocalVue();
localVue.use(Vuetify);

describe('App view', () => {
  const wrapper = mount(App, { localVue });

  describe('User', () => {
    it('should load test user', async () => {
      expect(wrapper.vm.user).toEqual({
        id: '1',
        user_name: 'me',
      });
    });
  });
});
