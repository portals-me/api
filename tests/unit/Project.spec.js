import { shallowMount } from '@vue/test-utils';
import Vue from 'vue';
import Vuetify from 'vuetify';
import VueRouter from 'vue-router';
import store from '@/store';

import * as firebase from '@firebase/testing';
import * as fs from 'fs';

Vue.use(Vuetify);

const router = new VueRouter();

/*
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
*/

Vue.use(Vuetify);

let testNumber = 0;
const projectIdBase = `firestore-emulator-${Date.now()}`;
const getProjectId = () => {
  return `${projectIdBase}-${testNumber}`;
};

const getFirestore = (auth) => {
  return firebase.initializeTestApp({
    projectId: getProjectId(),
    auth: auth,
  }).firestore();
};

beforeEach(async () => {
  testNumber ++;
  await firebase.loadFirestoreRules({
    projectId: getProjectId(),
    rules: fs.readFileSync('firestore.rules', 'utf8'),
  });
});

afterAll(async () => {
  await Promise.all(firebase.apps().map(app => app.delete()));
});

describe('Project', () => {
  const projectId = '1';
  const testUser = {
    uid: 'testUser',
    displayName: 'testUserName',
  };
  const firestore = getFirestore({ uid: testUser.uid });

  const testProject = {
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
  }

  let Project;
  beforeEach(async () => {
    await firestore.collection('projects').doc(projectId).set(testProject);

    let mockstore = firestore;
    jest.mock('@/instance/firestore', () => mockstore);
    Project = require('@/views/Project').default;
  });

  it('should render title', async () => {
    const wrapper = shallowMount(Project, {
      store,
      mocks: {
        $route: {
          params: { projectId }
        }
      }
    });
    await wrapper.vm.onMount();

    expect(wrapper.findAll('h2').length).toBe(1);
    expect(wrapper.find('h2').text().includes(testProject.title)).toBe(true);
  });
});
