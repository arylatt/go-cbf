FROM golang:1.20.3-alpine AS builder

COPY . /src

WORKDIR /src

RUN go build -o /app/track-cbf ./cmd/track-cbf

RUN mkdir -p /data

FROM gcr.io/distroless/static-debian11:nonroot

COPY --from=builder /app /app

WORKDIR /app

COPY --from=builder --chown=nonroot:nonroot /data /data

VOLUME /data

ENTRYPOINT [ "/app/track-cbf", "--output-dir", "/data" ]
