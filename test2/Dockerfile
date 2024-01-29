FROM oraclelinux:9

RUN set -ex; \
    groupdel input; \
    userdel systemd-coredump; \
    groupadd -g 1001 mysql; \
    useradd -u 1001 -r -g 1001 -s /sbin/nologin \
        -m -c "Default Application User" mysql

COPY mysqld.sh /usr/bin
COPY telemetry_agent.sh /usr/bin
COPY telemetry_supervisor.sh /usr/bin
COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
USER mysql
CMD ["/usr/bin/mysqld.sh"]
