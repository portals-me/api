import Vue from 'vue'
import Router from 'vue-router'

import SideBar from '@/components/SideBar.vue';
import UnsignedTopBar from '@/components/UnsignedTopBar.vue';
import TopBar from '@/components/TopBar.vue';

import Landing from '@/views/Landing.vue';
import Home from '@/views/Home.vue'
import Collection from '@/views/Collection.vue';
import SignIn from '@/views/SignIn.vue';
import SignUp from '@/views/SignUp.vue';
import User from '@/views/User.vue';

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
        default: SignUp,
        topbar: UnsignedTopBar,
      },
    },
    {
      path: '/signup/twitter-callback',
      name: 'signup-twitter-callback',
      components: {
        default: SignUp,
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
    {
      path: '/users/:userId',
      name: 'user',
      components: {
        default: User,
        topbar: TopBar,
      }
    },
  ]
});
