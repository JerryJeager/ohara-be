
# Ohara

This project was generated using [skelgo](https://github.com/JerryJeager/skelgo), a CLI tool for scaffolding Go REST API projects with [Gin](https://github.com/gin-gonic/gin), [GORM](https://gorm.io/), and PostgreSQL.

---

## 📁 Project Structure

```
.
├── .env, .env.example      # Environment variable files
├── .gitignore              # Git ignore rules
├── go.mod, go.sum          # Go module files
├── main.go                 # Application entry point
├── cmd/                    # Application startup logic (app.go)
├── config/                 # Database and environment configuration
├── docs/                   # Project documentation and migration guides
├── internal/               # Internal application logic
│   ├── db/                 # Database migrations and related files
│   ├── http/               # HTTP handlers and controllers
│   ├── models/             # Data models (e.g., User, OTP)
│   ├── service/            # Business logic and service layer
│   └── utils/              # Utility functions and helpers
│       └── emails/         # Email-related utilities
├── manualwire/             # Manual dependency wiring (DI)
├── middleware/             # Custom Gin middleware (auth, CORS)
```

### Package/Folder Explanations

- **cmd/**: Contains the main application logic and entry point for running the server.
- **config/**: Handles loading environment variables and setting up the database connection.
- **docs/**: Includes documentation files, such as migration instructions.
- **internal/db/**: Contains database migration files for schema changes.
- **internal/http/**: Implements HTTP handlers and controllers for routing API requests.
- **internal/models/**: Defines data models and structs used throughout the application.
- **internal/service/**: Contains business logic and service interfaces/implementations.
- **internal/utils/**: Utility functions for tasks like OTP generation, token handling, etc.
  - **internal/utils/emails/**: Helpers for sending and verifying emails.
- **manualwire/**: Provides manual dependency injection wiring for controllers and services.
- **middleware/**: Custom middleware for authentication, CORS, and other HTTP concerns.


---

## 🚀 Getting Started

### 1. Install Dependencies

```sh
go mod tidy
```

### 2. Set Up Environment

Copy `.env.example` to `.env` and update your environment variables as needed.

### 3. Database Migrations

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations.

- **Create a migration:**
  ```sh
  migrate create -ext psql -dir internal/db/migrations -seq MIGRATION_NAME
  ```

- **Run migrations:**
  ```sh
  migrate -path internal/db/migrations -database "postgresql://postgres:YOURPASSWORD@localhost:5432/YOURDB?sslmode=disable" -verbose up
  ```

---

## 🏃 Running the Server

```sh
go run main.go
```

The server will start on port `8080` by default.

---

## 🛠️ API Endpoints

- `POST /api/v1/users/signup` — Register a new user
- `POST /api/v1/users/verify-email` — Verify user email
- `POST /api/v1/users/login` — User login

---

## 📚 Documentation

See [`docs/migrations.md`](docs/migrations.md) for migration instructions and more.

---

## 📝 Customization

- Add your own models in `internal/models/`.
- Extend services and controllers as needed.
- Add new middleware in `middleware/`.
- Utility functions can be placed in `internal/utils/`.

---

## 🤝 Contributing

Feel free to open issues or submit pull requests!

---

## License

MIT
	