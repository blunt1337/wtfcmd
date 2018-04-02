import Vue from 'vue'
import VueRouter from 'vue-router'
import Layout from './layout'
import PagesRoutes from './routes'

Vue.use(VueRouter)
const router = new VueRouter({
	mode: 'history',
	routes: PagesRoutes,
	base: '/',
})

new Vue({
	el: '#main',
	router,
	render: h => h(Layout)
})