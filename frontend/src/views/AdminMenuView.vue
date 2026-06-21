<template>
  <div class="admin-menu-view">
    <div class="view-header">
      <div class="header-left">
        <h2 class="title-en">Admin Menu Dashboard ⚙️</h2>
        <p class="subtitle">研修メニューの追加、編集、削除を行います。変更は即座にデータベースおよび JSON に保存されます。</p>
      </div>
      <div class="header-right">
        <button @click="logout" class="btn-logout">ログアウト 🐾</button>
      </div>
    </div>

    <!-- 新規追加フォームカード -->
    <div class="form-card custom-card">
      <div class="card-header">
        <div class="header-left">
          <h3 class="card-title">{{ isEditingMode ? 'メニューを編集する ✏️' : '新規メニューを追加する ➕' }}</h3>
        </div>
        <div v-if="isEditingMode" class="header-right">
          <button @click="cancelEdit" class="btn-cancel">キャンセル</button>
        </div>
      </div>
      
      <form @submit.prevent="saveMenu" class="menu-form">
        <div class="form-grid">
          <div class="form-group">
            <label>研修名</label>
            <input type="text" v-model="form.name" placeholder="例: Gitで始めるバージョン管理" required />
          </div>
          <div class="form-group">
            <label>研修日数 (日)</label>
            <input type="number" v-model.number="form.days" min="1" placeholder="例: 3" required />
          </div>
          <div class="form-group">
            <label>難易度 (1〜5)</label>
            <select v-model.number="form.difficulty" required>
              <option :value="1">★☆☆☆☆</option>
              <option :value="2">★★☆☆☆</option>
              <option :value="3">★★★☆☆</option>
              <option :value="4">★★★★☆</option>
              <option :value="5">★★★★★</option>
            </select>
          </div>
          <div class="form-group">
            <label>ドキュメントURL</label>
            <input type="url" v-model="form.doc_link" placeholder="例: https://git-scm.com/..." />
          </div>
        </div>

        <div class="form-group">
          <label>研修の概要</label>
          <textarea v-model="form.summary" placeholder="研修で学ぶ概要を記入します..." rows="3" required></textarea>
        </div>

        <div class="form-grid">
          <div class="form-group">
            <label>身につくスキル (カンマ区切り)</label>
            <input type="text" v-model="form.skills" placeholder="例: Git基本操作, ブランチ管理" required />
          </div>
          <div class="form-group">
            <label>前提条件</label>
            <input type="text" v-model="form.prerequisites" placeholder="例: コマンドライン操作の基本" required />
          </div>
        </div>

        <div class="form-actions">
          <button type="submit" class="btn-submit">
            {{ isEditingMode ? '変更を保存する' : '新しく追加する' }}
          </button>
        </div>
      </form>
    </div>

    <!-- メニュー一覧 -->
    <div class="menu-list-section">
      <h3 class="section-title">登録済みの研修メニュー一覧 ({{ menus.length }})</h3>
      
      <div v-if="isLoading" class="status-msg">読み込み中...</div>
      
      <div v-else-if="menus.length === 0" class="status-msg empty-msg">
        登録されているメニューはありません。上のフォームから登録してください。
      </div>

      <div v-else class="menu-grid">
        <div v-for="menu in menus" :key="menu.id" class="menu-card custom-card">
          <div class="card-header">
            <div class="header-left">
              <span class="menu-id">ID: {{ menu.id }}</span>
              <h4 class="menu-title">{{ menu.name }}</h4>
            </div>
            <div class="header-right actions">
              <button @click="startEdit(menu)" class="btn-action-edit">✏️ 編集</button>
              <button @click="deleteMenu(menu)" class="btn-action-delete">🗑️ 削除</button>
            </div>
          </div>
          <div class="card-body">
            <div class="meta-row">
              <span class="days-badge">{{ menu.days }} Days</span>
              <div class="difficulty-stars">
                {{ "★".repeat(menu.difficulty) }}<span class="empty-stars">{{ "☆".repeat(5 - menu.difficulty) }}</span>
              </div>
            </div>
            <p class="summary-text">{{ menu.summary }}</p>
            <div class="detail-row">
              <strong>スキル:</strong> {{ menu.skills }}
            </div>
            <div class="detail-row">
              <strong>前提条件:</strong> {{ menu.prerequisites }}
            </div>
            <div v-if="menu.doc_link" class="detail-row">
              <strong>参考リンク:</strong> <a :href="menu.doc_link" target="_blank">{{ menu.doc_link }}</a>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';

const router = useRouter();
const menus = ref([]);
const isLoading = ref(true);
const isEditingMode = ref(false);

const initialForm = {
  id: null,
  name: '',
  days: 1,
  doc_link: '',
  summary: '',
  skills: '',
  prerequisites: '',
  difficulty: 3
};

const form = ref({ ...initialForm });

onMounted(async () => {
  const token = localStorage.getItem('admin_token');
  if (!token) {
    router.push('/admin/login');
    return;
  }
  axios.defaults.headers.common['X-Admin-Token'] = token;
  await fetchMenus();
});

const fetchMenus = async () => {
  isLoading.value = true;
  try {
    const res = await axios.get('/api/admin/menus');
    menus.value = res.data || [];
  } catch (e) {
    if (e.response && e.response.status === 401) {
      localStorage.removeItem('admin_token');
      router.push('/admin/login');
    } else {
      alert('メニュー一覧の取得に失敗しました');
    }
  } finally {
    isLoading.value = false;
  }
};

const logout = () => {
  localStorage.removeItem('admin_token');
  delete axios.defaults.headers.common['X-Admin-Token'];
  router.push('/admin/login');
};

const startEdit = (menu) => {
  isEditingMode.value = true;
  form.value = { ...menu };
  window.scrollTo({ top: 0, behavior: 'smooth' });
};

const cancelEdit = () => {
  isEditingMode.value = false;
  form.value = { ...initialForm };
};

const saveMenu = async () => {
  try {
    if (isEditingMode.value) {
      await axios.put(`/api/admin/menus/${form.value.id}`, form.value);
      alert('研修メニューを更新しました！');
    } else {
      await axios.post('/api/admin/menus', form.value);
      alert('研修メニューを追加しました！');
    }
    await fetchMenus();
    cancelEdit();
  } catch (e) {
    alert('保存に失敗しました');
  }
};

const deleteMenu = async (menu) => {
  if (confirm(`【重要】メニュー「${menu.name}」を削除しますか？\n※このメニューに紐づくすべての新人の計画、および進捗データも同時に削除されます。`)) {
    try {
      await axios.delete(`/api/admin/menus/${menu.id}`);
      alert('研修メニューを削除しました！');
      await fetchMenus();
      if (form.value.id === menu.id) {
        cancelEdit();
      }
    } catch (e) {
      alert('削除に失敗しました');
    }
  }
};
</script>

<style scoped>
* {
  box-sizing: border-box;
}

.admin-menu-view {
  max-width: 1000px;
  margin: 0 auto;
  padding-bottom: 80px;
}

.view-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2.5rem;
}

.title-en {
  font-size: 1.75rem;
  font-weight: 800;
  color: var(--text-main);
  margin-bottom: 0.5rem;
}

.subtitle {
  color: var(--text-muted);
  font-size: 0.95rem;
}

.btn-logout {
  background: #f1f5f9;
  border: 1px solid var(--border-color);
  color: var(--text-muted);
  padding: 0.5rem 1.25rem;
  border-radius: 8px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-logout:hover {
  background: #ef4444;
  color: white;
  border-color: #ef4444;
}

/* フォームカード */
.form-card {
  margin-bottom: 3rem;
  background: white;
}

.card-title {
  margin: 0;
  font-size: 1.2rem;
  font-weight: 800;
  color: var(--text-main);
}

.btn-cancel {
  background: transparent;
  border: none;
  color: var(--text-muted);
  font-weight: 700;
  cursor: pointer;
}

.btn-cancel:hover {
  text-decoration: underline;
}

.menu-form {
  padding: 2rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  text-align: left;
}

.form-group label {
  font-size: 0.8rem;
  font-weight: 700;
  color: var(--text-muted);
}

.form-group input,
.form-group select,
.form-group textarea {
  padding: 0.75rem;
  border: 2px solid var(--border-light);
  border-radius: 8px;
  font-size: 0.95rem;
  color: var(--text-main);
  outline: none;
  background: #f8fafc;
  transition: all 0.2s;
}

.form-group textarea {
  resize: vertical;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  border-color: var(--primary);
  background: white;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
}

.btn-submit {
  background: var(--primary);
  color: white;
  padding: 0.875rem 2rem;
  border-radius: 10px;
  font-weight: 700;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-submit:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 15px -4px rgba(79, 70, 229, 0.4);
}

/* 一覧セクション */
.section-title {
  font-size: 1.25rem;
  font-weight: 800;
  color: var(--text-main);
  margin-bottom: 1.5rem;
  text-align: left;
}

.status-msg {
  padding: 3rem;
  text-align: center;
  color: var(--text-muted);
  font-weight: 600;
}

.empty-msg {
  background: #f8fafc;
  border: 2px dashed var(--border-color);
  border-radius: 16px;
}

.menu-grid {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.menu-card {
  background: white;
  text-align: left;
}

.menu-id {
  font-size: 0.85rem;
  font-weight: 800;
  color: var(--text-muted);
  background: #e2e8f0;
  padding: 2px 6px;
  border-radius: 4px;
}

.menu-title {
  margin: 0;
  font-size: 1.15rem;
  font-weight: 800;
  color: var(--text-main);
}

.card-header .actions {
  display: flex;
  gap: 8px;
}

.btn-action-edit,
.btn-action-delete {
  padding: 4px 12px;
  border-radius: 6px;
  font-size: 0.8rem;
  font-weight: 700;
  border: 1px solid var(--border-color);
  background: white;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-action-edit:hover {
  background: #f1f5f9;
}

.btn-action-delete:hover {
  background: #fef2f2;
  color: #ef4444;
  border-color: #fee2e2;
}

.card-body {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.meta-row {
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

.difficulty-stars {
  font-size: 0.95rem;
  color: #0f172a;
}

.empty-stars {
  color: #cbd5e1;
}

.summary-text {
  margin: 0;
  color: #334155;
  font-size: 0.95rem;
  line-height: 1.6;
}

.detail-row {
  font-size: 0.875rem;
  color: var(--text-muted);
}

.detail-row strong {
  color: var(--text-main);
}

.detail-row a {
  color: var(--primary);
}

@media (max-width: 640px) {
  .view-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  .btn-logout {
    width: 100%;
  }
  .menu-form {
    padding: 1rem;
  }
}
</style>
