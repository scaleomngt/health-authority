import Vue from 'vue';
import App from './App.vue';
import router from './router';
import { Dialog,Button,Input,Select,Option,Table,TableColumn, } from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
// element ui 表格自定义表头样式
import { handerStyleEvent} from '@/utils/TableStyleFun/index.js'
Vue.prototype.$handerStyleEvent = handerStyleEvent

Vue.use(Input);
Vue.use(Select);
Vue.use(Option);
Vue.use(Button);
Vue.use(Dialog);
Vue.use(Table);
Vue.use(TableColumn);

Vue.config.productionTip = false;

new Vue({
  router,
  render(h) { return h(App); },
}).$mount('#app');
