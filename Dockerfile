FROM alpine:latest as certificates
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# RUNTIME

FROM scratch
COPY --from=certificates /etc/ssl/certs /etc/ssl/certs
ADD bin/rbl-control /
CMD ["/rbl-control"]