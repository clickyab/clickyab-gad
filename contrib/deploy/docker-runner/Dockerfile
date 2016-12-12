FROM registry.clickyab.ae/clickyab/ubuntu

RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y \
    supervisor \
    && rm -rf /var/lib/apt/lists/*

VOLUME /gad
