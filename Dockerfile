# Step 1: Build the application
FROM golang:latest as builder

WORKDIR /app
COPY . .
RUN go build -o bridge .

# Step 2: Create a minimal runtime image
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/bridge .

# You can add any additional dependencies or configuration files here if needed

CMD ["./bridge"]