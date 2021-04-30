module github.com/drone/autoscaler

go 1.12

replace github.com/docker/docker => github.com/docker/engine v17.12.0-ce-rc1.0.20200309214505-aa6a9891b09c+incompatible

require (
	github.com/99designs/basicauth-go v0.0.0-20160802081356-2a93ba0f464d
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.7 // indirect
	github.com/aws/aws-sdk-go v1.27.0
	github.com/bluele/slack v0.0.0-20171128075526-307046097ee9
	github.com/containerd/containerd v1.3.4 // indirect
	github.com/dchest/uniuri v0.0.0-20160212164326-8902c56451e9
	github.com/digitalocean/godo v1.60.0
	github.com/docker/distribution v0.0.0-20170726174610-edc3ab29cdff // indirect
	github.com/docker/docker v0.0.0-00010101000000-000000000000
	github.com/docker/go-connections v0.3.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/drone/drone-go v1.0.5-0.20190504210458-4d6116b897ba
	github.com/drone/envconfig v1.4.1
	github.com/drone/signal v0.0.0-20170915013802-ac5d07ef1315
	github.com/dustin/go-humanize v0.0.0-20171111073723-bb3d318650d4
	github.com/go-chi/chi v3.3.2+incompatible
	github.com/go-sql-driver/mysql v1.4.0
	github.com/golang/mock v1.5.0
	github.com/google/go-cmp v0.5.5
	github.com/gophercloud/gophercloud v0.0.0-20181014043407-c8947f7d1c51
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/h2non/gock v1.0.7
	github.com/hetznercloud/hcloud-go v1.4.0
	github.com/jmoiron/sqlx v0.0.0-20180228184624-cf35089a1979
	github.com/joho/godotenv v1.2.0
	github.com/kr/pretty v0.1.0
	github.com/lib/pq v1.0.0
	github.com/mattn/go-sqlite3 v1.6.0
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/nbio/st v0.0.0-20140626010706-e9e8d9816f32 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/packethost/packngo v0.1.0
	github.com/prometheus/client_golang v1.10.0
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.3
	github.com/sirupsen/logrus v1.6.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/oauth2 v0.0.0-20210413134643-5e61552d6c78
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	google.golang.org/api v0.45.0
	gotest.tools v2.2.0+incompatible // indirect
)
