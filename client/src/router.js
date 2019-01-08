import Vue from 'vue'
import Router from 'vue-router'

import SideBar from '@/components/SideBar';
import UnsignedTopBar from '@/components/UnsignedTopBar';
import TopBar from '@/components/TopBar';

import Landing from '@/views/Landing';
import Home from './views/Home'
import Project from '@/views/Project';
import Collection from '@/views/Collection';

Vue.use(Router)

export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'landing',
      components: {
        default: Landing,
        topbar: UnsignedTopBar,
      },
    },
    {
      path: '/dashboard',
      name: 'home',
      components: {
        default: Home,
        sidebar: SideBar,
        topbar: TopBar,
      },
    },
    {
      path: '/projects/:projectId',
      name: 'project',
      components: {
        default: Project,
        sidebar: SideBar,
        topbar: TopBar,
      },
    },
    {
      path: '/collections/:collectionId',
      name: 'collection',
      components: {
        default: Collection,
        sidebar: SideBar,
        topbar: TopBar,
      }
    }
  ]
});
