import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'
import SignIn from '@/views/SignIn.vue';
import Project from '@/views/Project.vue';

Vue.use(Router)

export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/signin',
      name: 'signin',
      component: SignIn,
    },
    {
      path: '/project/:projectId',
      name: 'project',
      component: Project,
    },
  ]
});
