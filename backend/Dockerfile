# Construye la API o el worker según el argumento CMD (api | worker).
#   docker build --build-arg CMD=api -t saas-api backend
#   docker build --build-arg CMD=worker -t saas-worker backend
FROM golang:1.25 AS build
ARG CMD=api
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/${CMD}

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/app /app
USER nonroot:nonroot
ENTRYPOINT ["/app"]
