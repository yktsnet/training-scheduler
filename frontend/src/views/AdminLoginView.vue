<template>
  <div class="admin-login-container">
    <div class="glass-card">
      <div class="login-header">
        <span class="lock-icon">🔒</span>
        <h2>Administrator Login</h2>
        <p>研修メニュー管理用のパスワードを入力してください</p>
        <div class="demo-tip">💡 デモ用のデフォルトパスワード: <code>admin123</code></div>
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="input-wrapper">
          <input 
            type="password" 
            v-model="password" 
            placeholder="パスワードを入力..." 
            required 
            ref="passwordInput"
            class="styled-input"
            :disabled="isLoading"
          />
        </div>

        <div v-if="errorMsg" class="error-msg">
          {{ errorMsg }}
        </div>

        <button type="submit" class="submit-btn" :disabled="isLoading">
          <span v-if="isLoading" class="spinner"></span>
          <span v-else>ログイン 🐾</span>
        </button>
      </form>

      <div class="back-link">
        <router-link to="/">← 一般画面へ戻る</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';

const password = ref('');
const errorMsg = ref('');
const isLoading = ref(false);
const passwordInput = ref(null);
const router = useRouter();

onMounted(() => {
  passwordInput.value?.focus();
});

const handleLogin = async () => {
  errorMsg.value = '';
  isLoading.value = true;

  try {
    const res = await axios.post('/api/admin/login', { password: password.value });
    localStorage.setItem('admin_token', res.data.token);
    router.push('/admin/menus');
  } catch (e) {
    if (e.response && e.response.data && e.response.data.error) {
      errorMsg.value = 'パスワードが正しくありません';
    } else {
      errorMsg.value = 'サーバー接続エラーが発生しました';
    }
  } finally {
    isLoading.value = false;
  }
};
</script>

<style scoped>
* {
  box-sizing: border-box;
}

.admin-login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 70vh;
  padding: 1rem;
}

.glass-card {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(226, 232, 240, 0.8);
  border-radius: 24px;
  padding: 3rem;
  width: 100%;
  max-width: 460px;
  box-shadow: 0 20px 40px -15px rgba(15, 23, 42, 0.1);
  text-align: center;
  transition: all 0.3s ease;
}

.lock-icon {
  font-size: 3rem;
  display: inline-block;
  margin-bottom: 1rem;
}

h2 {
  font-size: 1.5rem;
  font-weight: 800;
  color: var(--text-main);
  margin-bottom: 0.5rem;
}

p {
  color: var(--text-muted);
  font-size: 0.875rem;
  margin-bottom: 1.5rem;
}

.demo-tip {
  background: #f1f5f9;
  border-radius: 8px;
  padding: 0.5rem 1rem;
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--text-muted);
  margin-bottom: 2rem;
  display: inline-block;
}
.demo-tip code {
  color: var(--primary);
  font-weight: 700;
  background: white;
  padding: 2px 6px;
  border-radius: 4px;
  border: 1px solid var(--border-light);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.styled-input {
  width: 100%;
  padding: 1rem 1.25rem;
  border: 2px solid var(--border-light);
  border-radius: 12px;
  font-size: 1rem;
  color: var(--text-main);
  outline: none;
  background: #f8fafc;
  transition: all 0.2s;
}

.styled-input:focus {
  border-color: var(--primary);
  background: white;
  box-shadow: 0 0 0 4px rgba(79, 70, 229, 0.15);
}

.error-msg {
  color: #ef4444;
  font-size: 0.85rem;
  font-weight: 700;
  text-align: left;
  background: #fef2f2;
  padding: 0.75rem 1rem;
  border-radius: 10px;
  border: 1px solid #fee2e2;
}

.submit-btn {
  background: var(--primary);
  color: white;
  padding: 1rem;
  border-radius: 12px;
  font-weight: 700;
  border: none;
  cursor: pointer;
  font-size: 1rem;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px -6px rgba(79, 70, 229, 0.4);
}

.submit-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.back-link {
  margin-top: 2rem;
}

.back-link a {
  color: var(--text-muted);
  text-decoration: none;
  font-size: 0.875rem;
  font-weight: 600;
  transition: color 0.2s;
}

.back-link a:hover {
  color: var(--primary);
}

.spinner {
  width: 20px;
  height: 20px;
  border: 3px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
