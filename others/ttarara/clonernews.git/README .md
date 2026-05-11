# HackerNews Live UI

A minimal Hacker News client that fetches and displays stories, jobs, and polls using the [Hacker News Firebase API](https://github.com/HackerNews/API).

## 🚀 Features
- **Feeds**: Stories, Jobs, Polls
- **Incremental loading**: Only loads posts when needed (infinite scroll + "Load more")
- **Filters**:
  - Search by title (debounced typing)
  - Minimum score filter
- **Live Updates**:
  - Automatically polls `/updates.json` every 5 seconds
  - Inserts new posts **at the top** of the feed
  - Shows a banner indicating how many new items were added
- **Comments**: Nested comments with lazy loading
- **Responsive**: Works on desktop and mobile

## 📂 Project Structure
```
.
├── index.html
├── style.css
├── spinner.svg        # loader icon (inside assets/ if preferred)
├── js/
│   ├── api.js         # API layer (fetch from HN)
│   ├── store.js       # global app state
│   ├── ui.js          # rendering logic
│   ├── live.js        # auto live updates
│   └── main.js        # app bootstrap & glue code
└── assets/
    └── spinner.svg    
```

## 🛠️ Getting Started

1. Clone or download the repo:
   ```bash
   git clone https://platform.zone01.gr/git/ttarara/clonernews.git
   cd clonernews
   ```

2. Start a local server (Python 3):
   ```bash
   python3 -m http.server 5173
   ```
   Then open: [http://localhost:5173](http://localhost:5173)

   > ℹ️ You cannot open `index.html` directly with `file://` because ES modules need a web server.

3. Browse the feeds:
   - Use **Tabs** to switch between Stories, Jobs, Polls.
   - Use **Filters** to search or restrict by score.
   - Scroll down to load more.
   - New posts appear automatically at the top every ~5s.

## 🌀 Loader
- By default, the app uses `spinner.svg` as the loader icon.
- If you prefer, you can move it to `/assets/spinner.svg` and update the path in `index.html`.

## 📜 License
MIT — free to use, modify, and share.
