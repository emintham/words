import { useState, useEffect } from 'react';
import Login from './components/Login';
import Dashboard from './components/Dashboard';
import storage from './services/storage';
import api from './services/api';
import './App.css';

function App() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Validate session with server
    const checkSession = async () => {
      try {
        const response = await api.getCurrentUser();
        setUser(response.user);
        storage.setCurrentUser(response.user);
      } catch (err) {
        // Session invalid or expired, clear local storage
        storage.clearCurrentUser();
        setUser(null);
      } finally {
        setLoading(false);
      }
    };

    checkSession();
  }, []);

  const handleLogin = (userData) => {
    storage.setCurrentUser(userData);
    setUser(userData);
  };

  const handleLogout = async () => {
    try {
      await api.logout();
    } catch (err) {
      console.error('Logout error:', err);
    } finally {
      storage.clearCurrentUser();
      storage.clearCache();
      setUser(null);
    }
  };

  if (loading) {
    return (
      <div className="loading">
        <h2>Loading...</h2>
      </div>
    );
  }

  return (
    <div className="app">
      {user ? (
        <Dashboard user={user} onLogout={handleLogout} />
      ) : (
        <Login onLogin={handleLogin} />
      )}
    </div>
  );
}

export default App;
