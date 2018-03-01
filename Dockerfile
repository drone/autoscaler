FROM drone/ca-certs

ENV GODEBUG=netdns=go
ENV XDG_CACHE_HOME /var/lib/autoscaler
ENV DATABASE_DRIVER=sqlite3
ENV DATABASE_DATASOURCE=/var/lib/autoscaler/database.sqlite?cache=shared&mode=rwc&_busy_timeout=9999999

ADD release/linux/arm64/drone-autoscaler /bin/

EXPOSE 8080 80 443

ENTRYPOINT ["/bin/drone-autoscaler"]
