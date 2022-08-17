module github.com/drone/autoscaler

go 1.19

replace github.com/docker/docker => github.com/docker/engine v17.12.0-ce-rc1.0.20200309214505-aa6a9891b09c+incompatible

require (
	github.com/99designs/basicauth-go v0.0.0-20160802081356-2a93ba0f464d
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/aws/aws-sdk-go v1.44.186
	github.com/bluele/slack v0.0.0-20171128075526-307046097ee9
	github.com/dchest/uniuri v0.0.0-20160212164326-8902c56451e9
	github.com/digitalocean/godo v1.1.1
	github.com/docker/docker v0.0.0-00010101000000-000000000000
	github.com/docker/go-connections v0.3.0
	github.com/drone/drone-go v1.0.5-0.20190504210458-4d6116b897ba
	github.com/drone/envconfig v1.4.1
	github.com/drone/funcmap v0.0.0-20220929084810-72602997d16f
	github.com/drone/signal v0.0.0-20170915013802-ac5d07ef1315
	github.com/dustin/go-humanize v0.0.0-20171111073723-bb3d318650d4
	github.com/go-chi/chi v3.3.2+incompatible
	github.com/go-sql-driver/mysql v1.3.0
	github.com/golang/mock v1.6.0
	github.com/google/go-cmp v0.5.8
	github.com/gophercloud/gophercloud v0.0.0-20181014043407-c8947f7d1c51
	github.com/h2non/gock v1.0.7
	github.com/hetznercloud/hcloud-go v1.4.0
	github.com/jmoiron/sqlx v0.0.0-20180228184624-cf35089a1979
	github.com/joho/godotenv v1.2.0
	github.com/kr/pretty v0.1.0
	github.com/lib/pq v1.10.4
	github.com/mattn/go-sqlite3 v1.6.0
	github.com/packethost/packngo v0.1.0
	github.com/prometheus/client_golang v1.14.0
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.3
	github.com/sirupsen/logrus v1.6.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/oauth2 v0.0.0-20220622183110-fd043fe589d2
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4
	golang.org/x/time v0.0.0-20220224211638-0e9765cccd65
	google.golang.org/api v0.93.0
)

require (
	cloud.google.com/go/compute v1.7.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.7 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/containerd/containerd v1.3.4 // indirect
	github.com/docker/distribution v0.0.0-20170726174610-edc3ab29cdff // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/gogo/protobuf v1.1.1 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-querystring v0.0.0-20170111101155-53e6ce116135 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.1.0 // indirect
	github.com/googleapis/gax-go/v2 v2.4.0 // indirect
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/kr/text v0.1.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/nbio/st v0.0.0-20140626010706-e9e8d9816f32 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/tent/http-link-go v0.0.0-20130702225549-ac974c61c2f9 // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/net v0.1.0 // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220624142145-8cd45d7dbd1f // indirect
	google.golang.org/grpc v1.47.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gotest.tools v2.2.0+incompatible // indirect
	launchpad.net/gocheck v0.0.0-20140225173054-000000000087 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)
