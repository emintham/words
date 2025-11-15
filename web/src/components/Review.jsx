import { useState, useEffect } from 'react';
import api from '../services/api';

function Review({ user, onComplete }) {
  const [dueWords, setDueWords] = useState([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [wordDetails, setWordDetails] = useState(null);
  const [loading, setLoading] = useState(true);
  const [reviewing, setReviewing] = useState(false);
  const [showDefinition, setShowDefinition] = useState(false);

  useEffect(() => {
    loadDueWords();
  }, []);

  const loadDueWords = async () => {
    setLoading(true);
    try {
      const data = await api.getDueWords();
      setDueWords(data.words || []);
      if (data.words && data.words.length > 0) {
        loadWordDetails(data.words[0].word);
      }
    } catch (error) {
      console.error('Failed to load due words:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadWordDetails = async (word) => {
    try {
      const data = await api.getWord(word);
      setWordDetails(data);
    } catch (error) {
      console.error('Failed to load word details:', error);
    }
  };

  const handleReview = async (quality) => {
    setReviewing(true);
    try {
      await api.submitReview(dueWords[currentIndex].word, quality);

      // Move to next word
      if (currentIndex + 1 < dueWords.length) {
        const nextIndex = currentIndex + 1;
        setCurrentIndex(nextIndex);
        loadWordDetails(dueWords[nextIndex].word);
        setShowDefinition(false);
      } else {
        // All done!
        onComplete();
        setDueWords([]);
      }
    } catch (error) {
      console.error('Failed to submit review:', error);
    } finally {
      setReviewing(false);
    }
  };

  const qualityOptions = [
    { value: 0, label: '0 - Total blackout', description: 'Complete failure to recall' },
    { value: 1, label: '1 - Incorrect', description: 'Incorrect response; correct one remembered' },
    { value: 2, label: '2 - Hard', description: 'Correct response with serious difficulty' },
    { value: 3, label: '3 - Good', description: 'Correct response with difficulty' },
    { value: 4, label: '4 - Easy', description: 'Correct response with hesitation' },
    { value: 5, label: '5 - Perfect', description: 'Perfect response' },
  ];

  if (loading) {
    return (
      <div className="content-card">
        <p>Loading reviews...</p>
      </div>
    );
  }

  if (dueWords.length === 0) {
    return (
      <div className="content-card">
        <h2>🎉 All Done!</h2>
        <p style={{ color: '#666', marginTop: '20px' }}>
          You have no words due for review right now. Great job keeping up with your studies!
        </p>
        <p style={{ color: '#666', marginTop: '10px' }}>
          Come back later or add more words to your study list.
        </p>
      </div>
    );
  }

  const currentWord = dueWords[currentIndex];

  return (
    <div className="review-container">
      <div className="content-card" style={{ marginBottom: '20px', textAlign: 'center' }}>
        <p style={{ color: '#666', margin: 0 }}>
          Card {currentIndex + 1} of {dueWords.length}
        </p>
      </div>

      <div className="review-card">
        <h2>{currentWord.word}</h2>

        {!showDefinition ? (
          <>
            <p style={{ color: '#666', marginBottom: '30px' }}>
              Do you remember this word?
            </p>
            <button
              className="btn"
              onClick={() => setShowDefinition(true)}
            >
              Show Definition
            </button>
          </>
        ) : (
          <>
            {wordDetails && (
              <div style={{ textAlign: 'left', marginBottom: '30px' }}>
                {wordDetails.phonetic && (
                  <p style={{ color: '#666', marginBottom: '20px' }}>
                    {wordDetails.phonetic}
                  </p>
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
              </div>
            )}

            <h3 style={{ marginTop: '40px', marginBottom: '20px', textAlign: 'center' }}>
              How well did you remember?
            </h3>

            <div className="quality-buttons">
              {qualityOptions.map((option) => (
                <button
                  key={option.value}
                  className="quality-btn"
                  onClick={() => handleReview(option.value)}
                  disabled={reviewing}
                >
                  <div>{option.label}</div>
                  <div style={{ fontSize: '12px', color: '#666', marginTop: '5px' }}>
                    {option.description}
                  </div>
                </button>
              ))}
            </div>

            {reviewing && (
              <p style={{ marginTop: '20px', color: '#666' }}>
                Submitting review...
              </p>
            )}
          </>
        )}
      </div>
    </div>
  );
}

export default Review;
