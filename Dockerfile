FROM drone/ca-certs
EXPOSE 8080 80 443
VOLUME /data

ENV GODEBUG netdns=go
ENV XDG_CACHE_HOME /data
ENV DATABASE_DRIVER sqlite3
ENV DATABASE_DATASOURCE /data/database.sqlite?cache=shared&mode=rwc&_busy_timeout=9999999

ADD release/linux/amd64/drone-autoscaler /bin/
ENTRYPOINT ["/bin/drone-autoscaler"]
