FROM golang:1.26-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /ketide .

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /ketide ./ketide
COPY templates ./templates
COPY static ./static
EXPOSE 1324
ENTRYPOINT ["/app/ketide"]
