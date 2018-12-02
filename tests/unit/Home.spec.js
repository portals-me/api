import { mount, createLocalVue } from '@vue/test-utils';
import Home from '@/views/Home';
import Vuetify from 'vuetify';
import * as mock from '../test-server/mock';

const localVue = createLocalVue();
localVue.use(Vuetify);

beforeAll(() => {
  mock.server.start();
});

afterAll(() => {
  mock.server.shutdown();
});

describe('Home view', () => {
  const wrapper = mount(Home, { localVue });

  describe('Projects', () => {
    it('should load two test projects', async () => {
      await wrapper.vm.loadProjects();
      expect(wrapper.vm.projects.length).toEqual(2);
      expect(wrapper.vm.projects[0].title).toEqual('Project Meow');
      expect(wrapper.vm.projects[1].title).toEqual('Piyo-piyo');
    });
  });
});
