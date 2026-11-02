FROM node:22-alpine AS web
WORKDIR /web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM golang:1.27-bookworm AS api
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=web /web/dist ./web/dist
RUN CGO_ENABLED=0 go build -o /out/server ./cmd/server

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=api /out/server ./server
COPY --from=web /web/dist ./web/dist
ENV WEB_ROOT=/app/web/dist
ENV HTTP_ADDR=:8080
EXPOSE 8080
CMD ["./server"]
