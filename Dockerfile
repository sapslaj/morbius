FROM alpine:3
WORKDIR /
RUN apk add gcompat
COPY morbius /bin/morbius
COPY minimal-config.yaml /etc/morbius/config.yaml
ENV MORBIUS_CONFIG_FILE=/etc/morbius/config.yaml
EXPOSE 2055/udp 2056/udp 6343/udp 6060/tcp
CMD ["/bin/morbius"]
