# Thankly — Project Overview

## 1. Project Overview

**Gratitude Journal** is a subscription-based gratitude micro-blogging platform where users create an account and write **up to 3 gratitude entries every day**.

The product is designed like a simple personal blog rather than a traditional complex SaaS dashboard. Each day's gratitude entries appear as a small daily journal/post, and users can access their complete journal history from the day they joined.

The platform will have a **Free plan** and a **Premium subscription plan**.

### Core Idea

> **Write 3 things you're grateful for every day. Build a better mindset, one day at a time.**

---

# 2. Main Goals

The project has two goals:

### Product Goal

Create a simple, enjoyable daily gratitude-writing experience.

### Developer/Portfolio Goal

Demonstrate real-world **Go backend development** skills, including:

* REST API development
* Authentication & authorization
* PostgreSQL database design
* Business-rule enforcement
* Subscription/payment integration
* Redis
* Background jobs
* API validation
* Error handling
* Testing
* Docker
* Production deployment

---

# 3. Core Features

## User Authentication

Users can:

* Create an account
* Login
* Logout
* Refresh authentication token
* Update profile
* Change password

Authentication will be handled by the Go backend.

---

## Daily Gratitude

Each authenticated user can write **maximum 3 gratitude entries per day**.

Example:

### September 1

1. I'm grateful for my family.
2. I'm grateful for a peaceful morning.
3. I'm grateful for getting closer to my goals.

The backend must enforce the daily limit.

The frontend should never be responsible for enforcing this rule.

---

# 4. Journal History

Users can see their **complete gratitude history**.

There is no 30-day or 90-day history limitation.

Example:

```text
September 2026

Sep 01
  ☀️ Peaceful morning
  ❤️ Family
  🚀 Working toward my goals

Aug 31
  ☕ Morning coffee
  🌳 Evening walk
  😊 Good conversation

Aug 30
  ...
```

Users can:

* Browse previous entries
* Filter by date
* Search gratitude
* View individual days
* Edit their own entries
* Delete their own entries

---

# 5. Blog-like Experience

The application should feel more like a **personal micro-blog** than an admin dashboard.

Each day can be represented as a journal post:

```text
                 September 1, 2026

              Things I'm grateful for

        ☀️ Peaceful morning

        ❤️ Spending time with family

        🚀 Making progress toward my goals
```

This makes Astro a strong choice for the frontend.

---

# 6. Privacy

Each user's journal is private by default.

A user can only:

* Read their own gratitude
* Edit their own gratitude
* Delete their own gratitude

The Go backend must verify ownership before returning or modifying journal data.

Example:

```text
GET /api/gratitudes/:id

        ↓

Authenticate user
        ↓
Find gratitude
        ↓
Check user_id
        ↓
Owner?
   ↓        ↓
 YES        NO
Return    403
```

---

# 7. Subscription Model

The application will have two plans.

## Free Plan

* Account creation
* 3 gratitude entries per day
* Complete journal history
* Basic profile
* Basic streak

## Premium Plan

Premium should provide additional value without changing the core identity of the product.

Possible Premium features:

* Advanced gratitude statistics
* Calendar view
* Gratitude search
* Advanced streak analytics
* Journal export
* PDF/JSON export
* Custom journal themes
* Daily reminders
* Additional personalization

The **3 gratitude/day rule can remain for both plans** to keep the product concept consistent.

---

# 8. Subscription Flow

```text
User
  │
  ▼
Pricing Page
  │
  ▼
Choose Premium
  │
  ▼
Payment Provider
  │
  ▼
Successful Payment
  │
  ▼
Webhook
  │
  ▼
Go Backend
  │
  ▼
Update Subscription
  │
  ▼
PostgreSQL
```

The backend should never trust the frontend to determine whether a user is Premium.

Subscription status must come from the backend/database.

---

# 9. Recommended Tech Stack

## Frontend

### Astro

Astro will be used for:

* Landing page
* Blog-like journal interface
* Public pages
* Pricing page
* Authentication UI
* User profile
* Journal pages
* SEO
* Fast page delivery

Astro is especially suitable because the product is content-oriented and has a blog-like experience.

### React

React can optionally be used inside Astro for highly interactive components:

* Journal editor
* Search
* Calendar
* Statistics
* Subscription UI

Use React only where interactivity is actually needed.

### Tailwind CSS

For:

* Responsive UI
* Typography
* Components
* Dark mode
* Consistent design system

---

# 10. Backend

## Go

Go will be the primary backend language.

Recommended framework:

### Gin

Use Gin for:

* HTTP routing
* REST APIs
* Middleware
* Request handling
* Authentication middleware
* Error handling

The project should keep the business logic independent from Gin as much as possible.

---

# 11. Database

## PostgreSQL

PostgreSQL will be the primary database.

Why PostgreSQL?

* Excellent relational database
* Strong consistency
* Transactions
* Foreign keys
* Indexes
* Excellent Go support
* Suitable for subscription and user data
* Production-ready

---

# 12. Database Structure

Initial schema:

```text
users
-----
id
name
email
password_hash
created_at
updated_at
```

```text
gratitudes
---------
id
user_id
content
gratitude_date
created_at
updated_at
```

```text
subscriptions
-------------
id
user_id
provider
provider_customer_id
provider_subscription_id
plan
status
current_period_start
current_period_end
created_at
updated_at
```

Possible future tables:

```text
refresh_tokens
daily_reminders
journal_exports
user_preferences
subscription_events
```

---

# 13. Redis

Redis will be used where it provides real value.

Possible uses:

* API rate limiting
* Authentication-related temporary data
* Caching
* Background job queues
* Temporary subscription/payment state

Redis should not replace PostgreSQL as the source of truth for important journal data.

---

# 14. API Design

Example REST API:

### Authentication

```text
POST   /api/auth/register
POST   /api/auth/login
POST   /api/auth/refresh
POST   /api/auth/logout
```

### User

```text
GET    /api/me
PATCH  /api/me
PATCH  /api/me/password
```

### Gratitude

```text
POST   /api/gratitudes
GET    /api/gratitudes
GET    /api/gratitudes/:id
PATCH  /api/gratitudes/:id
DELETE /api/gratitudes/:id
```

### Journal

```text
GET /api/journal/today
GET /api/journal/history
GET /api/journal/:date
```

### Subscription

```text
GET  /api/subscription
POST /api/subscription/checkout
POST /api/subscription/cancel
```

### Webhook

```text
POST /api/webhooks/payment
```

---

# 15. Daily 3-Gratitude Rule

The most important business rule:

> A user cannot create more than 3 gratitude entries for the same calendar day.

Backend flow:

```text
POST /api/gratitudes
        │
        ▼
Authenticate user
        │
        ▼
Validate content
        │
        ▼
Check today's entries
        │
        ▼
Count >= 3?
   │             │
  YES            NO
   │             │
Return 409      Create
```

The database query can be optimized using an index such as:

```text
(user_id, gratitude_date)
```

The implementation should also consider concurrent requests so users cannot bypass the limit by sending multiple requests simultaneously.

---

# 16. Security

Security requirements:

* Passwords must be hashed
* Never store plain-text passwords
* JWT access token
* Refresh token mechanism
* Authentication middleware
* Ownership checks
* Input validation
* Rate limiting
* Secure HTTP headers
* CORS configuration
* SQL injection protection
* Payment webhook signature verification

---

# 17. Project Architecture

Recommended Go structure:

```text
backend/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── auth/
│   ├── user/
│   ├── gratitude/
│   ├── subscription/
│   ├── journal/
│   ├── middleware/
│   ├── database/
│   └── worker/
│
├── migrations/
│
├── pkg/
│
├── docs/
│
├── Dockerfile
├── go.mod
└── go.sum
```

Keep the project modular rather than putting everything into one large `main.go`.

---

# 18. Docker

Development environment:

```text
Docker Compose
│
├── Go API
├── PostgreSQL
└── Redis
```

This allows the complete backend environment to be started consistently.

---

# 19. Testing

The project should include:

### Unit Tests

* Gratitude daily limit
* Authentication
* Subscription rules
* Validation
* Journal service

### Integration Tests

* PostgreSQL
* API endpoints
* Authentication flow
* Gratitude CRUD

Important test:

```text
User creates gratitude #1 → SUCCESS
User creates gratitude #2 → SUCCESS
User creates gratitude #3 → SUCCESS
User creates gratitude #4 → REJECTED
```

---

# 20. Future Features

After the MVP:

* Daily email reminders
* Push notifications
* Mood tracking
* Gratitude streak
* Monthly statistics
* Calendar heatmap
* Public journal option
* Following other users
* Likes/reactions
* Comments
* AI-generated gratitude insights
* Mobile application
* Admin dashboard

These should be added later rather than making the first version unnecessarily complicated.

---

# 21. MVP Scope

The first version should focus on:

```text
1. User registration/login
2. User profile
3. Create 3 gratitude/day
4. Edit/delete gratitude
5. Complete journal history
6. Date-based journal view
7. Free/Premium subscription
8. Payment webhook
9. PostgreSQL
10. Go REST API
11. Astro frontend
12. Docker
13. Tests
```

This is enough to create a strong portfolio project.

---

# 22. Final Technology Stack

| Layer            | Technology                     |
| ---------------- | ------------------------------ |
| Frontend         | Astro                          |
| Interactive UI   | React                          |
| Styling          | Tailwind CSS                   |
| Backend          | Go                             |
| HTTP Framework   | Gin                            |
| API              | REST                           |
| Database         | PostgreSQL                     |
| Cache/Rate Limit | Redis                          |
| Authentication   | JWT + Refresh Token            |
| Payments         | Stripe                         |
| Documentation    | Swagger/OpenAPI                |
| Containerization | Docker                         |
| Testing          | Go testing + integration tests |
| CI/CD            | GitHub Actions                 |
| Deployment       | VPS / AWS / Railway / Render   |

---

# 23. Portfolio Positioning

For your CV/GitHub, describe it as:

**Gratitude Journal — Subscription-based Micro-Blogging Platform**

> Built a subscription-based gratitude journaling platform using Go, PostgreSQL, Redis, and Astro. Implemented secure authentication, user-level authorization, daily gratitude limits, complete journal history, REST APIs, subscription management, payment webhooks, Docker-based development, and automated testing.

This positioning makes the project demonstrate **real backend engineering skills**, rather than looking like a simple CRUD blog.
