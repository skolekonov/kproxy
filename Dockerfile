FROM alpine
RUN apk update && apk add ca-certificates
COPY kproxy /
CMD ["/kproxy"]
