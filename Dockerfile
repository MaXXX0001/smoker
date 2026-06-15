# Спільний multi-stage Dockerfile для всіх сервісів.
# Який бінар збирати — задає build-arg SERVICE (cosmos|chronos|chaos|orchestrator|gateway).
ARG SERVICE

FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG SERVICE
# GOWORK=off — у контейнері збираємо суто за go.mod, без workspace.
RUN CGO_ENABLED=0 GOWORK=off go build -trimpath -o /out/app ./cmd/${SERVICE}

FROM alpine:3.20
# tzdata — для time.LoadLocation (робочі години/локальний час),
# ca-certificates — для HTTPS-викликів зовнішніх API.
RUN apk add --no-cache tzdata ca-certificates
COPY --from=build /out/app /app
ENTRYPOINT ["/app"]
