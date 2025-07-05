# AGENTS.md

## Weather Boy Backend – Developer Guide

Welcome to the backend codebase for **Weather Boy**! This document is designed to help developers and agents quickly understand the architecture, core components, and operational flow of the backend service.

---

## Table of Contents

1. [Overview](#overview)
2. [Project Structure](#project-structure)
3. [Startup Flow](#startup-flow)
4. [Core Components](#core-components)
    - [HTTP Server & Routing](#http-server--routing)
    - [Health Checking](#health-checking)
    - [Database Layer](#database-layer)
    - [Observability & Telemetry](#observability--telemetry)
    - [Scheduler](#scheduler)
    - [Risk Scoring](#risk-scoring)
    - [Graceful Shutdown](#graceful-shutdown)
5. [Configuration](#configuration)
6. [Development & Operations](#development--operations)
7. [Extending the Codebase](#extending-the-codebase)
8. [References](#references)

---

## Overview

The Weather Boy backend is a Go service that periodically pulls multiple IMD data sources (bulletin PDFs, radar snapshots and nowcast JSON), stores them in PostgreSQL and exposes them via a small Fiber API. A risk score is calculated from these feeds and surfaced at `/v1/risk/:loc`. The project emphasizes observability via OpenTelemetry and structured logging using zap.

---

## Project Structure

- `cmd/weatherboyapi/` – Application entry point (`main.go`, `init.go`)
- `internal/`
    - `config/` – Environment loading and location list
    - `constants/` – Service-wide constants
    - `db/` – Database connection and utilities
    - `fetch/` – External data fetchers (nowcast, district warning, radar, river basin, AWS ARG, bulletin PDFs, etc.)
    - `parse/` – Parsing / AI summarisation logic (bulletin, radar, etc.)
    - `handlers/` – HTTP route handlers
    - `healthcheck/` – Health monitoring logic
    - `logger/` – Zap-based structured logging
    - `model/` – Database models
    - `opentelemetry/` – OTEL setup and utilities
    - `otelware/` – Fiber middleware for tracing/metrics
    - `repository/` – Database access helpers
    - `router/` – HTTP server and route registration
    - `scheduler/` – Cron jobs for fetching and parsing data
    - `score/` – Risk scoring logic
    - `shutdown/` – Graceful shutdown logic
    - `types/` – Shared types
    - `utils/` – Utility functions
- `migrations/` – SQL migrations managed by `goose`
- `Makefile` – Build, run, migrate, dev, and clean commands
- `Dockerfile` – Containerization

---

## Startup Flow

1. **Initialization (`init.go`):**
    - Sets up OpenTelemetry (tracing/metrics).
    - Initializes the PostgreSQL connection pool using environment variables.

2. **Main Application (`main.go`):**
    - Logs startup.
    - Starts the healthcheck goroutine.
    - Launches the HTTP server.
    - Waits for graceful shutdown on termination signals.

---

## Core Components

### HTTP Server & Routing

- **Framework:** [Fiber](https://gofiber.io/)
- **Server Setup:** `internal/router/startserver.go`
    - Configures middleware: recovery, CORS, OpenTelemetry.
    - Registers routes via `RegisterRoutes()`.
    - Listens on the address specified by `LISTEN_ADDR` env variable.
- **Routes:** Defined in `internal/router/router.go`
    - Example: `/health` endpoint for health checks.

### Health Checking

- **Logic:** `internal/healthcheck/healthcheck.go`
    - Periodically checks database and OTEL connectivity.
    - Updates a global `IsHealthy` flag.
- **Endpoint:** `/health` (see `internal/handlers/health.go`)
    - Returns `200 OK` if healthy, `418 Teapot` if shutting down.

### Database Layer

- **Driver:** [pgxpool](https://github.com/jackc/pgx)
- **Initialization:** `internal/db/db.go`
    - Parses DSN from environment.
    - Sets up connection pool and metrics.
    - Provides health check via `Ping`.
- **Access:** Use `db.GetDBDriver()` to access the pool.

### Observability & Telemetry

- **OpenTelemetry:** `internal/opentelemetry/opentelemetry.go`
    - Connects to OTEL collector via gRPC.
    - Sets up MeterProvider and TracerProvider.
- **Middleware:** `internal/otelware/`
    - Adds tracing and metrics to all HTTP requests.
    - Customizable via options.

### Scheduler

- **Cron Jobs:** `internal/scheduler/jobs.go`
    - Fetches bulletin PDFs daily at 18:30 IST.
    - Fetches nowcast JSON every 15 minutes with ±30 s jitter.
    - Jobs run in the Asia/Kolkata timezone.

### IMD Nowcast

- **Fetcher:** `internal/fetch/imdnowcast.go`
    - Builds `https://mausam.imd.gov.in/api/nowcast_district_api.php?id={DistrictID}` from static base URL.
    - Stores the raw JSON payload in `nowcast_raw` for auditing.
    - Extracts the `color` field and maps it to a POP value via `colorToPOP()`.
    - Persists `cat1..cat19` flags in `nowcast_category` linked to each nowcast row.

### District & State Warnings

- **Fetcher:** `internal/fetch/district_warning.go`
    - Hits `https://mausam.imd.gov.in/api/warnings_district_api.php?id={DistrictID}`.
    - Stores raw JSON in `district_warning_raw`.
    - Parses day-1–5 warning text and color codes into `district_warning`.

### Doppler Radar

- **Fetcher:** `internal/fetch/radar.go`
    - Downloads the latest Doppler radar composite PNG for the target station.
- **Parser:** `internal/parse/radar.go`
    - Converts IMD colour palette to dBZ.
    - Calculates max dBZ within 40 km radius for each location.
    - Persists to `radar` table (`location`, `max_dbz`, `captured_at`).

### Bulletin PDF

- **Fetcher:** `internal/fetch/bulletin.go`
    - Downloads the daily district bulletin PDF.
- **Parser / AI Summariser:** `internal/parse/bulletin.go`
    - Extracts text via PDF-to-text, summarises with Gemini Flash, stores structured forecast in `bulletin_parsed`.

### River Basin QPF

- **Fetcher:** `internal/fetch/riverbasin.go`
    - Retrieves JSON QPF for major river basins.
    - Persists to `river_basin_qpf` table.

### AWS / ARG Observations

- **Fetcher:** `internal/fetch/awsarg.go`
    - Ingests automatic weather-station / rain-gauge 1-hour precipitation totals.
    - Persists to `aws_arg` table.

### Risk Scoring

- **Logic:** `internal/score/score.go`
    - Bulletin heavy/very-heavy mention → +0.4
    - Radar ≥45 dBZ within 40 km → +0.4
    - Nowcast POP₁ₕ ≥0.7 → +0.2
    - Certain `catN` flags (2 and 3) → +0.1 each
    - Thresholds: ≥0.8 = RED, ≥0.5 = ORANGE, ≥0.3 = YELLOW, else GREEN.

### Graceful Shutdown

- **Handler:** `internal/shutdown/handleshutdown.go`
    - Listens for SIGTERM/SIGINT.
    - Marks service unhealthy, waits for draining.
    - Shuts down HTTP server, waits for goroutines, closes DB and OTEL providers.

---

## Configuration

All sensitive and environment-specific configuration is managed via environment variables and is automatically loaded from a `.env` file by the `config` package. Key variables include:

- `POSTGRES_HOST`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_PORT`
- `OTELCOL_HOST`, `OTELCOL_PORT`
- `LISTEN_ADDR`
- `OPENAI_API_KEY`, `DATA_DIR`, `METNET_BASE_URL`


## Development & Operations

- **Build:** `make build`
- **Run:** `make run` (loads `.env`, builds, and runs the binary)
- **Dev:** `make dev` (docker compose with hot reload via Air)
- **Migrate:** `make migrate` to apply DB schema changes
- **Test:** `go test ./...`
- **Clean:** `make clean`
- **Docker:** Use the provided `Dockerfile` for containerization.

---

## Extending the Codebase

- **Add New Endpoints:** Implement handler functions in `internal/handlers/`, register them in `internal/router/router.go`.
- **Add Database Logic:** Extend `internal/db/` and/or create repositories in `internal/repository/`.
- **Add Metrics/Tracing:** Use OpenTelemetry APIs or extend `otelware` middleware.
- **Add Health Checks:** Extend `internal/healthcheck/healthcheck.go` as needed.
- **Add Fetch Jobs:** Modify `internal/scheduler/jobs.go` and add code under `internal/fetch/`.
- **Tweak Scoring:** Update logic in `internal/score/score.go` as the risk model evolves.

---

## References

- [Fiber Documentation](https://docs.gofiber.io/)
- [pgx PostgreSQL Driver](https://github.com/jackc/pgx)
- [OpenTelemetry for Go](https://opentelemetry.io/docs/instrumentation/go/)

---

**For questions or contributions, please follow the project's code style and open a PR or issue as appropriate. Happy hacking!**
