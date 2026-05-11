# GraphQL Profile Dashboard

A lightweight front-end project that authenticates against the Zone01 platform and displays a student profile with key progress statistics and SVG-based graphs.

### Enjoy the app online

- [graphql](https://mrmarc0s.github.io/graphql/)

## What This Project Does

- Authenticates users with their Zone01 credentials
- Loads profile data through GraphQL queries
- Displays:
  - basic user information
  - total XP
  - audit ratio
  - completed projects
- Renders visual statistics using vanilla JavaScript and SVG

## Tech Stack

- HTML5
- CSS3
- Vanilla JavaScript (ES modules)
- Browser `fetch` API
- GraphQL API

## Getting Started

### Requirements

- Python 3 (used only to run a local static server)
- Internet access to reach `https://platform.zone01.gr`

### Run Locally

From the project root:

```bash
make run
```

Then open:

- [http://localhost:8000](http://localhost:8000)

## Configuration

API endpoints are configured in `static/src/config/config.js`.

Current defaults:

- Base URL: `https://platform.zone01.gr`
- Auth endpoint: `/api/auth/signin`
- GraphQL endpoint: `/api/graphql-engine/v1/graphql`

## Project Structure

```text
static/
  css/                           # Styling
  src/
    controllers/                 # Page-level logic
    services/                    # Auth, GraphQL, graph rendering
    config/                      # API URLs and localStorage keys
    utils/                       # UI and formatting helpers
index.html                       # Login page
Makefile                         # Local development command
profile.html                     # Profile dashboard page
```

## Notes

- This is a front-end only project; there is no local backend server.
- JWT tokens are stored in browser localStorage for session handling.

## Authors

- Christos Markos
