import { mount, createLocalVue } from '@vue/test-utils';
import Home from '@/views/Home';
import Vuetify from 'vuetify';

const localVue = createLocalVue();
localVue.use(Vuetify);

describe('Home container with test data', () => {
  const wrapper = mount(Home, { localVue });

  describe('Projects', () => {
    it('should load two test projects', async () => {
      expect(wrapper.vm.projects.length).toEqual(2);
      expect(wrapper.vm.projects[0].title).toEqual('Project Meow');
      expect(wrapper.vm.projects[1].title).toEqual('Piyo-piyo');
    });
  });
});
