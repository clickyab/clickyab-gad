FROM registry.clickyab.ae/clickyab/ubuntu

RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y \
    wget php-fpm php-mysql php php-memcached php-memcache php-redis php-gd  php-cli php-json \
    && wget https://github.com/Yelp/dumb-init/releases/download/v1.2.0/dumb-init_1.2.0_amd64.deb \
    && dpkg -i dumb-init_*.deb \
    && rm dumb-init_*.deb \
    && mkdir /app \
    && rm -rf /var/lib/apt/lists/*
    && sed -i '/daemonize /c daemonize = no' /etc/php/7.0/fpm/php-fpm.conf
    && sed -i 's/memory_limit = 128M/memory_limit = 256M/g' /etc/php/7.0/fpm/php.ini
    && sed -i 's/;opcache.enable=0/opcache.enable=1/g' /etc/php/7.0/fpm/php.ini
    && rm /etc/php/7.0/fpm/pool.d/www.conf

VOLUME /app

ADD init-script.sh /sbin/init-script.sh
RUN chmod a+x /sbin/init-script.sh
ADD www.conf /etc/php/7.0/fpm/pool.d/www.conf

WORKDIR /app
ENTRYPOINT ["/sbin/init-script.sh"]

CMD ["/bin/bash"]

