import { useState, useEffect } from 'react';
import api from '../services/api';
import Stats from './Stats';
import Review from './Review';
import WordList from './WordList';
import AddWord from './AddWord';

function Dashboard({ user, onLogout }) {
  const [currentView, setCurrentView] = useState('stats');
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    try {
      const data = await api.getUserStats();
      setStats(data);
    } catch (error) {
      console.error('Failed to load stats:', error);
    } finally {
      setLoading(false);
    }
  };

  const refreshData = () => {
    loadStats();
  };

  return (
    <div className="dashboard">
      <div className="dashboard-header">
        <h1>Welcome, {user.username}!</h1>
        <button className="btn btn-secondary" onClick={onLogout}>
          Logout
        </button>
      </div>

      <div className="dashboard-nav">
        <button
          className={`nav-btn ${currentView === 'stats' ? 'active' : ''}`}
          onClick={() => setCurrentView('stats')}
        >
          📊 Stats
        </button>
        <button
          className={`nav-btn ${currentView === 'review' ? 'active' : ''}`}
          onClick={() => setCurrentView('review')}
        >
          🎯 Review ({stats?.due_today || 0})
        </button>
        <button
          className={`nav-btn ${currentView === 'words' ? 'active' : ''}`}
          onClick={() => setCurrentView('words')}
        >
          📖 My Words
        </button>
        <button
          className={`nav-btn ${currentView === 'add' ? 'active' : ''}`}
          onClick={() => setCurrentView('add')}
        >
          ➕ Add Word
        </button>
      </div>

      {loading ? (
        <div className="content-card">
          <p>Loading...</p>
        </div>
      ) : (
        <>
          {currentView === 'stats' && <Stats stats={stats} />}
          {currentView === 'review' && (
            <Review user={user} onComplete={refreshData} />
          )}
          {currentView === 'words' && (
            <WordList user={user} onRefresh={refreshData} />
          )}
          {currentView === 'add' && (
            <AddWord user={user} onAdded={refreshData} />
          )}
        </>
      )}
    </div>
  );
}

export default Dashboard;
