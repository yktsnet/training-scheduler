<template>
  <div class="plan-view">
    <div class="view-header">
      <h2 class="title-en">Daily Training Log 📔</h2>
      <p class="subtitle">開始日からの歩みを記録します。カードをクリックして編集、外をクリックして保存。</p>
    </div>

    <div class="list-container">
      <div 
        v-for="day in timeline" 
        :key="day.date" 
        :id="'date-' + day.date"
        class="custom-card" 
        :class="{ 'card-active': day.isEditing, 'is-today': day.isToday }"
      >
        <div class="card-header">
          <div class="header-left">
            <div class="title-group">
              <span class="label-tag" :class="{ 'tag-today': day.isToday }">
                {{ day.isToday ? 'TODAY' : 'LOG' }}
              </span>
              <h3 class="item-name">{{ formatDateFull(day.date) }}</h3>
            </div>
          </div>
          
          <div class="header-right">
            <div class="segmented-control">
              <button 
                @click.stop="day.isEditing = false" 
                :class="{ 'is-selected': !day.isEditing }"
                class="control-btn"
              >👀 View</button>
              <button 
                @click.stop="day.isEditing = true" 
                :class="{ 'is-selected': day.isEditing }"
                class="control-btn"
              >✏️ Edit</button>
            </div>
          </div>
        </div>
        
        <div class="editor-area">
          <transition name="fade" mode="out-in">
            <textarea 
              v-if="day.isEditing"
              v-model="day.content" 
              @blur="exitEditMode(day)" 
              placeholder="今日の記録を入力..."
              spellcheck="false"
              v-focus
            ></textarea>

            <div 
              v-else 
              class="rendered-view" 
              v-html="linkify(day.content)"
              @click="enterEditMode(day)"
            ></div>
          </transition>
        </div>

        <div class="card-footer">
          <div v-if="day.isEditing" class="save-status">
            <span class="dot"></span> Editing Mode (Auto-saving)
          </div>
          <div v-else class="view-status">
            Click to edit content
          </div>
        </div>
      </div>
    </div>

    <button v-if="reports.length > 3" class="jump-today-btn" @click="scrollToToday">
      Jump to Today 👇
    </button>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import axios from 'axios';

const reports = ref([]);
const todayStr = new Date().toISOString().split('T')[0];
const todayDraft = ref({ date: todayStr, content: '', isEditing: true, isToday: true });

const vFocus = { mounted: (el) => el.focus() };

const fetchReports = async () => {
  try {
    const res = await axios.get('/api/reports');
    const data = Array.isArray(res.data) ? res.data : [];
    reports.value = data.map(r => ({ ...r, isEditing: false, isToday: r.date === todayStr }));
  } catch (e) {
    console.error("Fetch error:", e);
  }
};

const timeline = computed(() => {
  const list = [...reports.value];
  if (!list.some(r => r.date === todayStr)) {
    list.unshift(todayDraft.value);
  }
  return list.sort((a, b) => b.date.localeCompare(a.date));
});

const enterEditMode = (day) => {
  day.isEditing = true;
};

const exitEditMode = async (day) => {
  // 空白でなければ保存
  if (day.content.trim()) {
    try {
      await axios.post(`/api/reports/${day.date}`, { content: day.content });
    } catch (e) {
      console.error("Save failed", e);
    }
  }
  setTimeout(() => {
    day.isEditing = false;
    // 保存したのが「今日の枠」だった場合、リストを更新して正式なレコードとして表示
    if (day.date === todayStr) fetchReports();
  }, 100);
};

const formatDateFull = (dateStr) => {
  const d = new Date(dateStr);
  return d.toLocaleDateString('ja-JP', { 
    year: 'numeric', month: '2-digit', day: '2-digit', weekday: 'short' 
  });
};

const scrollToToday = () => {
  document.getElementById('date-' + todayStr)?.scrollIntoView({ behavior: 'smooth' });
};

const linkify = (text) => {
  if (!text) return '<span class="placeholder">記録がありません。クリックして記入を開始。</span>';
  const escaped = text.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
  const withBr = escaped.replace(/\n/g, '<br>');
  const urlPattern = /(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/ig;
  return withBr.replace(urlPattern, '<a href="$1" target="_blank" rel="noopener noreferrer">$1</a>');
};

onMounted(fetchReports);
</script>

<style scoped>
/* Tab2 (PlanView) 準拠のデザイン */
.plan-view { padding-bottom: 80px; max-width: 900px; margin: 0 auto; }
.view-header { margin-bottom: 2.5rem; }
.title-en { font-size: 1.75rem; font-weight: 900; color: #0f172a; margin-bottom: 0.5rem; }
.subtitle { color: #64748b; font-size: 0.95rem; }

.title-icon {
  vertical-align: -3px;
  margin-left: 6px;
  color: var(--primary);
}

.btn-icon {
  vertical-align: -2px;
  margin-right: 4px;
}

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

/* ヘッダー左側のレイアウト：番号エリアを消してスッキリ */
.header-left { display: flex; align-items: center; }

.item-name { font-size: 1.25rem; font-weight: 800; margin: 0; color: #0f172a; }
.label-tag { font-size: 0.65rem; font-weight: 800; color: #4f46e5; letter-spacing: 0.1em; display: block; }
.tag-today { color: #10b981; }

.segmented-control { display: flex; background: #f1f5f9; padding: 4px; border-radius: 10px; border: 1px solid #e2e8f0; }
.control-btn {
  padding: 6px 16px; border-radius: 7px; font-size: 0.75rem; font-weight: 700;
  color: #64748b; background: transparent; border: none; cursor: pointer;
}
.control-btn.is-selected { background: #fff; color: #0f172a; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }

textarea, .rendered-view {
  width: 100%; min-height: 250px; padding: 2rem; border: none; outline: none;
  font-family: 'JetBrains Mono', monospace; font-size: 1rem; line-height: 1.8;
}
textarea { background-image: linear-gradient(#f1f5f9 1px, transparent 1px); background-size: 100% 1.8em; }
.rendered-view { background: #fff; cursor: pointer; word-break: break-all; }
.rendered-view :deep(a) { color: #4f46e5; font-weight: 700; text-decoration: underline; }

.card-footer { padding: 1rem 1.75rem; border-top: 1px solid #f1f5f9; background: #fff; }
.save-status { font-size: 0.75rem; font-weight: 700; color: #10b981; display: flex; align-items: center; gap: 6px; }
.view-status { font-size: 0.75rem; color: #94a3b8; font-style: italic; }
.dot { width: 8px; height: 8px; background-color: #10b981; border-radius: 50%; animation: pulse 2s infinite; }

.jump-today-btn {
  position: fixed; bottom: 30px; right: 30px;
  background: #0f172a; color: white; padding: 12px 24px;
  border-radius: 50px; border: none; font-weight: 800; cursor: pointer;
  box-shadow: 0 4px 15px rgba(0,0,0,0.2); z-index: 1000;
}

@keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.4; } }
.fade-enter-active, .fade-leave-active { transition: opacity 0.1s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.placeholder { color: #94a3b8; font-style: italic; }

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
    min-height: 200px;
    font-size: 0.95rem;
  }
  .card-footer {
    padding: 0.75rem 1rem;
  }
  .jump-today-btn {
    bottom: 20px;
    right: 20px;
    padding: 10px 18px;
    font-size: 0.85rem;
  }
}
</style>
