module go-gin-ddd

go 1.17

require (
	cloud.google.com/go/storage v1.18.0
	github.com/gin-contrib/sessions v0.0.4
	github.com/gin-gonic/gin v1.7.7
	github.com/google/uuid v1.3.0
	github.com/joho/godotenv v1.4.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/ken109/gin-jwt v1.1.4
	github.com/noknow-hub/go_crypto v0.0.0-20190921184517-be89dab92b85
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292
	golang.org/x/oauth2 v0.0.0-20211005180243-6b3c2da341f1
	google.golang.org/api v0.58.0
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gorm.io/driver/mysql v1.3.2
	gorm.io/gorm v1.23.2
	jaytaylor.com/html2text v0.0.0-20211105163654-bc68cce691ba
	packages v0.0.0
)

require (
	cloud.google.com/go v0.97.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/gin-contrib/cors v1.3.1 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.10.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/goccy/go-json v0.9.4 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/googleapis/gax-go/v2 v2.1.1 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/gorilla/sessions v1.2.1 // indirect
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.8 // indirect
	github.com/lestrrat-go/blackmagic v1.0.0 // indirect
	github.com/lestrrat-go/httpcc v1.0.0 // indirect
	github.com/lestrrat-go/iter v1.0.1 // indirect
	github.com/lestrrat-go/jwx v1.2.19 // indirect
	github.com/lestrrat-go/option v1.0.0 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/ssor/bom v0.0.0-20170718123548-6386211fdfcf // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f // indirect
	golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20210924002016-3dee208752a0 // indirect
	google.golang.org/grpc v1.40.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace (
	github.com/go-playground/validator/v10 => github.com/ken109/validator/v10 v10.9.1-0.20211113103440-ec8ebbbebe05
	packages => ./packages
)
