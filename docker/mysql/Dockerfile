FROM mysql:5.7.25

EXPOSE 3306
COPY ./docker/mysql/init/* /docker-entrypoint-initdb.d/
COPY ./docker/mysql/my.cnf /etc/mysql/conf.d/my.cnf

CMD ["mysqld"]