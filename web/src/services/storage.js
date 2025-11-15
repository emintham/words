// Local storage service for managing user session
const STORAGE_KEYS = {
  CURRENT_USER: 'words_current_user',
  DUE_WORDS_CACHE: 'words_due_cache',
  LAST_SYNC: 'words_last_sync',
};

class StorageService {
  // User Management
  setCurrentUser(user) {
    localStorage.setItem(STORAGE_KEYS.CURRENT_USER, JSON.stringify(user));
  }

  getCurrentUser() {
    const user = localStorage.getItem(STORAGE_KEYS.CURRENT_USER);
    return user ? JSON.parse(user) : null;
  }

  clearCurrentUser() {
    localStorage.removeItem(STORAGE_KEYS.CURRENT_USER);
  }

  // Cache Management (for offline support)
  cacheDueWords(words) {
    localStorage.setItem(STORAGE_KEYS.DUE_WORDS_CACHE, JSON.stringify(words));
    localStorage.setItem(STORAGE_KEYS.LAST_SYNC, new Date().toISOString());
  }

  getCachedDueWords() {
    const words = localStorage.getItem(STORAGE_KEYS.DUE_WORDS_CACHE);
    return words ? JSON.parse(words) : null;
  }

  getLastSync() {
    return localStorage.getItem(STORAGE_KEYS.LAST_SYNC);
  }

  clearCache() {
    localStorage.removeItem(STORAGE_KEYS.DUE_WORDS_CACHE);
    localStorage.removeItem(STORAGE_KEYS.LAST_SYNC);
  }
}

export default new StorageService();
