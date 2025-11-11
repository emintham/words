# Words - Vocabulary Learning API

A Go-based REST API for vocabulary learning with local caching and spaced repetition.

## Features

- 📚 **Word Lookup:** Fetch definitions, synonyms, examples, and pronunciation
- 💾 **Local-First:** SQLite cache for fast, offline-capable lookups
- 🔄 **Auto-Sync:** Automatically caches external API responses
- 🧠 **Spaced Repetition:** (Phase 2) SM-2 algorithm for effective learning
- 👤 **Simple Auth:** (Phase 2) Username-based accounts, no passwords

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

### Phase 1 (Current)
- `GET /api/words/:word` - Look up word definition

### Phase 2 (Planned)
- `POST /api/users` - Create user account
- `GET /api/users/:username/review` - Get words due for review
- `POST /api/users/:username/review/:word` - Submit review rating
- `POST /api/users/:username/words/:word` - Add word to study list

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
- `words` - Base word entries
- `meanings` - Parts of speech
- `definitions` - Multiple definitions per meaning
- `phonetics` - Pronunciation guides
- `synonyms` / `antonyms` - Related words
- `source_urls` - Attribution

## Development Status

✅ **Phase 1 - Lookup API** (In Progress)
- [x] Project structure
- [x] Database schema
- [x] Local-first lookup logic
- [x] API endpoint
- [ ] Install dependencies
- [ ] Test integration
- [ ] Dataset import script

⏳ **Phase 2 - Spaced Repetition** (Planned)
- [ ] User accounts
- [ ] SM-2 algorithm
- [ ] Review scheduling
- [ ] Progress tracking

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
