import Vue from 'vue'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import router from './router'
import vueConfig from 'vue-config'
import App from './App.vue'
import store from '@/store'

const isDev = process.env.NODE_ENV === 'development';

Vue.use(Vuetify);
Vue.use(vueConfig, {
  API: 'https://v6bnqbi2hf.execute-api.ap-northeast-1.amazonaws.com/prod',
  isDev,
  providers: require('../../token/auth.json'),
});

new Vue({
  router,
  store,
  render: function (h) { return h(App) }
}).$mount('#app');
