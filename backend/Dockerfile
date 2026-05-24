FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG CMD=api
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out ./cmd/${CMD}

FROM gcr.io/distroless/static-debian12
COPY --from=builder /out /app
EXPOSE 8080
ENTRYPOINT ["/app"]
