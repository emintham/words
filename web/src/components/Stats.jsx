function Stats({ stats }) {
  return (
    <>
      <div className="stats-grid">
        <div className="stat-card">
          <h3>Total Words</h3>
          <p className="stat-value">{stats.total_words}</p>
        </div>
        <div className="stat-card">
          <h3>Due Today</h3>
          <p className="stat-value">{stats.due_today}</p>
        </div>
        <div className="stat-card">
          <h3>Learning</h3>
          <p className="stat-value">{stats.learning}</p>
        </div>
        <div className="stat-card">
          <h3>Reviewing</h3>
          <p className="stat-value">{stats.reviewing}</p>
        </div>
        <div className="stat-card">
          <h3>Mastered</h3>
          <p className="stat-value">{stats.mastered}</p>
        </div>
        <div className="stat-card">
          <h3>Total Reviews</h3>
          <p className="stat-value">{stats.total_reviews}</p>
        </div>
      </div>

      <div className="content-card">
        <h2>📈 Your Progress</h2>
        <p style={{ color: '#666', marginTop: '10px' }}>
          Keep up the great work! Consistent daily reviews will help you master new vocabulary.
        </p>

        {stats.due_today > 0 && (
          <div style={{ marginTop: '20px', padding: '20px', background: '#fff3cd', borderRadius: '8px' }}>
            <strong>🎯 You have {stats.due_today} word{stats.due_today > 1 ? 's' : ''} due for review today!</strong>
            <p style={{ margin: '10px 0 0 0' }}>
              Click on the "Review" tab to start practicing.
            </p>
          </div>
        )}

        {stats.total_words === 0 && (
          <div style={{ marginTop: '20px', padding: '20px', background: '#d1ecf1', borderRadius: '8px' }}>
            <strong>👋 Welcome!</strong>
            <p style={{ margin: '10px 0 0 0' }}>
              Get started by adding your first word in the "Add Word" tab.
            </p>
          </div>
        )}
      </div>
    </>
  );
}

export default Stats;
