
# Ohara

Ohara is a Go backend for turning websites into structured, queryable knowledge. It crawls a site, extracts meaningful content, chunks it, embeds it, and exposes search-style APIs that let apps ask questions about the website in natural language.

This project was originally bootstrapped with [skelgo](https://github.com/JerryJeager/skelgo), which gave the base Go API structure and conventions.

---

## Why Ohara?

Modern websites hold a lot of business knowledge: pricing, services, product details, FAQs, policies, and customer-facing messaging. Ohara helps turn that content into a usable knowledge layer for AI-driven experiences.

At a high level, Ohara does three things:

- indexes website content from a public URL
- stores the content as chunked, embeddable documents
- answers questions against that site through a query API

This makes the backend useful for:

- AI-powered site assistants
- customer support search
- website knowledge retrieval
- smart widgets or chat experiences powered by a brand’s own content

---

## Core capabilities

### Website ingestion

- accepts a website URL and metadata for a brand or business
- crawls the public site content
- extracts text from pages and filters noisy markup
- stores indexed page references and site profile details

### AI profile generation

- analyzes site content to infer summary, tone, and topical context
- stores that data with the website record
- helps downstream experiences better match the brand voice and domain

### Embedding and retrieval

- chunks page content into smaller segments
- embeds them for semantic search/retrieval
- supports document queries by website
- powers retrieval for embedded UI experiences and chat workflows

### Security and ownership model

- user-based ownership of websites
- JWT auth for protected endpoints
- allowed-origin checks for embeddable access
- local development toggles and website configuration controls

---

## Tech stack

- Go
- Gin HTTP framework
- GORM + PostgreSQL
- golang-migrate for database migrations
- Google GenAI / Gemini for profile generation and AI-related features
- HTML crawling and extraction via Go tools such as Colly and goquery
- Headless browser fallback for rendered or JavaScript-heavy pages via Rod

---

## Project structure

```text
.
├── cmd/                     # API startup and route registration
├── config/                 # environment and database/AI initialization
├── docs/                   # project docs, migration notes
├── internal/
│   ├── db/migrations/      # PostgreSQL migration files
│   ├── http/               # API handlers for users, websites, and documents
│   ├── models/             # database/domain models
│   ├── service/            # business logic and repository/service layers
│   └── utils/              # auth, token, email, and helper utilities
├── middleware/              # auth and CORS middleware
├── manualwire/              # manual dependency wiring
├── main.go                 # app entry point
├── go.mod                  # Go module definition
├── go.sum                  # dependency lockfile
├── README.md               # project overview
└── docs/                   # additional documentation
```

---

## Getting started

### 1. Install dependencies

```bash
go mod tidy
```

### 2. Configure environment

Create a local environment file and add the required values for your database, JWT secret, and AI provider configuration.

```bash
cp .env.example .env
```

Then update the variables in `.env` for your environment.

### 3. Run database migrations

This project uses [golang-migrate](https://github.com/golang-migrate/migrate).

Create a migration:

```bash
migrate create -ext psql -dir internal/db/migrations -seq MIGRATION_NAME
```

Run migrations:

```bash
migrate -path internal/db/migrations -database "postgresql://postgres:YOURPASSWORD@localhost:5432/YOURDB?sslmode=disable" -verbose up
```

### 4. Start the API

```bash
go run main.go
```

The server runs on port `8080` by default unless `PORT` is set.

---

## Main API flows

### User auth

- `POST /api/v1/users/auth/google` — Google auth flow
- `GET /api/v1/users` — fetch current user
- `POST /api/v1/users/auth/refresh` — refresh tokens

### Website management

- `POST /api/v1/websites` — create a website record
- `GET /api/v1/websites` — list websites for the authenticated user
- `GET /api/v1/websites/:website_id/indexed` — list indexed pages
- `PATCH /api/v1/websites/:website_id/local-dev` — enable or disable local dev mode
- `PATCH /api/v1/websites/:website_id/allowed-origins` — update allowed origins
- `DELETE /api/v1/websites/:website_id` — delete a website

### Document and retrieval endpoints

- `POST /api/v1/documents/query/:website_id` — query indexed content for a website
- `GET /api/v1/websites/:website_id/widget` — widget embed endpoint

---

## Operational notes

Ohara is designed to be the backend layer behind a website intelligence experience. The app records the site metadata, crawls the public content, enriches it with AI-generated profile information, and exposes retrieval interfaces that can power a frontend widget, assistant, or embedded search experience.

The project is intentionally structured around a few clear responsibilities:

- HTTP layer for API interactions
- service layer for business logic
- database layer for persistence and migrations
- AI-assisted enrichment for profile and retrieval workflows

---

## Documentation

- [docs/migrations.md](docs/migrations.md) — migration reference and database workflow notes

---

## Contributing

Contributions are welcome. If you are extending the product, keep the project’s core theme in mind: website intelligence, retrieval, and accessible AI-powered knowledge experiences.

---

## License

MIT
