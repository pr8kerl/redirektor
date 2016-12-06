import Vue from 'vue'
import Materials from 'vue-materials'
import App from './App.vue'
import VueResource from 'vue-resource'

Vue.use(VueResource)
Vue.use(Materials)

var vm = new Vue({
  el: '#app',
  render: h => h(App)
})
