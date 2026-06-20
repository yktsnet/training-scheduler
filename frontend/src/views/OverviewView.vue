<template>
  <div class="overview-view">
    <div class="view-header">
      <h2 class="title-en">Training Dashboard 🚀</h2>
      <p class="subtitle">現在の進捗とコンディションを客観的に把握・記録します。</p>
    </div>

    <div v-if="isLoading" class="status-msg">Loading...</div>
    
    <div v-else class="list-container">
      <div 
        v-for="item in items" 
        :key="item.id" 
        class="dash-card" 
        :class="{ 'is-completed': item.is_completed }"
      >
        <div class="dash-header">
          <div class="header-left">
            <span class="item-index">{{ item.menu_id }}</span>
            <div class="title-group">
              <h3 class="item-name">{{ item.menu_name }}</h3>
              <span class="status-badge" :class="item.is_completed ? 'status-done' : 'status-running'">
                {{ item.is_completed ? 'COMPLETED' : 'RUNNING' }}
              </span>
            </div>
          </div>
          <div class="header-right">
            <button @click="toggleComplete(item)" class="toggle-btn" :title="item.is_completed ? '進行中に戻す' : '完了にする'">
              {{ item.is_completed ? '✅' : '🏃' }}
            </button>
          </div>
        </div>

        <div class="dash-body">
          <div class="progress-section">
            <div class="day-stats">
              <span class="day-count">Day {{ calculateBusinessDay(item.start_date) }}</span>
              <span class="day-total">/ {{ item.target_days }} days</span>
            </div>
            <div class="bar-container">
              <div class="bar-fill" :style="{ width: calculateProgress(item) + '%' }"></div>
            </div>
          </div>

          <div class="settings-area">
            <div class="input-group">
              <label>📅 開始日</label>
              <input type="date" v-model="item.start_date" @change="save(item)">
            </div>
            <div class="input-group">
              <label>⏱️ 目標日数</label>
              <input type="number" v-model="item.target_days" @blur="save(item)">
            </div>
          </div>

          <div class="subjective-section">
            <div class="section-label">
              <span>進捗の手応え</span>
              <span class="condition-tag" :class="'off-' + Math.round(item.offset_days || 3)">
                {{ getOffsetLabel(item.offset_days) }}
              </span>
            </div>
            <input 
              type="range" min="1" max="5" step="1" 
              v-model.number="item.offset_days" 
              @change="save(item)"
              class="flat-slider"
            >
            <div class="slider-ticks">
              <span>苦戦</span><span>遅れ</span><span>順調</span><span>快調</span><span>爆速</span>
            </div>
          </div>

          <div class="memo-section">
            <label>📝 振り返り・メモ</label>
            <textarea 
              v-model="item.status_memo" 
              placeholder="今の状況を記録（クリックで入力）" 
              @focus="handleMemoFocus(item)"
              @blur="save(item)"
              class="flat-textarea"
              rows="2"
            ></textarea>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { eachDayOfInterval, isWeekend, parseISO, isAfter, startOfDay } from 'date-fns';

const items = ref([]);
const isLoading = ref(true);

const offsetLabels = {
  1: "🚨 かなり遅延", 2: "⚠️ 少し遅延", 3: "✅ 予定通り", 4: "✨ 前倒し", 5: "⚡ 爆速進行"
};

const getOffsetLabel = (val) => offsetLabels[Math.round(val)] || "-";

const handleMemoFocus = (item) => {
  if (item.status_memo === '未着手') item.status_memo = '';
};

const fetchOverviews = async () => {
  try {
    const res = await axios.get('/api/overview');
    items.value = res.data.map(i => {
      if (i.offset_days === 0) i.offset_days = 3;
      if (i.status_memo === '未着手') i.status_memo = '';
      return i;
    });
  } catch (e) { console.error(e); }
  finally { isLoading.value = false; }
};

const calculateBusinessDay = (startDateStr) => {
  if (!startDateStr) return '-';
  try {
    const start = parseISO(startDateStr);
    const today = startOfDay(new Date());
    if (isAfter(start, today)) return 0;
    const days = eachDayOfInterval({ start, end: today });
    return days.filter(day => !isWeekend(day)).length;
  } catch (e) { return '-'; }
};

const calculateProgress = (item) => {
  const current = calculateBusinessDay(item.start_date);
  if (current === '-' || !item.target_days) return 0;
  return Math.min(Math.max((current / item.target_days) * 100, 0), 100);
};

const save = async (item) => {
  try {
    await axios.post(`/api/overview/${item.id}`, item);
  } catch (e) { console.error(e); }
};

const toggleComplete = (item) => {
  item.is_completed = !item.is_completed;
  save(item);
};

onMounted(fetchOverviews);
</script>

<style scoped>
.overview-view { padding-bottom: 60px; max-width: 1100px; margin: 0 auto; }
.view-header { margin-bottom: 2.5rem; text-align: left; }
.title-en { font-size: 1.75rem; font-weight: 800; color: #0f172a; margin-bottom: 0.5rem; }
.subtitle { color: #64748b; font-size: 0.95rem; }

/* 2列グリッドの指定 */
.list-container { 
  display: grid; 
  grid-template-columns: repeat(2, 1fr); 
  gap: 2rem; 
  align-items: start;
}

/* カードのデザイン */
.dash-card {
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  transition: border-color 0.2s;
  overflow: hidden;
  text-align: left;
}

/* ヘッダー */
.dash-header {
  padding: 1.25rem 1.5rem;
  background: #fff;
  border-bottom: 1px solid #f1f5f9;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left { display: flex; align-items: center; gap: 1rem; }
.item-index { font-size: 1rem; font-weight: 700; color: #94a3b8; }
.title-group { display: flex; flex-direction: column; gap: 2px; }
.item-name { font-size: 1.1rem; font-weight: 800; color: #0f172a; margin: 0; }

.status-badge {
  font-size: 0.6rem; font-weight: 800; padding: 1px 6px; border-radius: 4px; width: fit-content;
}
.status-running { background: #e0e7ff; color: #4338ca; }
.status-done { background: #d1fae5; color: #065f46; }

.toggle-btn {
  background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 8px;
  width: 36px; height: 36px; display: flex; align-items: center; justify-content: center;
  cursor: pointer; transition: 0.2s;
}
.toggle-btn:hover { background: #f1f5f9; border-color: #cbd5e1; }

/* ボディ */
.dash-body { padding: 1.5rem; display: flex; flex-direction: column; gap: 1.5rem; }

/* 経過日数 */
.day-stats { margin-bottom: 0.5rem; display: flex; align-items: baseline; gap: 6px; }
.day-count { font-size: 2.2rem; font-weight: 800; color: #0f172a; line-height: 1; }
.day-total { font-size: 0.85rem; color: #94a3b8; font-weight: 600; }

.bar-container { height: 6px; background: #f1f5f9; border-radius: 10px; overflow: hidden; }
.bar-fill { height: 100%; background: #4f46e5; transition: width 0.6s ease; }

/* 設定項目 */
.settings-area { display: grid; grid-template-columns: 1.2fr 1fr; gap: 1rem; }
.input-group label { font-size: 0.7rem; font-weight: 700; color: #94a3b8; margin-bottom: 0.4rem; display: block; }
.input-group input {
  width: 100%; padding: 0.6rem; border: 1px solid #e2e8f0; border-radius: 8px;
  font-size: 0.85rem; color: #334155; outline: none; background: #fcfcfd;
}

/* 主観評価 */
.section-label { 
  display: flex; justify-content: space-between; align-items: center;
  font-size: 0.7rem; font-weight: 700; color: #64748b; margin-bottom: 0.75rem;
}
.condition-tag { font-weight: 800; color: #4f46e5; font-size: 0.75rem; }

.flat-slider { width: 100%; height: 4px; border-radius: 10px; appearance: none; background: #e2e8f0; outline: none; }
.flat-slider::-webkit-slider-thumb { 
  appearance: none; width: 16px; height: 16px; background: #fff; border: 2px solid #4f46e5; border-radius: 50%; cursor: pointer;
}
.slider-ticks { display: flex; justify-content: space-between; font-size: 0.6rem; color: #cbd5e1; margin-top: 0.4rem; font-weight: 600; }

/* メモ欄 */
.memo-section label { font-size: 0.7rem; font-weight: 700; color: #94a3b8; margin-bottom: 0.5rem; display: block; }
.flat-textarea {
  width: 100%; border: 1px solid #f1f5f9; border-radius: 8px; background: #fcfcfd;
  padding: 10px; font-size: 0.85rem; line-height: 1.5; color: #334155; resize: none; outline: none;
}
.flat-textarea:focus { border-color: #e2e8f0; background: #fff; }

.is-completed { opacity: 0.6; filter: grayscale(0.5); }

/* レスポンシブ：スマホでの調整 */
@media (max-width: 850px) {
  .list-container { grid-template-columns: 1fr; }
}

@media (max-width: 640px) {
  .overview-view { padding-bottom: 40px; }
  .view-header { margin-bottom: 1.5rem; }
  .title-en { font-size: 1.4rem; }
  .subtitle { font-size: 0.85rem; }
  .list-container { gap: 1.5rem; }
  
  .dash-header {
    padding: 1rem;
  }
  .dash-body {
    padding: 1rem;
    gap: 1.25rem;
  }
  .day-count {
    font-size: 1.8rem;
  }
  .settings-area {
    grid-template-columns: 1fr;
    gap: 0.75rem;
  }
}
</style>
