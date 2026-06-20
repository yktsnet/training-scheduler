<template>
  <div class="menu-view">
    <div class="view-header">
      <h2 class="title-en">Curriculum Selection 🔍</h2>
      <p class="subtitle">研修項目を選択してください。番号順に進めていくのがスムーズです。</p>
      
      <!-- 追加：チェックが0件のときに表示するアラートメッセージ -->
      <transition name="fade">
        <div v-if="selectedIds.length === 0" class="selection-guide-alert">
          💡 まずは研修項目を1つ以上チェックして、下部の「Generate Roadmap」ボタンを押してください。
        </div>
      </transition>
    </div>

    <div class="menu-list">
      <div 
        v-for="(menu, index) in menus" 
        :key="menu.id" 
        class="menu-card" 
        :class="{ 'is-selected': selectedIds.includes(menu.id) }"
      >
        <div class="card-header">
          <div class="header-left">
            <div class="index-number">{{ index + 1 }}</div>
            
            <div class="checkbox-container">
              <input type="checkbox" :id="'m-'+menu.id" :value="menu.id" v-model="selectedIds">
              <label :for="'m-'+menu.id" class="custom-check"></label>
            </div>
          </div>

          <div class="header-main-content">
            <h3 class="menu-name">{{ menu.name }}</h3>
            <div class="header-meta">
              <span class="days-badge">{{ menu.days }} Days</span>
              <div class="difficulty-stars">
                {{ "★".repeat(menu.difficulty) }}<span class="empty-stars">{{ "☆".repeat(5-menu.difficulty) }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="card-content">
          <div class="info-block-main">
            <h4 class="block-label">📝 SUMMARY</h4>
            <p class="block-text">{{ menu.summary }}</p>
            
            <div class="link-action">
              <a :href="menu.doc_link" target="_blank" class="reference-card-btn">
                <span class="icon">📖</span> View Reference Docs
              </a>
            </div>
          </div>
          
          <div class="info-block-sub">
            <div class="sub-item">
              <h4 class="block-label">🛠️ SKILLS</h4>
              <p class="block-text">{{ menu.skills }}</p>
            </div>
            <div class="sub-item">
              <h4 class="block-label">🎓 PREREQUISITE</h4>
              <p class="block-text">{{ menu.prerequisites }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="footer-container">
      <button 
        class="btn-generate" 
        @click="generatePlans" 
        :disabled="selectedIds.length === 0"
      >
        Generate Roadmap ✨
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useRouter } from 'vue-router';

const menus = ref([]);
const selectedIds = ref([]);
const router = useRouter();

const fetchMenus = async () => {
  try {
    const res = await axios.get('/api/menus');
    menus.value = res.data;
    
    // 既にプランがあるものを初期チェック
    const planRes = await axios.get('/api/plans');
    const existingMenuNames = planRes.data.map(p => p.menu_name);
    selectedIds.value = menus.value
      .filter(m => existingMenuNames.includes(m.name))
      .map(m => m.id);
  } catch (e) {
    console.error("Fetch error:", e);
  }
};

const generatePlans = async () => {
  try {
    await axios.post('/api/menus/select', { menu_ids: selectedIds.value });
    router.push('/plan');
  } catch (e) {
    console.error(e);
  }
};

onMounted(fetchMenus);
</script>

<style scoped>
.menu-view { padding-bottom: 140px; }
.view-header { margin-bottom: 2.5rem; }
.title-en { font-size: 1.75rem; font-weight: 900; color: #0f172a; margin-bottom: 0.5rem; }
.subtitle { color: #64748b; font-size: 0.95rem; }

.title-icon {
  vertical-align: -3px;
  margin-left: 6px;
  color: var(--primary);
}

.label-icon {
  vertical-align: -2px;
  margin-right: 6px;
  color: #64748b;
}

.btn-icon {
  vertical-align: -2px;
  margin-right: 6px;
}

.menu-list { display: flex; flex-direction: column; gap: 1.5rem; }

.menu-card {
  background: #ffffff;
  border: 2px solid #cbd5e1;
  border-radius: 16px;
  padding: 1.75rem;
  transition: all 0.2s ease;
}
.menu-card.is-selected { border-color: #4f46e5; background-color: #f8faff; box-shadow: 0 0 0 1px #4f46e5; }

.card-header {
  display: flex;
  align-items: center;
  gap: 1.25rem;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 2px solid #f1f5f9;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 15px;
}

.header-main-content {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1.5rem;
}

/* インデックス番号のスタイル */
.index-number {
  font-family: 'Inter', sans-serif;
  font-size: 1.2rem; font-weight: 900; color: #64748b;
  min-width: 20px;
}
.is-selected .index-number { color: #4f46e5; }

.menu-name {
  font-size: 1.3rem;
  font-weight: 800;
  margin: 0;
  color: #0f172a;
}

.header-meta {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.days-badge {
  background: #e2e8f0;
  color: #4f46e5;
  font-size: 0.75rem;
  font-weight: 800;
  padding: 4px 12px;
  border-radius: 6px;
}

.card-content { display: grid; grid-template-columns: 1.4fr 1fr; gap: 2.5rem; }
.block-label { font-size: 0.7rem; font-weight: 800; color: #475569; margin-bottom: 0.75rem; }
.block-text { font-size: 0.95rem; line-height: 1.6; margin: 0; color: #1e293b; }

.link-action { margin-top: 1.5rem; }
.reference-card-btn {
  display: inline-flex; align-items: center; gap: 8px;
  background: #f1f5f9; color: #4f46e5; text-decoration: none;
  padding: 10px 18px; border-radius: 10px; font-size: 0.85rem; font-weight: 700; border: 1px solid #e2e8f0;
}
.reference-card-btn:hover { background: #4f46e5; color: white; }

.info-block-sub { display: flex; flex-direction: column; gap: 1.5rem; padding-left: 1.5rem; border-left: 2px solid #f1f5f9; }

.checkbox-container { position: relative; width: 30px; height: 30px; }
.checkbox-container input { display: none; }
.custom-check { display: block; width: 30px; height: 30px; border: 2px solid #94a3b8; border-radius: 8px; background: white; cursor: pointer; }
.checkbox-container input:checked + .custom-check { background: #4f46e5; border-color: #4f46e5; }
.checkbox-container input:checked + .custom-check::after { content: "✓"; color: white; position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); font-weight: 900; }

/* フッターボタン対策 */
.footer-container {
  position: fixed; bottom: 0; left: 0; right: 0; padding: 25px;
  background: rgba(255, 255, 255, 0.9); backdrop-filter: blur(10px);
  border-top: 1px solid #e2e8f0; display: flex; justify-content: center; z-index: 1000;
}
.btn-generate {
  background: #10b981; color: white !important;
  padding: 16px 80px; font-size: 1.1rem; font-weight: 800; border-radius: 14px;
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3); border: none; cursor: pointer;
}
.btn-generate:disabled { background: #cbd5e1; color: #94a3b8 !important; cursor: not-allowed; box-shadow: none; }

/* 追加：未選択時のガイドアラートのスタイル */
.selection-guide-alert {
  margin-top: 1rem;
  padding: 12px 20px;
  background-color: #fff7ed; /* 温かみのある薄いオレンジ */
  border: 1px solid #ffedd5;
  border-radius: 10px;
  color: #ea580c;            /* 視認性の良い濃いオレンジ */
  font-size: 0.9rem;
  font-weight: 700;
  display: inline-block;
  box-shadow: 0 2px 4px rgba(234, 88, 12, 0.05);
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

/* ========== モバイル対応のレスポンシブスタイル ========== */
@media (max-width: 768px) {
  .menu-view { padding-bottom: 100px; }
  .view-header { margin-bottom: 1.5rem; }
  .title-en { font-size: 1.4rem; }
  .subtitle { font-size: 0.85rem; }
  
  .menu-card {
    padding: 1.25rem;
    border-radius: 12px;
  }
  
  /* カードヘッダーのスマホレイアウト */
  .card-header {
    display: flex;
    flex-direction: row;
    align-items: flex-start;
    gap: 0.75rem;
    margin-bottom: 1rem;
    padding-bottom: 0.75rem;
  }
  .header-left {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-top: 2px;
  }
  .header-main-content {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 0.35rem;
  }
  .menu-name {
    font-size: 1.15rem;
    line-height: 1.4;
  }
  .header-meta {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }
  .days-badge {
    align-self: auto;
    padding: 2px 8px;
    font-size: 0.7rem;
  }
  .difficulty-stars {
    font-size: 0.85rem;
  }
  
  /* カードコンテンツの1カラム化 */
  .card-content {
    grid-template-columns: 1fr;
    gap: 1.25rem;
  }
  
  .info-block-sub {
    padding-left: 0;
    border-left: none;
    border-top: 1px solid #f1f5f9;
    padding-top: 1.25rem;
    gap: 1.25rem;
  }
  
  .link-action {
    margin-top: 1rem;
  }
  
  /* フッターボタンの最適化 */
  .footer-container {
    padding: 15px;
  }
  .btn-generate {
    width: 100%;
    padding: 12px 24px;
    font-size: 1rem;
    border-radius: 10px;
  }
}
</style>
