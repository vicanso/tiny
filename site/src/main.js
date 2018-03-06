import Vue from 'vue'
import Element from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'

import App from './App.vue'

Vue.use(Element)

Vue.config.productionTip = false

new Vue({
  render: h => h(App)
}).$mount('#app')
