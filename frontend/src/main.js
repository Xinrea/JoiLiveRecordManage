import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify'
import VueRouter from 'vue-router'
import FileList from './components/FileList'
import axios from 'axios'

Vue.use(VueRouter)
Vue.config.productionTip = false

const router = new VueRouter({
  routes: [
    {
      path: '/',
      component: FileList,
    }
  ]
})

new Vue({
  vuetify,
  router,
  axios,
  render: h => h(App)
}).$mount('#app')
