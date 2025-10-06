FROM hub.byted.org/infcplibrary/golang:1.20.5-bullseye-auth AS builder
WORKDIR /go/src/code.byted.org/epscp/vetes-api
COPY . .
RUN make build-linux

FROM hub.byted.org/base/security.debian11:latest
RUN apt update \
     && apt install -y --no-install-recommends ca-certificates curl locales \
     && apt clean \
     && rm -rf /var/lib/apt/lists/*
ENV LANG=C.UTF-8
RUN locale-gen
USER nobody
WORKDIR /app
COPY --from=builder /go/src/code.byted.org/epscp/vetes-api/bin/vetes-api .
ENTRYPOINT ["/app/vetes-api"]
