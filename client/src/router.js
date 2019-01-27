import Vue from 'vue'
import Router from 'vue-router'

import SideBar from '@/components/SideBar';
import UnsignedTopBar from '@/components/UnsignedTopBar';
import TopBar from '@/components/TopBar';

import Landing from '@/views/Landing';
import Home from './views/Home'
import Collection from '@/views/Collection';
import SignIn from '@/views/SignIn';

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
      path: '/signup',
      name: 'signup',
      components: {
        default: SignIn,
        topbar: UnsignedTopBar,
      },
    },
    {
      path: '/signup/twitter-callback',
      name: 'signup-twitter-callback',
      components: {
        default: SignIn,
        topbar: UnsignedTopBar,
      },
    },
    {
      path: '/signin',
      name: 'signin',
      components: {
        default: SignIn,
        topbar: UnsignedTopBar,
      },
    },
    {
      path: '/signin/twitter-callback',
      name: 'signin-twitter-callback',
      components: {
        default: SignIn,
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
      path: '/collections/:collectionId',
      name: 'collection',
      components: {
        default: Collection,
        sidebar: SideBar,
        topbar: TopBar,
      }
    },
  ]
});
