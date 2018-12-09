import { shallowMount } from '@vue/test-utils';
import Vue from 'vue';
import Vuetify from 'vuetify';
import VueRouter from 'vue-router';
import store from '@/store';

import * as firebase from '@firebase/testing';
import * as fs from 'fs';

Vue.use(Vuetify);
Vue.use(VueRouter);

const router = new VueRouter();

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

/*
describe('Home without User', () => {
  const mockstore = getFirestore(null);

  let Home;
  beforeEach(async () => {
    jest.mock('@/instance/firestore', () => mockstore);
    Home = require('@/views/Home').default;
  });

  it('should not load', async () => {
    shallowMount(Home, { router, store });
  });
});
*/

describe('Home view with Test User', () => {
  const testUser = {
    uid: 'testUser',
    displayName: 'testUserName',
  };
  const firestore = getFirestore({ uid: testUser.uid });

  beforeEach(async () => {
    store.commit('setUser', testUser);
    store.commit('setInitialized');

    await firestore.collection('projects').doc('test-1').set({
      title: 'Project Meow',
      description: 'ぞうの卵はおいしいぞう。ぞうの卵はおいしいぞう。ぞうの卵はおいしいぞう。',
      media: [
        'document',
        'picture',
        'movie',
      ],
      owner: testUser.uid,
      cover: {
        sort: 'solid',
        color: 'teal darken-2',
      }
    });
    await firestore.collection('projects').doc('test-2').set({
      title: 'Piyo-piyo',
      description: '',
      media: [
        'document',
      ],
      owner: testUser.uid,
      cover: {
        sort: 'solid',
        color: 'orange lighten-2',
      }
    });
  });

  describe('Projects', () => {
    let Home = null;

    beforeEach(async () => {
      let mockstore = firestore;
      jest.mock('@/instance/firestore', () => mockstore);
      Home = require('@/views/Home').default;
    });

    it('should load two test projects', async () => {
      const wrapper = shallowMount(Home, { router, store });
      await wrapper.vm.onMount();

      expect(wrapper.vm.projects.length).toBe(2);
      expect(wrapper.vm.projects[0].title).toBe('Project Meow');
      expect(wrapper.vm.projects[1].title).toBe('Piyo-piyo');
    });

    it('should render two projects', async () => {
      const wrapper = shallowMount(Home, { router, store });
      await wrapper.vm.onMount();

      expect(wrapper.findAll({ name: 'v-card' }).wrappers.some(wrapper => wrapper.text().includes('Project Meow'))).toBe(true);
      expect(wrapper.findAll({ name: 'v-card' }).wrappers.some(wrapper => wrapper.text().includes('Piyo-piyo'))).toBe(true);
    });

    it('should create new project', async () => {
      const wrapper = shallowMount(Home, { router, store });
      await wrapper.vm.onMount();

      const number = wrapper.vm.projects.length;
      wrapper.vm.form.title = 'New Project';
      await wrapper.vm.createProject();

      expect(wrapper.vm.projects.length).toBe(number + 1);
      expect(wrapper.vm.projects[0].title).toBe('New Project');
    });
  });
});
