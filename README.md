# Order Service

This service manages orders.  
It creates orders, calls the payment service via gRPC, and updates order status.

## Tech stack
- Go
- Gin (HTTP API)
- gRPC
- PostgreSQL

## Run

```bash
go run cmd/app/main.go
