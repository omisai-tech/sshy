FROM alpine:3.19

RUN apk add --no-cache openssh-client git ca-certificates

COPY sshy /usr/local/bin/sshy

RUN chmod +x /usr/local/bin/sshy

ENTRYPOINT ["sshy"]
