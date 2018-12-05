import { mount, createLocalVue } from '@vue/test-utils';
import Home from '@/views/Home';
import Vuetify from 'vuetify';
import adapter from 'axios-mock-adapter';
import * as axios from 'axios';
import vueConfig from 'vue-config'

const stub = new adapter(axios);
const API = 'http://localhost:5000';

stub.onGet(`${API}/projects`).reply(200, [
  {
    id: '1',
    title: 'Project Meow',
    description: 'ぞうの卵はおいしいぞう。ぞうの卵はおいしいぞう。ぞうの卵はおいしいぞう。',
    media: [
      'document',
      'picture',
      'movie',
    ],
    cover: {
      sort: 'solid',
      color: 'teal darken-2',
    }
  },
  {
    id: '2',
    title: 'Piyo-piyo',
    description: '',
    media: [
      'document',
    ],
    cover: {
      sort: 'solid',
      color: 'orange lighten-2',
    }
  },
]);

const localVue = createLocalVue();
localVue.use(Vuetify);
localVue.use(vueConfig, {
  API: API,
  axios: axios,
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
