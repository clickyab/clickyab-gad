FROM registry.clickyab.ae/clickyab/ubuntu

RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y \
    wget \
    && wget https://github.com/Yelp/dumb-init/releases/download/v1.2.0/dumb-init_1.2.0_amd64.deb \
    && dpkg -i dumb-init_*.deb \
    && mkdir /app \
    && rm dumb-init_*.deb \
    && rm -rf /var/lib/apt/lists/*

VOLUME /app

ADD init-script.sh /sbin/init-script.sh
RUN chmod a+x /sbin/init-script.sh

WORKDIR /app
ENTRYPOINT ["/sbin/init-script.sh"]

CMD ["/bin/bash"]

