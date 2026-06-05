import { createRouter, createWebHistory } from 'vue-router'
import MenuView from '../views/MenuView.vue'
import PlanView from '../views/PlanView.vue'
import ReportView from '../views/ReportView.vue' // 1. 追加
import OverviewView from '../views/OverviewView.vue'

const routes = [
  { path: '/', component: MenuView },
  { path: '/plan', component: PlanView }, 
  { path: '/report', component: ReportView }, // 2. calendarからreportへ変更
  { path: '/overview', component: OverviewView },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})
