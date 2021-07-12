FROM ubuntu:16.04

WORKDIR /root
COPY dplatformos dplatformos
COPY dplatformos-cli dplatformos-cli
COPY dplatformos.toml dplatformos*.toml ./
COPY entrypoint.sh entrypoint.sh

CMD ["/root/dplatformos", "-f", "/root/dplatformos.toml"]
