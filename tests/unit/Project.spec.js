import { shallowMount, createLocalVue } from '@vue/test-utils';
import Project from '@/views/Project';
import Vuetify from 'vuetify';
import VueRouter from 'vue-router';
import * as mock from '../test-server/mock';
import router from '@/router';

const localVue = createLocalVue();
localVue.use(Vuetify);
localVue.use(VueRouter);

router.push('/projects/1');

beforeAll(() => {
  mock.server.start();
});

afterAll(() => {
  mock.server.shutdown();
});

describe('Project view', () => {
  const wrapper = shallowMount(Project, { localVue, router });

  describe('Project collections', () => {
    it('should have a project', async () => {
      await wrapper.vm.loadProject();
      expect(wrapper.vm.project.id).toEqual('1');
    });

    it('should have a collection', async () => {
      await wrapper.vm.loadProject();
      expect(wrapper.vm.project.collections.length).toBeGreaterThan(0);
      expect(wrapper.vm.project.collections[0].items.length).toBeGreaterThan(0);
    });

    it('should have comments', async () => {
      await wrapper.vm.loadProject();
      expect(wrapper.vm.project.comments.length).toBeGreaterThan(0);
      expect(wrapper.vm.project.comments[0].id).toEqual('1');
      expect(wrapper.vm.project.comments[0].owner.user_name).toEqual('me');
    });
  });
});
