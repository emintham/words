import { useState } from 'react';
import api from '../services/api';

function AddWord({ user, onAdded }) {
  const [word, setWord] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [wordDetails, setWordDetails] = useState(null);

  const handleSearch = async (e) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setWordDetails(null);
    setLoading(true);

    try {
      const data = await api.getWord(word.toLowerCase());
      setWordDetails(data);
    } catch (err) {
      setError(err.message || 'Word not found. Please try another word.');
    } finally {
      setLoading(false);
    }
  };

  const handleAdd = async () => {
    setError('');
    setSuccess('');
    setLoading(true);

    try {
      await api.addWordToStudyList(user.username, word.toLowerCase());
      setSuccess(`"${word}" has been added to your study list!`);
      setWord('');
      setWordDetails(null);
      onAdded();
    } catch (err) {
      setError(err.message || 'Failed to add word. It might already be in your list.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="content-card">
      <h2>➕ Add New Word</h2>
      <p style={{ color: '#666', marginBottom: '30px' }}>
        Search for a word and add it to your study list.
      </p>

      {error && <div className="error-message">{error}</div>}
      {success && (
        <div style={{
          background: '#d4edda',
          color: '#155724',
          padding: '12px',
          borderRadius: '8px',
          marginBottom: '20px',
          border: '1px solid #c3e6cb'
        }}>
          {success}
        </div>
      )}

      <form onSubmit={handleSearch}>
        <div className="form-group">
          <label htmlFor="word">Word</label>
          <input
            type="text"
            id="word"
            value={word}
            onChange={(e) => setWord(e.target.value)}
            placeholder="Enter a word (e.g., serendipity)"
            required
          />
        </div>

        <button type="submit" className="btn" disabled={loading}>
          {loading ? 'Searching...' : 'Search Word'}
        </button>
      </form>

      {wordDetails && (
        <div style={{ marginTop: '30px', padding: '20px', background: '#f8f9fa', borderRadius: '8px' }}>
          <h3 style={{ margin: '0 0 10px 0' }}>{wordDetails.word}</h3>
          {wordDetails.phonetic && (
            <p style={{ color: '#666', margin: '0 0 20px 0' }}>{wordDetails.phonetic}</p>
          )}

          {wordDetails.meanings && wordDetails.meanings.map((meaning, idx) => (
            <div key={idx} style={{ marginBottom: '20px' }}>
              <h4 style={{ color: '#667eea', marginBottom: '10px' }}>
                {meaning.partOfSpeech}
              </h4>
              {meaning.definitions && meaning.definitions.slice(0, 2).map((def, defIdx) => (
                <div key={defIdx} style={{ marginBottom: '10px' }}>
                  <p style={{ margin: '5px 0' }}>
                    <strong>{defIdx + 1}.</strong> {def.definition}
                  </p>
                  {def.example && (
                    <p style={{ marginLeft: '20px', fontStyle: 'italic', color: '#666' }}>
                      "{def.example}"
                    </p>
                  )}
                </div>
              ))}
            </div>
          ))}

          <button
            className="btn btn-success"
            onClick={handleAdd}
            disabled={loading}
            style={{ marginTop: '20px' }}
          >
            {loading ? 'Adding...' : 'Add to Study List'}
          </button>
        </div>
      )}
    </div>
  );
}

export default AddWord;
