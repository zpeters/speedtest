FROM iron/base
WORKDIR /app
COPY bin/speedtest /app/
ENTRYPOINT ["./speedtest"]
