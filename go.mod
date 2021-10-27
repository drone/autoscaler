module github.com/drone/autoscaler

go 1.12

replace github.com/docker/docker => github.com/docker/engine v17.12.0-ce-rc1.0.20200309214505-aa6a9891b09c+incompatible

require (
	cloud.google.com/go v0.28.0 // indirect
	github.com/99designs/basicauth-go v0.0.0-20160802081356-2a93ba0f464d
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/Microsoft/go-winio v0.4.7 // indirect
	github.com/aws/aws-sdk-go v1.13.5
	github.com/beorn7/perks v0.0.0-20160804104726-4c0e84591b9a // indirect
	github.com/bluele/slack v0.0.0-20171128075526-307046097ee9
	github.com/containerd/containerd v1.3.4 // indirect
	github.com/dchest/uniuri v0.0.0-20160212164326-8902c56451e9
	github.com/digitalocean/godo v1.1.1
	github.com/docker/distribution v0.0.0-20170726174610-edc3ab29cdff // indirect
	github.com/docker/docker v0.0.0-00010101000000-000000000000
	github.com/docker/go-connections v0.3.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/drone/drone-go v1.0.5-0.20190504210458-4d6116b897ba
	github.com/drone/envconfig v1.4.1
	github.com/drone/funcmap v0.0.0-20210903193859-704120d6923c
	github.com/drone/signal v0.0.0-20170915013802-ac5d07ef1315
	github.com/dustin/go-humanize v0.0.0-20171111073723-bb3d318650d4
	github.com/go-chi/chi v3.3.2+incompatible
	github.com/go-ini/ini v1.32.0 // indirect
	github.com/go-sql-driver/mysql v1.3.0
	github.com/gogo/protobuf v0.0.0-20170307180453-100ba4e88506 // indirect
	github.com/golang/mock v1.3.1
	github.com/google/go-cmp v0.4.0
	github.com/google/go-querystring v0.0.0-20170111101155-53e6ce116135 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gophercloud/gophercloud v0.0.0-20181014043407-c8947f7d1c51
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/h2non/gock v1.0.7
	github.com/hetznercloud/hcloud-go v1.4.0
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jmespath/go-jmespath v0.0.0-20160202185014-0b12d6b521d8 // indirect
	github.com/jmoiron/sqlx v0.0.0-20180228184624-cf35089a1979
	github.com/joho/godotenv v1.2.0
	github.com/kr/pretty v0.0.0-20160823170715-cfb55aafdaf3
	github.com/kr/text v0.0.0-20160504234017-7cafcd837844 // indirect
	github.com/lib/pq v1.0.0
	github.com/mattn/go-sqlite3 v1.6.0
	github.com/matttproud/golang_protobuf_extensions v1.0.0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/nbio/st v0.0.0-20140626010706-e9e8d9816f32 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/packethost/packngo v0.1.0
	github.com/pkg/errors v0.8.1 // indirect
	github.com/prometheus/client_golang v0.8.0
	github.com/prometheus/common v0.0.0-20180110214958-89604d197083 // indirect
	github.com/prometheus/procfs v0.0.0-20180212145926-282c8707aa21 // indirect
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.3
	github.com/sirupsen/logrus v1.4.2
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/tent/http-link-go v0.0.0-20130702225549-ac974c61c2f9 // indirect
	golang.org/x/crypto v0.0.0-20190621222207-cc06ce4a13d4
	golang.org/x/lint v0.0.0-20190313153728-d0100b6bd8b3 // indirect
	golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	golang.org/x/tools v0.0.0-20190524140312-2c0ae7006135 // indirect
	google.golang.org/api v0.0.0-20180921000521-920bb1beccf7
	google.golang.org/appengine v1.4.0 // indirect
	google.golang.org/grpc v1.30.0 // indirect
	gopkg.in/ini.v1 v1.51.0 // indirect
	gotest.tools v2.2.0+incompatible // indirect
	honnef.co/go/tools v0.0.0-20190523083050-ea95bdfd59fc // indirect
)

replace github.com/drone/funcmap => github.com/iainlane/funcmap v0.0.0-20211116113722-13f662008062
