module github.com/MrColorado/backend/server

go 1.22.0

replace github.com/MrColorado/backend/internal => ../internal

replace github.com/MrColorado/backend/logger => ../logger

require (
	github.com/MrColorado/backend/internal v0.0.0-00010101000000-000000000000
	github.com/MrColorado/backend/logger v0.0.0-00010101000000-000000000000
	github.com/aws/aws-sdk-go-v2 v1.25.3
	github.com/aws/aws-sdk-go-v2/config v1.27.7
	github.com/aws/aws-sdk-go-v2/credentials v1.17.7
	github.com/aws/aws-sdk-go-v2/service/s3 v1.51.4
	github.com/friendsofgo/errors v0.9.2
	github.com/lib/pq v1.10.9
	github.com/nats-io/nats.go v1.33.1
	github.com/volatiletech/sqlboiler/v4 v4.16.2
	github.com/volatiletech/strmangle v0.0.6
	google.golang.org/grpc v1.62.1
	google.golang.org/protobuf v1.33.0
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.1 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.15.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.3.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.17.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.20.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.23.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.28.4 // indirect
	github.com/aws/smithy-go v1.20.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/klauspost/compress v1.17.7 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/volatiletech/inflect v0.0.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240311173647-c811ad7063a7 // indirect
	gopkg.in/errgo.v2 v2.1.0 // indirect
)
