FROM debian:latest

COPY . /app

WORKDIR /home

RUN apt-get update \
  && apt-get install -y sqlite3 \
  && ln -s /app/static \
  && /app/dbinit.sh

CMD UPD_PORTS=1961,1962,1963 /app/udptest