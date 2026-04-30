# 🏨 Hotel Booking System

A RESTful API for hotel management built with **Go**, **Gin**, **GORM**, and **PostgreSQL**.

## 📋 Features

- Hotel & Room management (CRUD)
- User authentication (JWT)
- Reservations with date conflict validation
- Payment processing
- Favorite rooms per user
- Pagination & filters on all list endpoints
- Role-based access (admin / guest)
- Database migrations with golang-migrate

---

## 🚀 Quick Start

### 1. Clone & configure

```bash
git clone <your-repo>
cd hotel-booking
cp .env.example .env
# Edit .env with your database credentials
```

### 2. Install tools

```bash
make install-tools   # installs golang-migrate
go mod tidy
```

### 3. Create database & run migrations

```bash
make db-create
make migrate-up
```

### 4. Run server

```bash
make run
# Server starts on http://localhost:8080
```

---

## 📁 Project Structure

```
hotel-booking/
├── main.go                    # Entry point
├── Makefile                   # Migration & build commands
├── .env.example               # Environment template
├── config/
│   ├── config.go              # Config loader
│   └── database.go            # DB connection & AutoMigrate
├── models/
│   ├── user.go                # User model + DTOs
│   ├── hotel.go               # Hotel model + DTOs
│   ├── room.go                # Room model + DTOs
│   ├── reservation.go         # Reservation model + DTOs
│   ├── payment.go             # Payment model + DTOs
│   └── amenity.go             # Amenity + FavoriteRoom models
├── handlers/
│   ├── router.go              # All routes registered here
│   ├── auth_handler.go        # Register / Login / Me
│   ├── hotel_handler.go       # Hotels CRUD
│   ├── room_handler.go        # Rooms CRUD + Favorites
│   ├── reservation_handler.go # Reservations CRUD
│   ├── payment_handler.go     # Payments
│   └── amenity_handler.go     # Amenities
├── middleware/
│   └── auth.go                # JWT auth + AdminOnly middlewares
├── migrations/
│   ├── 000001_create_hotels.*
│   ├── 000002_create_users.*
│   ├── 000003_create_rooms.*
│   ├── 000004_create_reservations_payments.*
│   └── 000005_create_favorite_rooms.*
└── utils/
    ├── jwt.go                 # Token generation & parsing
    └── response.go            # Unified JSON response helpers
```

---

## 🔌 API Endpoints

### Auth
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | /auth/register | ❌ | Register new user |
| POST | /auth/login | ❌ | Login, get JWT |
| GET | /auth/me | ✅ | Get current user |

### Hotels
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | /hotels | ❌ | List hotels (filter: city, country, stars) |
| GET | /hotels/:id | ❌ | Get hotel by ID |
| POST | /hotels | 🔑 Admin | Create hotel |
| PUT | /hotels/:id | 🔑 Admin | Update hotel |
| DELETE | /hotels/:id | 🔑 Admin | Delete hotel |

### Rooms
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | /rooms | ❌ | List rooms (filter: hotel_id, type, price, capacity) |
| GET | /rooms/:id | ❌ | Get room by ID |
| POST | /rooms | 🔑 Admin | Create room |
| PUT | /rooms/:id | 🔑 Admin | Update room |
| DELETE | /rooms/:id | 🔑 Admin | Delete room |
| GET | /rooms/favorites | ✅ | Get my favorite rooms |
| PUT | /rooms/:id/favorites | ✅ | Add to favorites |
| DELETE | /rooms/:id/favorites | ✅ | Remove from favorites |

### Reservations
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | /reservations | ✅ | List my reservations |
| GET | /reservations/:id | ✅ | Get reservation |
| POST | /reservations | ✅ | Create reservation |
| PUT | /reservations/:id | ✅ | Update reservation |
| DELETE | /reservations/:id | ✅ | Cancel reservation |

### Payments
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | /payments | ✅ | Pay for reservation |
| GET | /payments/:id | ✅ | Get payment |

### Amenities
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | /amenities | ❌ | List amenities |
| POST | /amenities | 🔑 Admin | Create amenity |

---

## 📦 Example Requests

### Register
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com","password":"secret123"}'
```

### Create Hotel (admin token required)
```bash
curl -X POST http://localhost:8080/hotels \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Grand Hotel","address":"Main St 1","city":"Almaty","country":"Kazakhstan","stars":5}'
```

### List Rooms with filters & pagination
```bash
curl "http://localhost:8080/rooms?hotel_id=1&type=suite&min_price=100&page=1&limit=5"
```

### Create Reservation
```bash
curl -X POST http://localhost:8080/reservations \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"room_id":1,"check_in":"2025-06-01T14:00:00Z","check_out":"2025-06-05T12:00:00Z","guests":2}'
```

---

## 🔧 Makefile Commands

```bash
make run               # Start server
make build             # Build binary
make migrate-up        # Apply all migrations
make migrate-down      # Rollback all
make migrate-down-1    # Rollback last migration
make migrate-status    # Current version
make migrate-create NAME=add_indexes  # New migration
make migrate-force VERSION=3          # Force version
make db-create         # Create database
make db-drop           # Drop database
make install-tools     # Install golang-migrate CLI
make tidy              # go mod tidy
```

---

## 🌿 Git Branch Strategy (Practicums)

```
main
├── practicum/2-rest-api        # In-memory REST API practice
├── practicum/3-gin-api         # Gin-based bookstore API
├── practicum/4-postgresql      # PostgreSQL + GORM
├── practicum/5-authentication  # JWT authentication
├── practicum/6-favorites       # Favorites feature
└── practicum/7-migrations      # golang-migrate + Makefile
```
