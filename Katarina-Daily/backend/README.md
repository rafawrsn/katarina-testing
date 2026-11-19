# Fashion E-Commerce Backend

A production-oriented REST API backend for a fashion e-commerce platform
(UNIQLO-style catalog: products sold as color/size variants), built in
Go with Gin, GORM and PostgreSQL.

## Stack

- **Language:** Go 1.22
- **Framework:** Gin
- **ORM:** GORM
- **Database:** PostgreSQL only
- **Auth:** JWT (HS256) + bcrypt password hashing
- **Architecture:** Handler → Service → Repository → Database, with
  DTOs at the HTTP boundary so `models` never leak password hashes or
  internal-only fields into responses.

## Project layout

```
backend/
├── cmd/server/main.go        # entrypoint: wiring + graceful shutdown
├── internal/
│   ├── config/                # env var loading
│   ├── database/               # GORM connection + AutoMigrate
│   ├── middleware/             # auth, RBAC, CORS, rate limit, logging, recovery
│   ├── models/                 # GORM models (1 file per table, mostly)
│   ├── repositories/           # DB access, incl. transaction-aware *Tx methods
│   ├── services/                # business logic (checkout, stock, webhook idempotency...)
│   ├── handlers/                # HTTP layer: parse, validate, call service, respond
│   ├── dto/                     # request/response shapes, never expose models directly
│   └── routes/                  # route registration
├── pkg/
│   ├── shipping/                # ShippingService abstraction + mock provider
│   └── payment/                 # PaymentGateway abstraction + Midtrans implementation
├── .env.example
├── Dockerfile
└── docker-compose.yml
```

## Running locally

### Option A — Docker Compose (recommended)

```bash
cp .env.example .env   # optional, docker-compose.yml already sets dev defaults
docker compose up --build
```

This starts Postgres 16 and the API on `http://localhost:8080`.
`AutoMigrate` runs automatically on boot — no separate migration step.

### Option B — Go + local Postgres

```bash
cp .env.example .env
# edit .env: point DATABASE_URL at your local Postgres

go mod tidy      # fetches dependencies; requires network access
go run ./cmd/server
```

Health check: `GET /health`.

> **This code was written without network access to a Go toolchain**,
> so `go mod tidy`/`go build` have not been run against it in this
> environment. The versions pinned in `go.mod` are real published
> versions, and the code was cross-checked with static analysis for
> import correctness, matching method/field names, and balanced
> syntax — but you should run `go mod tidy && go build ./...` after
> download and fix any straggling typos before deploying.

## Environment variables

See `.env.example` for the full list. The important ones:

| Variable | Purpose |
|---|---|
| `DATABASE_URL` | PostgreSQL connection string (Railway sets this automatically in production) |
| `JWT_SECRET` | HMAC signing key for auth tokens — set a real secret in production |
| `PAYMENT_SERVER_KEY` / `PAYMENT_CLIENT_KEY` / `PAYMENT_ENVIRONMENT` | Midtrans sandbox/production credentials |
| `SHIPPING_ORIGIN_CITY` / `SHIPPING_ORIGIN_PROVINCE` | Store's ship-from location, used for rate calculation |
| `CORS_ALLOWED_ORIGINS` | Comma-separated list of allowed frontend origins |

## API surface

All responses use the envelope:

```json
{ "success": true, "message": "...", "data": {}, "meta": {} }
{ "success": false, "message": "...", "error": "SOME_CODE" }
```

Key routes (see `internal/routes/routes.go` for the full list):

- `POST /api/v1/auth/register`, `POST /api/v1/auth/login`, `GET /api/v1/auth/me`
- `GET/POST/PUT/DELETE /api/v1/categories`, `/api/v1/products`, `/api/v1/variants` (writes are admin-only)
- `GET/POST/PUT/DELETE /api/v1/addresses` (customer, own data only)
- `GET/POST/PUT/DELETE /api/v1/cart`, `/api/v1/cart/items/:id`
- `GET /api/v1/shipping/rates?address_id=...` — dynamic shipping quotes
- `POST /api/v1/orders` — **this is checkout**: turns the current cart into an order and opens a payment transaction
- `GET /api/v1/orders`, `GET /api/v1/orders/:id`, `POST /api/v1/orders/:id/cancel`
- `POST /api/v1/payments/webhook` — Midtrans server-to-server notification (no JWT; authenticity comes from signature verification)
- `GET/POST/DELETE /api/v1/wishlist` — customer's saved products (`DELETE /api/v1/wishlist/:product_id`)
- `GET /api/v1/products/:id/reviews` (public), `POST /api/v1/products/:id/reviews` (customer, purchase-verified), `DELETE /api/v1/reviews/:id` (own review)
- `/api/v1/admin/*` — products, categories, variants, orders, shipments, users, review moderation, all behind `role=admin`

## Design decisions worth knowing about

- **Checkout never trusts the client for money.** `POST /api/v1/orders`
  re-reads live variant prices/stock, recomputes weight, re-quotes
  shipping, and re-validates any voucher — the request body only
  contains an address id, courier/service choice, and an optional
  voucher code.
- **Stock is decremented once, at confirmed payment**, not at
  checkout. `internal/services/payment_service.go` locks the order and
  payment rows (`SELECT ... FOR UPDATE`) inside a transaction before
  touching anything, so a duplicate webhook delivery is a safe no-op —
  no double order, no double stock decrement, no duplicate payment
  record. If stock runs out between checkout and payment confirmation
  (a race with another customer), the order is cancelled and the
  payment is flagged `FAILED` with a logged note that a manual refund
  is required — there's no auto-refund call in this MVP.
- **Order shipping fields and order-item product/price fields are
  snapshots**, copied in at checkout time, so a later address edit or
  catalog price change never rewrites history.
- **Shipping and payment are behind interfaces**
  (`pkg/shipping.Service`, `pkg/payment.Gateway`). The bundled shipping
  provider is a documented, formula-based mock (origin/destination/
  weight/courier all affect the quote) standing in for a real
  Indonesian courier aggregator — swap the constructor in `main.go` to
  go live. The payment gateway is a real Midtrans Snap + Core API
  implementation with proper SHA512 webhook signature verification;
  it needs real sandbox/production keys to actually call Midtrans.
- **Admin order status changes follow a small state machine**
  (`internal/services/order_service.go`) so an order can't jump from
  `DELIVERED` back to `PENDING_PAYMENT`, and cancelling a `PAID`/
  `PROCESSING`/`SHIPPED` order restores the stock that was deducted at
  payment time.
- **Reviews require a verified purchase.** `POST
  /api/v1/products/:id/reviews` takes an `order_id`; the service
  checks that exact order belongs to the requesting customer, actually
  contains the product, and has reached a paid-or-later status before
  the review is accepted (`OrderRepository.HasPurchasedProductInOrder`).
  One review per (user, product, order).

## What's not implemented

- Real shipping-provider and payment-provider credentials/integration
  testing (needs live sandbox keys and network access this environment
  didn't have).
- Automated tests. Given the scope of this build, hand-written tests
  weren't included — `go test ./...` scaffolding would be a reasonable
  next step, particularly around the checkout and webhook idempotency
  logic.
