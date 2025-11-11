# Words - Vocabulary Learning API

A Go-based REST API for vocabulary learning with local caching and spaced repetition.

## Features

- 📚 **Word Lookup:** Fetch definitions, synonyms, examples, and pronunciation
- 💾 **Local-First:** SQLite cache for fast, offline-capable lookups
- 🔄 **Auto-Sync:** Automatically caches external API responses
- 🧠 **Spaced Repetition:** ✅ SM-2 algorithm for effective learning
- 👤 **Simple Auth:** ✅ Username-based accounts, no passwords
- 📊 **Progress Tracking:** ✅ Learning statistics, review history, and streaks
- 🎯 **Smart Scheduling:** ✅ Adaptive review intervals based on performance

## Quick Start

### 1. Install Dependencies
```bash
go get github.com/gin-gonic/gin github.com/mattn/go-sqlite3
```

### 2. Build
```bash
go build -o api cmd/api/main.go
```

### 3. Run
```bash
./api
# Server starts on http://localhost:8080
```

### 4. Test
```bash
# Look up a word (will fetch from API and cache locally)
curl http://localhost:8080/api/words/mordant | jq

# Second lookup will be instant (from cache)
curl http://localhost:8080/api/words/mordant | jq
```

## API Endpoints

### Phase 1 - Word Lookup ✅
- `GET /api/words/:word` - Look up word definition

### Phase 2 - Spaced Repetition ✅
**User Management:**
- `POST /api/users` - Create user account
- `GET /api/users/:username` - Get user details
- `GET /api/users/:username/stats` - Get learning statistics

**Vocabulary:**
- `POST /api/users/:username/words/:word` - Add word to study list
- `GET /api/users/:username/words` - Get all user's words (optional: `?status=learning|reviewing|mastered`)

**Reviews:**
- `GET /api/users/:username/review` - Get words due for review
- `POST /api/users/:username/review/:word` - Submit review rating (body: `{"quality": 0-5}`)
- `GET /api/users/:username/review/:word/history` - Get review history for a word

## Architecture

```
User Request → API Handler → Service Layer → Local DB (SQLite)
                                    ↓ (if not found)
                              External API (dictionaryapi.dev)
                                    ↓
                              Cache & Return
```

## Data Sources

- **Primary API:** [DictionaryAPI.dev](https://dictionaryapi.dev/) (Free, no API key)
- **Planned Datasets:** WordNet, Wordset, OPTED for bulk imports

## Database

SQLite with normalized schema:

**Phase 1 - Dictionary Data:**
- `words` - Base word entries
- `meanings` - Parts of speech
- `definitions` - Multiple definitions per meaning
- `phonetics` - Pronunciation guides
- `synonyms` / `antonyms` - Related words
- `source_urls` - Attribution

**Phase 2 - Learning System:**
- `users` - User accounts
- `user_words` - Words being studied (with SM-2 metadata)
- `review_history` - Complete audit trail of all reviews

## Development Status

✅ **Phase 1 - Lookup API** (Complete)
- [x] Project structure
- [x] Database schema
- [x] Local-first lookup logic
- [x] API endpoint
- [x] Install dependencies
- [x] Test integration
- [x] Dataset import script

✅ **Phase 2 - Spaced Repetition** (Complete)
- [x] User accounts (username-based)
- [x] SM-2 algorithm implementation
- [x] Review scheduling with quality ratings
- [x] Progress tracking and statistics
- [x] Word status management (learning/reviewing/mastered)
- [x] Review history tracking
- [x] Comprehensive API endpoints

## Project Structure

```
words/
├── cmd/api/              # Application entry point
├── internal/             # Private application code
│   ├── database/         # DB initialization & migrations
│   ├── handlers/         # HTTP request handlers
│   ├── models/           # Data structures
│   └── services/         # Business logic
├── pkg/                  # Public library code
│   └── dictionary/       # External API client
└── PROGRESS.md           # Detailed progress notes
```

## License

MIT

## Contributing

See `PROGRESS.md` for current status and next steps.
