FROM golang:1.23 AS build
WORKDIR /src
COPY . .
RUN go build -trimpath -o /out/harbor-sfu ./cmd/harbor-sfu
FROM gcr.io/distroless/base-debian12
COPY --from=build /out/harbor-sfu /harbor-sfu
EXPOSE 8093/tcp 10000/udp
ENTRYPOINT ["/harbor-sfu"]
