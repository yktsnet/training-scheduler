import { createRouter, createWebHistory } from 'vue-router'
import MenuView from '../views/MenuView.vue'
import PlanView from '../views/PlanView.vue'
import ReportView from '../views/ReportView.vue'
import OverviewView from '../views/OverviewView.vue'
import AdminLoginView from '../views/AdminLoginView.vue'
import AdminMenuView from '../views/AdminMenuView.vue'
import axios from 'axios'

const routes = [
  { path: '/', component: MenuView },
  { path: '/plan', component: PlanView }, 
  { path: '/report', component: ReportView }, 
  { path: '/overview', component: OverviewView },
  { path: '/admin/login', component: AdminLoginView },
  { path: '/admin/menus', component: AdminMenuView },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})

// ナビゲーションガード：ロードマップ未生成時の遷移制限
router.beforeEach(async (to, from, next) => {
  // 管理者ルートの場合はガードを通さない
  if (to.path.startsWith('/admin')) {
    return next()
  }

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
        return next('/')
      }
    } catch (e) {
      console.error("Navigation guard error:", e)
    }
  }
  next()
})

