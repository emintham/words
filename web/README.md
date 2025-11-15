# Words - Vocabulary Learning PWA

A Progressive Web App for learning vocabulary with spaced repetition, built with React and Vite.

## Features

- 📱 **Progressive Web App** - Install on mobile and desktop
- 🔄 **Spaced Repetition** - SM-2 algorithm for optimal learning
- 📊 **Progress Tracking** - Track your learning statistics
- ⚡ **Offline Support** - Review words even without internet
- 🎯 **Quality Ratings** - Rate your recall from 0-5
- 📖 **Word Management** - Add and organize your vocabulary

## Quick Start

### 1. Start the Go API Backend

```bash
# From the project root
cd /home/user/words
./api
# API runs on http://localhost:9090
```

### 2. Start the Web App

```bash
# From the web directory
cd web
pnpm install
pnpm run dev
# Web app runs on http://localhost:5173
```

### 3. Open in Browser

Navigate to `http://localhost:5173` and create a username to get started!

## Project Structure

```
web/
├── src/
│   ├── components/       # React components
│   │   ├── Login.jsx    # Login screen
│   │   ├── Dashboard.jsx # Main dashboard
│   │   ├── Stats.jsx    # Statistics view
│   │   ├── Review.jsx   # Review system with SM-2
│   │   ├── WordList.jsx # Word list view
│   │   └── AddWord.jsx  # Add new words
│   ├── services/        # API and storage services
│   │   ├── api.js       # Backend API client
│   │   └── storage.js   # Local storage management
│   ├── App.jsx          # Main app component
│   ├── App.css          # Styles
│   └── main.jsx         # Entry point
└── vite.config.js       # Vite + PWA configuration
```

## Technologies

- **React 18** - UI framework
- **Vite** - Build tool and dev server
- **Vite PWA Plugin** - Service worker and manifest
- **Workbox** - Service worker runtime caching

## PWA Features

- **Offline Support**: Cached words remain available
- **Installable**: Add to home screen on mobile/desktop
- **Fast**: Service worker caches assets and API responses

## Development

```bash
pnpm run dev      # Start dev server
pnpm run build    # Build for production
pnpm run preview  # Preview production build
```

## Browser Support

- Chrome/Edge 90+
- Firefox 88+
- Safari 14+
- iOS Safari 14+
- Android Chrome 90+

## License

MIT
