<template>
  <div class="plan-view">
    <div class="view-header">
      <h2 class="title-en">Roadmap Details ✍️</h2>
      <p class="subtitle">カードをクリックして編集、外をクリックして保存。番号はメニューの共通IDと連動しています。</p>
    </div>

    <div v-if="isLoading" class="status-msg">Loading plans...</div>
    
    <div v-else class="list-container">
      <div 
        v-for="plan in plans" 
        :key="plan.id" 
        class="custom-card" 
        :class="{ 'card-active': plan.isEditing }"
      >
        <div class="card-header">
          <div class="header-left">
            <span class="index-text">{{ plan.displayIndex }}</span>
            <div class="title-group">
              <span class="label-tag">CURRICULUM</span>
              <h3 class="item-name">{{ plan.menu_name }}</h3>
            </div>
          </div>
          
          <div class="header-right">
            <div class="segmented-control">
              <button 
                @click.stop="plan.isEditing = false" 
                :class="{ 'is-selected': !plan.isEditing }"
                class="control-btn"
              >👁️ View</button>
              <button 
                @click.stop="plan.isEditing = true" 
                :class="{ 'is-selected': plan.isEditing }"
                class="control-btn"
              >✏️ Edit</button>
            </div>
          </div>
        </div>
        
        <div class="editor-area">
          <transition name="fade" mode="out-in">
            <textarea 
              v-if="plan.isEditing"
              v-model="plan.content" 
              @blur="exitEditMode(plan)"
              placeholder="計画を入力..."
              spellcheck="false"
              v-focus
            ></textarea>

            <div 
              v-else 
              class="rendered-view" 
              v-html="linkify(plan.content)"
              @click="enterEditMode(plan)"
            ></div>
          </transition>
        </div>

        <div class="card-footer">
          <div v-if="plan.isEditing" class="save-status">
            <span class="dot"></span> Editing Mode (Auto-saving)
          </div>
          <div v-else class="view-status">
            Click to edit content
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';

const plans = ref([]);
const isLoading = ref(true);

// カスタムディレクティブ: 自動フォーカス
const vFocus = {
  mounted: (el) => el.focus()
};

const fetchPlans = async () => {
  try {
    // マスターメニューと現在の計画を同時に取得
    const [menusRes, plansRes] = await Promise.all([
      axios.get('/api/menus'),
      axios.get('/api/plans')
    ]);

    const allMenus = menusRes.data;

    // 各プランに対して、マスターデータ上のインデックスを紐付ける
    plans.value = plansRes.data.map(p => {
      // menu_name をキーにして、マスターリストの何番目かを探す
      const masterIndex = allMenus.findIndex(m => m.name === p.menu_name);
      return { 
        ...p, 
        isEditing: false,
        displayIndex: masterIndex !== -1 ? masterIndex + 1 : '?' 
      };
    });

    // 表示順もマスターデータのインデックス順にソートする（整合性のため）
    plans.value.sort((a, b) => a.displayIndex - b.displayIndex);

  } catch (e) {
    console.error("Data fetch error", e);
  } finally {
    isLoading.value = false;
  }
};

const enterEditMode = (plan) => {
  plan.isEditing = true;
};

const exitEditMode = async (plan) => {
  try {
    await axios.post(`/api/plans/${plan.id}`, { content: plan.content });
  } catch (e) {
    console.error("Save failed", e);
  }
  setTimeout(() => {
    plan.isEditing = false;
  }, 100);
};

const linkify = (text) => {
  if (!text) return '<span class="placeholder">内容が空です。クリックして入力を開始。</span>';
  const escaped = text.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
  const withBr = escaped.replace(/\n/g, '<br>');
  const urlPattern = /(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/ig;
  return withBr.replace(urlPattern, '<a href="$1" target="_blank" rel="noopener noreferrer">$1</a>');
};

onMounted(fetchPlans);
</script>

<style scoped>
/* 前回確定した高品質デザイン（style.cssに依存しない独立したルール）を維持 */
.plan-view { padding-bottom: 80px; max-width: 900px; margin: 0 auto; }
.view-header { margin-bottom: 2.5rem; }
.title-en { font-size: 1.75rem; font-weight: 900; color: #0f172a; margin-bottom: 0.5rem; }
.subtitle { color: #64748b; font-size: 0.95rem; }

.list-container { display: flex; flex-direction: column; gap: 3rem; }

.custom-card {
  background: #fff; border: 2px solid #cbd5e1; border-radius: 16px;
  overflow: hidden; transition: all 0.3s ease;
}
.custom-card.card-active {
  border-color: #4f46e5;
  box-shadow: 0 10px 25px -5px rgba(79, 70, 229, 0.1);
}

.card-header {
  background: #f8fafc; padding: 1.25rem 1.75rem; border-bottom: 2px solid #e2e8f0;
  display: flex; justify-content: space-between; align-items: center;
}

.header-left { display: flex; align-items: center; gap: 1.25rem; }

/* インデックス番号のデザイン：Tab1と完全に同期 */
.index-text { font-size: 1.2rem; font-weight: 900; color: #94a3b8; min-width: 1.5rem; }
.card-active .index-text { color: #4f46e5; }

.item-name { font-size: 1.25rem; font-weight: 800; margin: 0; color: #0f172a; }
.label-tag { font-size: 0.65rem; font-weight: 800; color: #4f46e5; letter-spacing: 0.1em; display: block; }

.segmented-control { display: flex; background: #f1f5f9; padding: 4px; border-radius: 10px; border: 1px solid #e2e8f0; }
.control-btn {
  padding: 6px 16px; border-radius: 7px; font-size: 0.75rem; font-weight: 700;
  color: #64748b; background: transparent; border: none; cursor: pointer;
}
.control-btn.is-selected { background: #fff; color: #0f172a; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }

textarea, .rendered-view {
  width: 100%; min-height: 350px; padding: 2rem; border: none; outline: none;
  font-family: 'JetBrains Mono', monospace; font-size: 1rem; line-height: 1.8;
}
textarea { background-image: linear-gradient(#f1f5f9 1px, transparent 1px); background-size: 100% 1.8em; }
.rendered-view { background: #fff; cursor: pointer; word-break: break-all; }
.rendered-view :deep(a) { color: #4f46e5; font-weight: 700; text-decoration: underline; }

.card-footer { padding: 1rem 1.75rem; border-top: 1px solid #f1f5f9; background: #fff; }
.save-status { font-size: 0.75rem; font-weight: 700; color: #10b981; display: flex; align-items: center; gap: 6px; }
.view-status { font-size: 0.75rem; color: #94a3b8; font-style: italic; }
.dot { width: 8px; height: 8px; background-color: #10b981; border-radius: 50%; animation: pulse 2s infinite; }

@keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.4; } }
.fade-enter-active, .fade-leave-active { transition: opacity 0.1s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

/* ========== モバイル対応のレスポンシブスタイル ========== */
@media (max-width: 768px) {
  .plan-view { padding-bottom: 40px; }
  .view-header { margin-bottom: 1.5rem; }
  .title-en { font-size: 1.4rem; }
  .subtitle { font-size: 0.85rem; }
  
  .list-container { gap: 1.5rem; }
  
  .card-header {
    padding: 1rem;
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  .header-left {
    width: 100%;
    align-items: flex-start;
  }
  .header-right {
    width: 100%;
  }
  .segmented-control {
    width: 100%;
    display: flex;
  }
  .control-btn {
    flex: 1;
    text-align: center;
    padding: 8px 16px;
    font-size: 0.8rem;
  }
  .item-name {
    font-size: 1.1rem;
  }
  textarea, .rendered-view {
    padding: 1.25rem;
    min-height: 250px;
    font-size: 0.95rem;
  }
  .card-footer {
    padding: 0.75rem 1rem;
  }
}
</style>
