FROM registry.clickyab.ae/clickyab/ubuntu

RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y \
    supervisor \
    && rm -rf /var/lib/apt/lists/*

ADD init-script.sh /sbin/init-script.sh
RUN chmod a+x /sbin/init-script.sh

ENTRYPOINT [/sbin/init-script.sh]
CMD server

VOLUME /gad
