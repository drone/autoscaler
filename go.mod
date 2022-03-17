module github.com/drone/autoscaler

go 1.17

replace github.com/docker/docker => github.com/docker/engine v17.12.0-ce-rc1.0.20200309214505-aa6a9891b09c+incompatible

require (
	github.com/99designs/basicauth-go v0.0.0-20160802081356-2a93ba0f464d
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/aws/aws-sdk-go v1.13.5
	github.com/bluele/slack v0.0.0-20171128075526-307046097ee9
	github.com/dchest/uniuri v0.0.0-20160212164326-8902c56451e9
	github.com/digitalocean/godo v1.1.1
	github.com/docker/docker v0.0.0-00010101000000-000000000000
	github.com/docker/go-connections v0.3.0
	github.com/drone/drone-go v1.0.5-0.20190504210458-4d6116b897ba
	github.com/drone/envconfig v1.4.1
	github.com/drone/funcmap v0.0.0-20210903193859-704120d6923c
	github.com/drone/signal v0.0.0-20170915013802-ac5d07ef1315
	github.com/dustin/go-humanize v0.0.0-20171111073723-bb3d318650d4
	github.com/go-chi/chi v3.3.2+incompatible
	github.com/go-sql-driver/mysql v1.3.0
	github.com/golang/mock v1.3.1
	github.com/google/go-cmp v0.5.5
	github.com/gophercloud/gophercloud v0.0.0-20181014043407-c8947f7d1c51
	github.com/h2non/gock v1.0.7
	github.com/hetznercloud/hcloud-go v1.4.0
	github.com/jmoiron/sqlx v0.0.0-20180228184624-cf35089a1979
	github.com/joho/godotenv v1.2.0
	github.com/kr/pretty v0.1.0
	github.com/lib/pq v1.10.4
	github.com/mattn/go-sqlite3 v1.6.0
	github.com/packethost/packngo v0.1.0
	github.com/prometheus/client_golang v1.11.1
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.3
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.7.0
	github.com/yandex-cloud/go-genproto v0.0.0-20220314102905-1acaee8ca7eb
	github.com/yandex-cloud/go-sdk v0.0.0-20220314105123-d0c2a928feb6
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	google.golang.org/api v0.0.0-20180921000521-920bb1beccf7
)

require (
	cloud.google.com/go v0.34.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.7 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/containerd/containerd v1.3.4 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/distribution v0.0.0-20170726174610-edc3ab29cdff // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-ini/ini v1.32.0 // indirect
	github.com/gogo/protobuf v1.1.1 // indirect
	github.com/golang-jwt/jwt/v4 v4.1.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-querystring v0.0.0-20170111101155-53e6ce116135 // indirect
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.0.0-20160202185014-0b12d6b521d8 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/kr/text v0.1.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/nbio/st v0.0.0-20140626010706-e9e8d9816f32 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.26.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/tent/http-link-go v0.0.0-20130702225549-ac974c61c2f9 // indirect
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20210603081109-ebe580a85c40 // indirect
	golang.org/x/text v0.3.5 // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	google.golang.org/appengine v1.4.0 // indirect
	google.golang.org/genproto v0.0.0-20211021150943-2b146023228c // indirect
	google.golang.org/grpc v1.41.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/ini.v1 v1.51.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
	gotest.tools v2.2.0+incompatible // indirect
	launchpad.net/gocheck v0.0.0-20140225173054-000000000087 // indirect
)

replace github.com/drone/funcmap => github.com/iainlane/funcmap v0.0.0-20211116113722-13f662008062
