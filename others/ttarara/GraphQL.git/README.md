# GraphQL Profile Page

A modern, interactive profile page that displays user information from the Zone01 GraphQL API, featuring beautiful SVG graphs and statistics.

## Features

- **Authentication**: Secure JWT-based login with support for both username and email
- **User Profile**: Display user information, XP, progress, and grades
- **Interactive Graphs**: Multiple SVG-based statistics graphs including:
  - XP Progress Over Time (Line Chart)
  - XP by Project (Bar Chart)
  - Pass/Fail Ratio (Pie Chart)
  - Audit Ratio (Horizontal Bar Chart)
- **Modern UI/UX**: Beautiful, responsive design with smooth animations
- **GraphQL Integration**: Uses normal queries, nested queries, and queries with arguments

## Project Structure

```
graphql/
├── index.html          # Login page
├── profile.html        # Profile page
├── css/
│   └── styles.css     # Main stylesheet
├── js/
│   ├── app.js         # Main application logic
│   ├── auth.js        # Authentication utilities
│   ├── graphql.js     # GraphQL query functions
│   └── graphs.js      # SVG graph generation
├── netlify.toml       # Netlify configuration
└── README.md          # This file
```

## Getting Started

### How to run it:

   **Using Python 3**
   ```bash
   python3 -m http.server 8000
   ```
   Then open your browser and go to: `http://localhost:8000`
   To stop the server, press `Ctrl+C` in the terminal

**Note:** You must use a web server (not just open the HTML file directly) because the application makes API calls that require proper HTTP requests.

### Authentication

The application uses JWT authentication:
- Supports both `username:password` and `email:password` authentication

## GraphQL Queries

The application demonstrates three types of GraphQL queries:

1. **Normal Queries**: Simple queries like fetching user information
2. **Nested Queries**: Queries that traverse relationships (e.g., transaction with user)
3. **Queries with Arguments**: Queries using filters and variables (e.g., object by ID)

## Statistics Graphs

All graphs are created using SVG and include:

- **XP Over Time**: Shows cumulative XP progression with area chart
- **XP by Project**: Bar chart displaying top 10 projects by XP
- **Pass/Fail Ratio**: Pie chart showing success/failure distribution
- **Audit Ratio**: Horizontal bar chart categorizing audit results

## Hosting on Netlify

The site will be live at `https://ttarara-graphql.netlify.app/`

## Technologies Used

- HTML5
- CSS3 (with CSS Variables and Flexbox/Grid)
- Vanilla JavaScript (ES6+)
- GraphQL
- SVG for data visualization
- JWT for authentication

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)


## Authors

Tarara Theocharoula 🧑🏻‍💻 🧑🏻‍💻
