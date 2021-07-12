FROM ubuntu:16.04

WORKDIR /root
COPY dplatformos dplatformos
COPY dplatformos-cli dplatformos-cli
COPY dplatformos.toml ./

CMD ["/root/dplatformos", "-f", "/root/dplatformos.toml"]
