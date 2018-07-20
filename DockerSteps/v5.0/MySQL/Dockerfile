FROM mysql

ENV MYSQL_ROOT_PASSWORD=maplewish \
    MYSQL_USER=mwish \
    MYSQL_PASSWORD=maplewish \
    MYSQL_ALLOW_EMPTY_PASSWORD=yes \
    MYSQL_DATABASE=go_test

ADD initdb.sql /docker-entrypoint-initdb.d

EXPOSE 3306

CMD ["./dbinit.sh"]