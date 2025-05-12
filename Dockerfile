# Stage 1: Build the application
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache nodejs npm

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Install dependencies and generate output.css
RUN npm install  && \
    npx @tailwindcss/cli -i ./cmd/web/assets/css/input.css -o ./cmd/web/assets/css/output.css

# Install templ and generate templates
RUN go install github.com/a-h/templ/cmd/templ@v0.2.793 && \
    templ generate

# Build the application
RUN go build -o main ./cmd/api/main.go

# Stage 2: Final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the built binary and .env file
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8082

ENTRYPOINT ["./main"]