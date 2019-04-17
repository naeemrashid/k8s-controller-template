FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ADD controller-template .
CMD ["./controller-template"]
