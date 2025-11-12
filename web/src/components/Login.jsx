import { useState } from 'react';
import api from '../services/api';

function Login({ onLogin }) {
  const [username, setUsername] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      // Try to get existing user first
      let user;
      try {
        user = await api.getUser(username);
      } catch (err) {
        // User doesn't exist, create new one
        user = await api.createUser(username);
      }

      onLogin(user);
    } catch (err) {
      setError(err.message || 'Failed to login. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-container">
      <div className="login-box">
        <h1>📚 Words</h1>
        <p>Learn vocabulary with spaced repetition</p>

        {error && <div className="error-message">{error}</div>}

        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="username">Username</label>
            <input
              type="text"
              id="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Enter your username"
              required
              minLength={3}
              maxLength={20}
              pattern="[a-zA-Z0-9_]+"
              title="Username can only contain letters, numbers, and underscores"
            />
          </div>

          <button type="submit" className="btn" disabled={loading}>
            {loading ? 'Loading...' : 'Continue'}
          </button>
        </form>

        <p style={{ marginTop: '20px', fontSize: '14px', color: '#999' }}>
          No password needed! Just enter a username to get started.
        </p>
      </div>
    </div>
  );
}

export default Login;
