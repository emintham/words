import { useState, useEffect } from 'react';
import api from '../services/api';

function WordList({ user, onRefresh }) {
  const [words, setWords] = useState([]);
  const [filter, setFilter] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadWords();
  }, [filter]);

  const loadWords = async () => {
    setLoading(true);
    try {
      const data = await api.getUserWords(user.username, filter);
      setWords(data.words || []);
    } catch (error) {
      console.error('Failed to load words:', error);
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (dateStr) => {
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric'
    });
  };

  return (
    <div className="content-card">
      <div style={{ marginBottom: '20px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h2>📖 My Words</h2>
        <select
          value={filter}
          onChange={(e) => setFilter(e.target.value)}
          style={{
            padding: '8px 16px',
            borderRadius: '8px',
            border: '2px solid #e0e0e0',
            fontSize: '14px'
          }}
        >
          <option value="">All Words</option>
          <option value="learning">Learning</option>
          <option value="reviewing">Reviewing</option>
          <option value="mastered">Mastered</option>
        </select>
      </div>

      {loading ? (
        <p>Loading words...</p>
      ) : words.length === 0 ? (
        <p style={{ color: '#666', textAlign: 'center', padding: '40px' }}>
          No words found. Add some words to get started!
        </p>
      ) : (
        <ul className="word-list">
          {words.map((word) => (
            <li key={word.id} className="word-item">
              <div className="word-info">
                <h3>{word.word}</h3>
                <div className="word-meta">
                  <span>Added: {formatDate(word.added_at)}</span>
                  <span>Next Review: {formatDate(word.next_review_date)}</span>
                  <span>Interval: {word.interval_days} days</span>
                  <span>Ease: {word.ease_factor.toFixed(2)}</span>
                </div>
              </div>
              <span className={`badge badge-${word.status}`}>
                {word.status}
              </span>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

export default WordList;
