import Vue from 'vue'
import Router from 'vue-router'
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
      name: 'home',
      component: Home,
      meta: {
        requiresAuth: true,
      }
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
      component: Project,
      meta: {
        requiresAuth: true,
      }
    },
  ]
});
