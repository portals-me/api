import Vue from 'vue'
import Router from 'vue-router'

import SideBar from '@/components/SideBar';
import UnsignedTopBar from '@/components/UnsignedTopBar';
import TopBar from '@/components/TopBar';

import Landing from '@/views/Landing';
import Home from './views/Home'
import SignIn from '@/views/SignIn';
import SignUp from '@/views/SignUp';
import Project from '@/views/Project';

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
      path: '/signin',
      name: 'signin',
      component: SignIn,
    },
    {
      path: '/signup',
      name: 'signup',
      component: SignUp,
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
  ]
});
