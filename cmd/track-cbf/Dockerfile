FROM golang:1.20.3-alpine AS builder

COPY . /src

WORKDIR /src

RUN go build -o /app/track-cbf ./cmd/track-cbf

FROM gcr.io/distroless/static-debian11:nonroot

COPY --from=builder /app /app

WORKDIR /app

ENTRYPOINT [ "/app/track-cbf" ]
