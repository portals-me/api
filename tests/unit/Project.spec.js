import { shallowMount } from '@vue/test-utils';
import Project from '@/views/Project';
import Vue from 'vue';
import Vuetify from 'vuetify';
import router from '@/router';
import adapter from 'axios-mock-adapter';
import * as axios from 'axios';
import vueConfig from 'vue-config'

const stub = new adapter(axios);
const API = 'http://localhost:5000';

stub.onGet(`${API}/projects/1`).reply(200, {
  id: '1',
  title: 'Project Meow',
  description: 'ぞうの卵はおいしいぞう。ぞうの卵はおいしいぞう。ぞうの卵はおいしいぞう。',
  owner: '1',
  media: [
    'document',
    'picture',
    'movie',
  ],
  cover: {
    sort: 'solid',
    color: 'teal darken-2',
  },
  collections: [
    {
      id: '1',
      name: 'My Collection',
      items: [
        {
          id: '1',
          type: 'share',
          entity: {
            source: 'ogp',
            ogp: {
              image: 'https://avatars2.githubusercontent.com/u/1870091?s=400&amp;v=4',
              title: 'myuon/pyxis',
              description: 'Pyxis. Contribute to myuon/pyxis development by creating an account on GitHub.',
              url: 'https://github.com/myuon/pyxis',
            }
          }
        },
        {
          id: '2',
          type: 'share',
          entity: {
            source: 'twitter',
            ogp: {
              image: 'https://pbs.twimg.com/media/DtTpA7TU4AAQNx9.jpg:large',
              title: 'みょん on Twitter',
              description: '“様子',
              url: 'https://twitter.com/myuon_myon/status/1068735246793756673'
            }
          }
        },
      ]
    }
  ],
  comments: [
    {
      id: '1',
      owner: {
        id: '1',
        user_name: 'me',
      },
      created_at: '2018/12/01 20:30:45',
      message: 'へいへいほー\nYes\nにゃんぽよ',
    },
    {
      id: '2',
      owner: {
        id: '2',
        user_name: 'he',
      },
      created_at: '2018/12/01 20:40:45',
      message: 'てすぽよ',
    },
    {
      id: '3',
      owner: {
        id: '3',
        user_name: 'she',
      },
      created_at: '2018/12/01 20:45:45',
      message: 'ぽよい',
    },
    {
      id: '4',
      owner: {
        id: '1',
        user_name: 'me',
      },
      created_at: '2018/12/01 21:50:55',
      message: 'へい',
    },
  ]
});

Vue.use(Vuetify);
Vue.use(router);
Vue.use(vueConfig, {
  API: API,
  axios: axios,
});

router.push('/projects/1');

describe('Project view', () => {
  const wrapper = shallowMount(Project, { router });

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
