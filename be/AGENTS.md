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
    - [Graceful Shutdown](#graceful-shutdown)
5. [Configuration](#configuration)
6. [Billing Logic](#billing-logic)
7. [Development & Operations](#development--operations)
8. [Extending the Codebase](#extending-the-codebase)
9. [References](#references)

---

## Overview

The Weather Boy backend is a Go-based service designed for reliability, observability, and maintainability. It leverages the Fiber web framework for HTTP APIs, PostgreSQL for persistence, and OpenTelemetry for distributed tracing and metrics.

---

## Project Structure

- `cmd/weatherboyapi/` – Application entry point (`main.go`, `init.go`)
- `internal/`
    - `constants/` – Service-wide constants
    - `db/` – Database connection and utilities
    - `handlers/` – HTTP route handlers
    - `healthcheck/` – Health monitoring logic
    - `logger/` – Logging setup
    - `opentelemetry/` – OTEL setup and utilities
    - `otelware/` – Fiber middleware for tracing/metrics
    - `repository/` – (Placeholder for data access logic)
    - `router/` – HTTP server and route registration
    - `shutdown/` – Graceful shutdown logic
    - `types/` – Shared types
    - `utils/` – Utility functions
- `Makefile` – Build, run, test, and clean commands
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

### Graceful Shutdown

- **Handler:** `internal/shutdown/handleshutdown.go`
    - Listens for SIGTERM/SIGINT.
    - Marks service unhealthy, waits for draining.
    - Shuts down HTTP server, waits for goroutines, closes DB and OTEL providers.

---

## Configuration

All sensitive and environment-specific configuration is managed via environment variables, typically loaded from a `.env` file. Key variables include:

- `POSTGRES_HOST`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_PORT`
- `OTELCOL_HOST`, `OTELCOL_PORT`
- `LISTEN_ADDR`


## Development & Operations

- **Build:** `make build`
- **Run:** `make run` (loads `.env`, builds, and runs the binary)
- **Test:** `make test`
- **Clean:** `make clean`
- **Docker:** Use the provided `Dockerfile` for containerization.

---

## Extending the Codebase

- **Add New Endpoints:** Implement handler functions in `internal/handlers/`, register them in `internal/router/router.go`.
- **Add Database Logic:** Extend `internal/db/` and/or create repositories in `internal/repository/`.
- **Add Metrics/Tracing:** Use OpenTelemetry APIs or extend `otelware` middleware.
- **Add Health Checks:** Extend `internal/healthcheck/healthcheck.go` as needed.

---

## References

- [Fiber Documentation](https://docs.gofiber.io/)
- [pgx PostgreSQL Driver](https://github.com/jackc/pgx)
- [OpenTelemetry for Go](https://opentelemetry.io/docs/instrumentation/go/)
- [Project Billing Algorithm](docs/billing_algorithm.md)

---

**For questions or contributions, please follow the project's code style and open a PR or issue as appropriate. Happy hacking!**
