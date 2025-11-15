// API service for connecting to Go backend
const API_BASE_URL = '/api';

class ApiService {
  async request(endpoint, options = {}) {
    const url = `${API_BASE_URL}${endpoint}`;
    const config = {
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

  // User Management
  async createUser(username) {
    return this.request('/users', {
      method: 'POST',
      body: JSON.stringify({ username }),
    });
  }

  async getUser(username) {
    return this.request(`/users/${username}`);
  }

  async getUserStats(username) {
    return this.request(`/users/${username}/stats`);
  }

  // Word Lookup
  async getWord(word) {
    return this.request(`/words/${word}`);
  }

  // Vocabulary Management
  async addWordToStudyList(username, word) {
    return this.request(`/users/${username}/words/${word}`, {
      method: 'POST',
    });
  }

  async getUserWords(username, status = null) {
    const query = status ? `?status=${status}` : '';
    return this.request(`/users/${username}/words${query}`);
  }

  // Review System
  async getDueWords(username) {
    return this.request(`/users/${username}/review`);
  }

  async submitReview(username, word, quality) {
    return this.request(`/users/${username}/review/${word}`, {
      method: 'POST',
      body: JSON.stringify({ quality }),
    });
  }

  async getReviewHistory(username, word) {
    return this.request(`/users/${username}/review/${word}/history`);
  }
}

export default new ApiService();
