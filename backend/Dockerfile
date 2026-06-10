FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o /backend ./cmd/main.go
RUN go build -o /migrator ./cmd/migrations/main.go

FROM alpine:3.20

RUN adduser -D appuser
RUN apk --no-cache add ca-certificates curl
RUN mkdir -p /app/uploads/images && chown -R appuser:appuser /app/uploads

COPY --from=build /backend /app/backend
COPY --from=build /migrator /app/migrator
COPY --from=build /src/migrations /app/migrations

USER appuser
WORKDIR /app

CMD ["./backend"]
