import { createRouter, createWebHistory } from 'vue-router'
import MenuView from '../views/MenuView.vue'
import PlanView from '../views/PlanView.vue'
import ReportView from '../views/ReportView.vue' // 1. 追加
import OverviewView from '../views/OverviewView.vue'
import axios from 'axios'

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

// ナビゲーションガード：ロードマップ未生成時の遷移制限
router.beforeEach(async (to, from, next) => {
  const activeAnimalId = localStorage.getItem('active_animal_id')
  if (!activeAnimalId) {
    return next()
  }

  // 「1. Select Menu」以外のタブにアクセスしようとした場合
  if (to.path !== '/') {
    try {
      // プランの有無をバックエンドから確認
      const res = await axios.get('/api/plans', {
        headers: { 'X-User-Id': activeAnimalId }
      })
      const hasRoadmap = res.data && res.data.length > 0
      
      if (!hasRoadmap) {
        alert("カリキュラムが選択されていません。まずは「1. Select Menu」でカリキュラムを選択し、Roadmapを生成してください。")
        return next('/')
      }
    } catch (e) {
      console.error("Navigation guard error:", e)
    }
  }
  next()
})

