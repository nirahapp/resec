# Build layer
FROM golang:1.21 AS builder
WORKDIR /go/src/github.com/nirahapp/resec
COPY . .
ARG RESEC_VERSION
ENV RESEC_VERSION ${RESEC_VERSION:-local-dev}
RUN echo $RESEC_VERSION
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'main.Version=${RESEC_VERSION}'" -a -installsuffix cgo -o build/resec  .

# Run layer
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/nirahapp/resec/build/resec .
CMD ["./resec"]
