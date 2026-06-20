import { createApp } from 'vue'
import App from './App.vue'
import { router } from './router' // さっき作ったrouterを読み込む
import './style.css'

const app = createApp(App)
app.use(router)
app.mount('#app')
