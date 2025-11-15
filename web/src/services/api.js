// API service for connecting to Go backend
const API_BASE_URL = '/api';

class ApiService {
  async request(endpoint, options = {}) {
    const url = `${API_BASE_URL}${endpoint}`;
    const config = {
      credentials: 'include', // Send cookies with requests
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    };

    try {
      const response = await fetch(url, config);

      // Try to parse JSON response
      let data;
      try {
        data = await response.json();
      } catch (parseError) {
        // If JSON parsing fails, throw a generic error
        throw new Error('Failed to communicate with server. Please ensure the backend is running.');
      }

      if (!response.ok) {
        throw new Error(data.error || 'Request failed');
      }

      return data;
    } catch (error) {
      console.error(`API Error [${endpoint}]:`, error);
      throw error;
    }
  }

  // Authentication
  async login(username) {
    return this.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username }),
    });
  }

  async logout() {
    return this.request('/auth/logout', {
      method: 'POST',
    });
  }

  async getCurrentUser() {
    return this.request('/auth/me');
  }

  // User Management
  async createUser(username) {
    return this.request('/users', {
      method: 'POST',
      body: JSON.stringify({ username }),
    });
  }

  async getUser() {
    return this.request('/user');
  }

  async getUserStats() {
    return this.request('/user/stats');
  }

  // Word Lookup
  async getWord(word) {
    return this.request(`/words/${word}`);
  }

  // Vocabulary Management
  async addWordToStudyList(word) {
    return this.request(`/words/${word}`, {
      method: 'POST',
    });
  }

  async getUserWords(status = null) {
    const query = status ? `?status=${status}` : '';
    return this.request(`/words${query}`);
  }

  // Review System
  async getDueWords() {
    return this.request('/review');
  }

  async submitReview(word, quality) {
    return this.request(`/review/${word}`, {
      method: 'POST',
      body: JSON.stringify({ quality }),
    });
  }

  async getReviewHistory(word) {
    return this.request(`/review/${word}/history`);
  }
}

export default new ApiService();
