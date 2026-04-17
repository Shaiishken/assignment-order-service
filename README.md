# Order Service

Сервис для работы с заказами.  
Создает заказ, отправляет запрос в payment-service через gRPC и обновляет статус.

## Что используется
- Go
- Gin (HTTP)
- gRPC
- PostgreSQL

## Запуск

```bash
go run cmd/app/main.go
