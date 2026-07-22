# Social Network

A Facebook-like social network application built with a JavaScript frontend framework and a Go backend service, utilizing SQLite, WebSockets, and Docker containerization.

---

## Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
  - [Authentication](#authentication)
  - [Followers & Profiles](#followers--profiles)
  - [Posts & Media](#posts--media)
  - [Groups & Events](#groups--events)
  - [Real-Time Chat](#real-time-chat)
  - [Notifications](#notifications)
- [Architecture & Tech Stack](#architecture--tech-stack)
  - [Frontend](#frontend)
  - [Backend](#backend)
  - [Allowed Packages](#allowed-packages)
- [Database & Migrations](#database--migrations)
- [Containerization (Docker)](#containerization-docker)
- [Learning Objectives](#learning-objectives)
- [Getting Started](#getting-started)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

This project requires building a full-stack social network web application. The application features real-time messaging, user authentication, customizable privacy settings, media uploads, group management, and interactive notifications.

---

## Key Features

### Authentication
- **Registration & Login:** Secure authentication using sessions and cookies.
- **User Fields:**
  - Required: Email, Password, First Name, Last Name, Date of Birth
  - Optional: Avatar/Image, Nickname, About Me
- **Session Persistence:** Users remain logged in across sessions until explicitly logging out.

### Followers & Profiles
- **Follow System:** Send follow requests. Auto-follow applies for public profiles; manual approval required for private profiles.
- **Profile Types:**
  - **Public:** Profile info, user activity, posts, and followers/following lists are visible to all users.
  - **Private:** Profile info and activity are restricted to followers only.
- **Privacy Toggle:** Users can switch their profile setting between public and private at any time.

### Posts & Media
- **Post Creation:** Users can publish posts and comments with optional image/GIF attachments.
- **Supported Media:** Handling of `JPEG`, `PNG`, and `GIF` file types saved to the database/file system.
- **Post Privacy Levels:**
  - `Public`: Visible to all users.
  - `Almost Private`: Visible to followers only.
  - `Private`: Visible to specifically selected followers.

### Groups & Events
- **Group Management:** Create groups with a title and description, invite users, or respond to join requests.
- **Group Interaction:** Members can create group-specific posts, comments, and events.
- **Events:** Events contain a title, description, date/time, and response options (`Going` / `Not going`).

### Real-Time Chat
- **Private Messaging:** Real-time messaging between users who follow each other (or messaging users with public profiles).
- **Features:** Emoji support included.
- **Group Chat:** Dedicated real-time chat rooms for group members.

### Notifications
- **Universal Access:** Visible across all pages, distinct from private chat messages.
- **Notification Triggers:**
  - Incoming follow requests for private profiles.
  - Group invitations.
  - Group join requests (for group creators).
  - New group event creations.

---

## Architecture & Tech Stack

### Frontend
- Developed using a JavaScript Framework (e.g., Next.js, Vue.js, Svelte, Mithril, or another framework of your choice).
- Focus on responsiveness, client-side rendering, and API communication with the backend.

### Backend
- Modular backend app (written in Go) handling server logic, HTTP protocol middleware, and database operations.
- Option to use Caddy or build a custom web server in Go.

### Allowed Packages
- Go standard library
- `github.com/gorilla/websocket`
- `golang-migrate` / `sql-migration` / `migration`
- `sqlite3`
- `golang.org/x/crypto/bcrypt`
- `github.com/gofrs/uuid` or `github.com/google/uuid`

---

## Database & Migrations

Data persistence is managed using **SQLite**. Database schema modifications must be implemented using database migrations.

### Project Directory Structure
```text
backend
├── pkg
│   ├── db
│   │   ├── migrations
│   │   │   └── sqlite
│   │   │       ├── 000001_create_users_table.down.sql
│   │   │       ├── 000001_create_users_table.up.sql
│   │   │       ├── 000002_create_posts_table.down.sql
│   │   │       └── 000002_create_posts_table.up.sql
│   │   └── sqlite
│   │       └── sqlite.go
│   │
│   └── ...other_pkgs.go
│
└── server.go
```

- Migrations path example: `file://backend/pkg/db/migrations/sqlite`
- `sqlite.go` handles database connectivity, applying migrations on startup, and helper routines.

---

## Containerization (Docker)

The application is split into two separate Docker containers:

1. **Backend Container:** Runs server logic, listens for client requests, executes migrations, and interacts with the SQLite database.
2. **Frontend Container:** Serves client-side static assets and single-page application code to browser clients.

*Ensure backend and frontend containers expose their respective application ports for inter-container communication and host access.*

---

## Learning Objectives

This project provides hands-on experience with:
- Authentication via sessions and cookies
- Docker containerization and image creation
- Database schema design, SQL manipulation, and migrations
- Basic encryption techniques (e.g., password hashing with bcrypt)
- Real-time bidirectional communication using WebSockets

---
