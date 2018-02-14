FROM drone/ca-certs

ENV GODEBUG=netdns=go
ENV XDG_CACHE_HOME /var/lib/autoscaler
ENV DRONE_DATABASE_PATH /var/lib/autoscaler/snapshot.db

ADD release/linux/arm64/drone-autoscaler /bin/

EXPOSE 8080 80 443

ENTRYPOINT ["/bin/drone-autoscaler"]
