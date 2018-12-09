import { shallowMount } from '@vue/test-utils';
import Vue from 'vue';
import Vuetify from 'vuetify';
import store from '@/store';

import * as firebase from '@firebase/testing';
import * as fs from 'fs';

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
  const $route = {
    params: { projectId }
  };

  const testUser = {
    uid: 'testUser',
    display_name: 'testUserName',
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
    articles: [
      'article-1',
      'article-2',
    ]
  };
  const testComments = [
    {
      id: 'comment-1',
      owner: testUser.uid,
      message: 'へいへいほー にゃんぽよ',
    },
    {
      id: 'comment-2',
      owner: 'anonymous',
      message: 'てすぽよ',
    },
  ];
  const testArticles = [
    {
      id: 'article-1',
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
      id: 'article-2',
      type: 'share',
      entity: {
        source: 'ogp',
        ogp: {
          image: 'https://pbs.twimg.com/media/DtTpA7TU4AAQNx9.jpg:large',
          title: 'みょん on Twitter',
          description: '“様子',
          url: 'https://twitter.com/myuon_myon/status/1068735246793756673'
        }
      }
    },
  ];

  let Project;
  beforeEach(async () => {
    await firestore.collection('users').doc(testUser.uid).set(testUser);
    await firestore.collection('users').doc('anonymous').set({
      display_name: 'anonymous',
    });

    await firestore.collection('projects').doc(projectId).set(testProject);
    testComments.forEach(async (comment) => {
      await firestore.collection('projects').doc(projectId).collection('comments').doc(comment.id).set(comment);
    });
    testArticles.forEach(async (article) => {
      await firestore.collection('articles').doc(article.id).set(article);
    });

    let mockstore = firestore;
    jest.mock('@/instance/firestore', () => mockstore);
    Project = require('@/views/Project').default;
  });

  it('should load project', async () => {
    const wrapper = shallowMount(Project, { store, mocks: { $route } });
    await wrapper.vm.onMount();

    expect(wrapper.vm.project.title).toEqual(testProject.title);
    expect(wrapper.vm.project.description).toEqual(testProject.description);
  });

  it('should load comments and comment users', async () => {
    const wrapper = shallowMount(Project, { store, mocks: { $route } });
    await wrapper.vm.onMount();

    expect(wrapper.vm.project.comments.length).toBe(testComments.length);
    expect(wrapper.vm.project.comments[0].message).toBe(testComments[0].message);
    expect(wrapper.vm.project.comments[1].message).toBe(testComments[1].message);
    expect(wrapper.vm.project.comments[0].owner.display_name).toBe(testUser.display_name);
    expect(wrapper.vm.project.comments[1].owner.display_name).toBe('anonymous');
  });

  it('should load articles', async () => {
    const wrapper = shallowMount(Project, { store, mocks: { $route } });
    await wrapper.vm.onMount();

    expect(wrapper.vm.project.articles.length).toBe(testArticles.length);
    expect(wrapper.vm.project.articles[0].entity).toEqual(testArticles[0].entity);
    expect(wrapper.vm.project.articles[1].entity).toEqual(testArticles[1].entity);
  });
});
