<template>
  <div class="app-container">
    
    <!-- アニマル選択モーダル（管理者ルートでは非表示） -->
    <div v-if="showAnimalModal && !isAdminRoute" class="animal-overlay">
      <div class="animal-modal">
        <h2>Who are you? 🐾</h2>
        <p>あなたのアニマルを選んでください</p>
        
        <div class="animal-grid">
          <button v-for="u in users" :key="u.id" @click="login(u)" 
                  class="animal-btn" :class="{ active: activeAnimal?.id === u.id }">
            {{ u.emoji }}
          </button>
          
          <div class="divider"></div>
          <button v-for="emoji in availableEmojis" :key="emoji" @click="createUser(emoji)" 
                  class="animal-btn new-animal">
            {{ emoji }} +
          </button>
        </div>

        <div class="admin-link-wrapper">
          <router-link to="/admin/login" @click="showAnimalModal = false" class="admin-link">🐾 管理者としてログイン</router-link>
        </div>

        <div v-if="activeAnimal" class="danger-zone">
          <button @click="deleteUser" class="btn-delete">🗑️ このアニマルを初期化(削除)する</button>
          <button @click="showAnimalModal = false" class="btn-close">戻る</button>
        </div>
      </div>
    </div>

    <!-- 管理者用ルートの場合：ヘッダーやモーダル制御なしで直接描画 -->
    <template v-if="isAdminRoute">
      <main class="content-area">
        <router-view :key="$route.fullPath" />
      </main>
    </template>

    <!-- 一般アニマルユーザー用ルートの場合 -->
    <template v-else-if="activeAnimal && !showAnimalModal">
      <header class="main-header">
        <h1 class="brand">Training Scheduler 📅</h1>
        <div class="current-animal">Playing as: {{ activeAnimal.emoji }}</div>
        
        <nav class="tab-nav">
          <router-link to="/" class="nav-item">
            <span class="icon">🔍</span> 1. Select Menu
          </router-link>
          <router-link 
            to="/plan" 
            class="nav-item" 
            :class="{ 'nav-disabled': !hasRoadmap }"
            @click.prevent="handleTabClick('/plan', $event)"
          >
            <span class="icon">✍️</span> 2. Edit Plan
          </router-link>
          <router-link 
            to="/report" 
            class="nav-item" 
            :class="{ 'nav-disabled': !hasRoadmap }"
            @click.prevent="handleTabClick('/report', $event)"
          >
            <span class="icon">📔</span> 3. Daily Log
          </router-link>
          <router-link 
            to="/overview" 
            class="nav-item" 
            :class="{ 'nav-disabled': !hasRoadmap }"
            @click.prevent="handleTabClick('/overview', $event)"
          >
            <span class="icon">🚀</span> 4. Overview
          </router-link>
        </nav>
      </header>

      <main class="content-area">
        <router-view :key="$route.fullPath" />
      </main>

      <button class="fab-settings" @click="showAnimalModal = true" aria-label="Settings">🐾</button>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import axios from 'axios';

const showAnimalModal = ref(true);
const activeAnimal = ref(null);
const users = ref([]);
const hasRoadmap = ref(false);

const router = useRouter();
const route = useRoute();

const isAdminRoute = computed(() => route.path.startsWith('/admin'));

const allEmojis = ['🦁','🐰','🦊','🐼','🐨','🐯','🐸','🐵','🐧','🦉','🐺','🐴','🐗','🐢','🐍','🐬','🦖','🦍','🦥','🦦'];

const availableEmojis = computed(() => {
  const used = users.value.map(u => u.emoji);
  return allEmojis.filter(e => !used.includes(e));
});

const checkRoadmapStatus = async () => {
  if (!activeAnimal.value) {
    hasRoadmap.value = false;
    return;
  }
  try {
    const res = await axios.get('/api/plans');
    hasRoadmap.value = res.data && res.data.length > 0;
  } catch (e) {
    hasRoadmap.value = false;
  }
};

const handleTabClick = (path, event) => {
  if (!hasRoadmap.value) {
    return;
  }
  router.push(path);
};

const fetchUsers = async () => {
  try {
    const res = await axios.get('/api/users');
    users.value = res.data;
    
    const savedId = localStorage.getItem('active_animal_id');
    if (savedId) {
      const found = users.value.find(u => u.id === parseInt(savedId));
      if (found) await login(found);
    }
  } catch (e) {
    console.error(e);
  }
};

const login = async (user) => {
  activeAnimal.value = user;
  localStorage.setItem('active_animal_id', user.id);
  axios.defaults.headers.common['X-User-Id'] = user.id; 
  showAnimalModal.value = false;
  await checkRoadmapStatus();
};

const createUser = async (emoji) => {
  try {
    const res = await axios.post('/api/users', { emoji });
    await fetchUsers();
    login(res.data);
  } catch (e) {
    console.error(e);
  }
};

const deleteUser = async () => {
  if(confirm(`【警告】${activeAnimal.value.emoji} の全データ（計画・日報・進捗）を削除します。元には戻せません。本当によろしいですか？`)) {
    try {
      await axios.delete(`/api/users/${activeAnimal.value.id}`);
      localStorage.removeItem('active_animal_id');
      activeAnimal.value = null;
      delete axios.defaults.headers.common['X-User-Id'];
      await fetchUsers();
    } catch (e) {
      console.error(e);
    }
  }
};

onMounted(fetchUsers);

watch(() => activeAnimal.value, checkRoadmapStatus);
watch(() => route.path, checkRoadmapStatus);
</script>

<style scoped>
/* ========== 元のスタイル ========== */
.app-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.main-header {
  margin-bottom: 3rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.icon {
  margin-right: 6px;
}

.brand {
  font-size: 1.8rem;
  font-weight: 800;
  color: var(--text-main, #0f172a);
  letter-spacing: -0.025em;
  margin-bottom: 0.5rem; /* 少し詰めました */
}

.tab-nav {
  display: inline-flex;
  background: #e2e8f0;
  padding: 4px;
  border-radius: 12px;
  gap: 4px;
}

.nav-item {
  text-decoration: none;
  color: var(--text-muted, #64748b);
  padding: 10px 24px;
  border-radius: 8px;
  font-size: 0.95rem;
  font-weight: 600;
  transition: all 0.2s;
}

.nav-item:hover {
  color: var(--text-main, #0f172a);
}

.router-link-active {
  background: white;
  color: var(--primary, #4f46e5) !important;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
}

.nav-disabled {
  opacity: 0.5;
  cursor: not-allowed !important;
  color: #94a3b8 !important;
  background: transparent !important;
  box-shadow: none !important;
}

.nav-disabled:hover {
  color: #94a3b8 !important;
}

.content-area {
  animation: fadeIn 0.4s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.current-animal {
  font-size: 1rem;
  font-weight: 700;
  color: #64748b;
  margin-bottom: 1.5rem;
}

.fab-settings {
  position: fixed;
  top: 30px;
  right: 30px;
  width: 54px;
  height: 54px;
  border-radius: 50%;
  background: #fff;
  border: 2px solid #cbd5e1;
  font-size: 1.6rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  cursor: pointer;
  transition: all 0.2s;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
}

.fab-settings:hover {
  transform: scale(1.1) rotate(15deg);
  box-shadow: 0 6px 16px rgba(0,0,0,0.15);
}

.animal-overlay {
  position: fixed; top: 0; left: 0; width: 100vw; height: 100vh;
  background: rgba(15, 23, 42, 0.8); backdrop-filter: blur(8px);
  display: flex; justify-content: center; align-items: flex-start; z-index: 2000;
  overflow-y: auto;
  padding: 2rem 0;
}

.animal-modal {
  background: #fff; padding: 3rem; border-radius: 20px;
  max-width: 600px; width: 90%; text-align: center;
  box-shadow: 0 20px 40px rgba(0,0,0,0.2);
  margin: auto;
}

.animal-grid {
  display: flex; flex-wrap: wrap; gap: 10px; justify-content: center; margin: 2rem 0;
}

.animal-btn {
  font-size: 2rem; padding: 10px; border: 2px solid transparent;
  background: #f1f5f9; border-radius: 12px; cursor: pointer; transition: 0.2s;
}

.animal-btn:hover { background: #e2e8f0; transform: translateY(-2px); }
.animal-btn.active { border-color: #4f46e5; background: #e0e7ff; }
.new-animal { font-size: 1.5rem; opacity: 0.6; }
.new-animal:hover { opacity: 1; }

.divider { width: 100%; height: 1px; margin: 10px 0; }

.danger-zone {
  display: flex; justify-content: space-between; align-items: center;
  margin-top: 2rem; border-top: 1px solid #e2e8f0; padding-top: 1.5rem;
}

.btn-delete { color: #ef4444; background: none; border: none; cursor: pointer; font-weight: 700; font-size: 0.9rem; }
.btn-delete:hover { text-decoration: underline; }
.btn-close { background: #cbd5e1; padding: 10px 20px; border-radius: 8px; border: none; cursor: pointer; font-weight: 700; }
.btn-close:hover { background: #94a3b8; color: white; }

.admin-link-wrapper {
  margin-top: 1.5rem;
  border-top: 1px dashed #e2e8f0;
  padding-top: 1rem;
}
.admin-link {
  color: var(--text-muted, #64748b);
  font-size: 0.85rem;
  font-weight: 700;
  text-decoration: none;
  transition: color 0.2s;
}
.admin-link:hover {
  color: var(--primary, #4f46e5);
}

/* ========== モバイル対応のレスポンシブスタイル ========== */
@media (max-width: 640px) {
  .app-container {
    padding: 1rem 0.75rem;
  }
  .main-header {
    margin-bottom: 1.5rem;
  }
  .brand {
    font-size: 1.5rem;
  }
  .current-animal {
    margin-bottom: 1rem;
  }
  .tab-nav {
    display: flex;
    width: 100%;
    overflow-x: auto;
    white-space: nowrap;
    border-radius: 8px;
    padding: 2px;
    gap: 2px;
    scrollbar-width: none; /* Firefox */
  }
  .tab-nav::-webkit-scrollbar {
    display: none; /* Safari/Chrome */
  }
  .nav-item {
    padding: 8px 12px;
    font-size: 0.8rem;
    flex: 1;
    text-align: center;
    min-width: max-content;
    border-radius: 6px;
  }
  .fab-settings {
    top: 15px;
    right: 15px;
    width: 44px;
    height: 44px;
    font-size: 1.3rem;
  }
  .animal-modal {
    padding: 1.5rem;
    border-radius: 16px;
  }
  .animal-grid {
    margin: 1rem 0;
    gap: 8px;
  }
  .animal-btn {
    font-size: 1.6rem;
    padding: 8px;
  }
}
</style>
