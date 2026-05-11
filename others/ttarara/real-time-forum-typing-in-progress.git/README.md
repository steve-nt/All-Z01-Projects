# Tech Forum

A modern, real-time forum application built with Go, SQLite, and vanilla JavaScript. Features include user authentication, posts, comments, likes/dislikes, private messaging, notifications, and real-time updates via WebSocket.

## ✨ Features

- **User Authentication**
  - Email/password registration and login
  - Secure session management
  - Password reset functionality

- **Posts & Comments**
  - Create, edit, and delete posts
  - Comment on posts
  - Like/dislike posts and comments
  - Category filtering

- **Real-time Communication**
  - Private messaging between users
  - Real-time message delivery via WebSocket
  - Online/offline user status
  - Users Status box showing recently active users (last hour)

- **Notifications**
  - Real-time notifications for likes, comments, and messages
  - Notification count badge
  - Mark notifications as read

- **User Interface**
  - Modern, responsive dark theme
  - Drag-and-drop Users Status box
  - Minimizable sidebar
  - Real-time updates without page refresh

## 🛠️ Tech Stack

### Backend
- **Go 1.23+** - Main server language
- **SQLite** - Database
- **Gorilla WebSocket** - Real-time communication

### Frontend
- **Vanilla JavaScript** - No framework dependencies
- **Tailwind CSS** - Styling
- **WebSocket API** - Real-time updates

## 📋 Prerequisites

- Go 1.23 or higher
- Node.js and npm (for CSS building)
- SQLite3

## 🚀 Installation

### 1. Clone the repository

```bash
git clone <repository-url>
cd real-time-forum
```

### 2. Install Go dependencies

```bash
go mod download
```

### 3. Install frontend dependencies

```bash
npm install
```

### 4. Build CSS

```bash
npm run build-css-prod
```

Or for development with watch mode:

```bash
npm run build-css
```

### 5. Initialize the database

The database will be automatically created when you first run the application.

## ⚙️ Configuration

### 🔧 Environment Variables

The application uses the following default settings:

- **Server Port**: `8081` (configurable in `main.go`)
- **Database**: `forum.db` (SQLite database file)
- **Session Duration**: 7 days

## Running the Application

### Development Mode

1. Build the CSS (if not already built):
   ```bash
   npm run build-css-prod
   ```

2. Run the Go server:
   ```bash
   go run main.go
   ```

3. Open your browser and navigate to:
   ```
   http://localhost:8081
   ```

### Production Mode

1. Build the application:
   ```bash
   go build -o forum main.go
   ```

2. Run the executable:
   ```bash
   ./forum
   ```

### Using Docker

1. Build and run with Docker Compose:
   ```bash
   docker-compose up --build
   ```

2. The application will be available at:
   ```
   http://localhost:8080
   ```

## 📁 Project Structure

```
real-time-forum/
├── frontend/
│   ├── app.js              # Main frontend JavaScript
│   ├── index.html          # Main HTML file
│   ├── src/
│   │   └── input.css       # Tailwind CSS source
│   └── dist/
│       └── output.css      # Compiled CSS
├── internals/
│   ├── database/
│   │   ├── database.go     # Database connection and initialization
│   │   ├── table.sql       # Database schema
│   │   └── sqlstruct.go    # SQL structs
│   ├── handlers/
│   │   ├── routes.go       # Route setup
│   │   ├── authHandler.go  # Authentication handlers
│   │   ├── postHandler.go  # Post CRUD operations
│   │   ├── commentHandler.go # Comment operations
│   │   ├── messageHandler.go # Private messaging
│   │   ├── notificationHandler.go # Notifications
│   │   └── wsHandler.go    # WebSocket handler
│   └── utils/
│       └── utils.go        # Utility functions
├── main.go                 # Application entry point
├── go.mod                  # Go dependencies
├── package.json            # Node.js dependencies
└── docker-compose.yml      # Docker configuration
```

## API Endpoints

### 🔐 Authentication
- `POST /login` - User login
- `POST /register` - User registration
- `GET /logout` - User logout
- `GET /api/auth/status` - Check authentication status

### Posts
- `GET /api/posts` - Get all posts
- `GET /api/post?id={id}` - Get single post
- `POST /new-post` - Create new post
- `POST /api/posts/edit` - Edit post
- `POST /api/posts/delete` - Delete post
- `GET /api/posts/filtered` - Get filtered posts by category

### Comments
- `GET /api/comments?post_id={id}` - Get comments for a post
- `POST /api/comments/create` - Create comment
- `POST /api/comments/edit` - Edit comment
- `POST /api/comments/delete` - Delete comment

### Likes/Dislikes
- `POST /api/posts/like` - Like/dislike a post
- `POST /api/comments/like` - Like/dislike a comment

### 💌 Messages
- `GET /api/messages/users` - Get list of users for messaging (shows online users and offline users active in last hour)
- `GET /api/messages?user_id={id}` - Get messages with a user
- `POST /api/messages/send` - Send a message

### Notifications
- `GET /api/notifications` - Get user notifications
- `POST /api/notifications/mark-read` - Mark notification as read
- `POST /api/notifications/mark-all-read` - Mark all as read
- `GET /api/notifications/count` - Get unread count

### WebSocket
- `WS /ws` - WebSocket connection for real-time updates

## 💡 Key Features Explained

### Users Status Box
- Shows online users and offline users who were active in the last hour
- Displays unread message counts
- Click on a user to start a private chat
- Minimizable and draggable interface
- Only visible when logged in

### Real-time Updates
- New messages appear instantly
- User online/offline status updates in real-time
- Notifications delivered immediately
- No page refresh required

### Private Messaging
- Direct messages between users
- Message history with timestamps
- Unread message indicators
- Real-time message delivery

## Development

### Building CSS

For development (with watch mode):
```bash
npm run build-css
```

For production (minified):
```bash
npm run build-css-prod
```

### 🗄️ Database Schema

The database schema is defined in `internals/database/table.sql` and includes:
- Users table
- Posts table
- Comments table
- LikesDislikes table
- PrivateMessages table
- Notifications table
- Sessions table

## Security Features

- Password hashing with bcrypt
- Session-based authentication
- SQL injection prevention (parameterized queries)
- XSS protection (HTML escaping)
- CSRF protection via session validation

## 📜 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## 👨‍💻 Authors
For questions or issues, please contact us on Discord:

[Charoula Tarara](https://discordapp.com/users/1242540766879023160)
[Georgia Marouli](https://discordapp.com/users/1277216244910522371)
[Andriana Stas](https://discordapp.com/users/780150798927134740)
[Sofia Busho](https://discordapp.com/users/1276592724979613697)

> © 2025 Xaroula Tarara, Georgia Marouli, Andriana Stas and Sofia Busho for Zone01Athens Projects
