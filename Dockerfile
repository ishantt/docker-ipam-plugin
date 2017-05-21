FROM alpine
RUN apk update
RUN mkdir -p /run/docker/plugins
COPY docker-ipam-plugin docker-ipam-plugin
CMD ["/docker-ipam-plugin"]
