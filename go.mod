module github.com/WaffleHacks/mailer

go 1.17

require (
	github.com/aws/aws-sdk-go-v2 v1.16.2
	github.com/aws/aws-sdk-go-v2/config v1.14.0
	github.com/aws/aws-sdk-go-v2/credentials v1.9.0
	github.com/aws/aws-sdk-go-v2/service/sesv2 v1.12.0
	github.com/cenkalti/backoff/v4 v4.1.3
	github.com/getsentry/sentry-go v0.12.0
	github.com/go-chi/chi/v5 v5.0.7
	github.com/gobwas/glob v0.2.3
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.3
	github.com/joho/godotenv v1.4.0
	github.com/k3a/html2text v1.0.8
	github.com/mailgun/mailgun-go/v4 v4.6.0
	github.com/matoous/go-nanoid v1.5.0
	go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws v0.32.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.32.0
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/exporters/jaeger v1.7.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.7.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.7.0
	go.opentelemetry.io/otel/sdk v1.7.0
	go.opentelemetry.io/otel/trace v1.7.0
	go.uber.org/zap v1.21.0
	google.golang.org/genproto v0.0.0-20220222213610-43724f9ea8cf
	google.golang.org/grpc v1.46.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)

require (
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.11.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.15.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.7.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.10.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.15.0 // indirect
	github.com/aws/smithy-go v1.11.2 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.7.0 // indirect
	go.opentelemetry.io/proto/otlp v0.16.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/net v0.0.0-20211008194852-3b03d305991f // indirect
	golang.org/x/sys v0.0.0-20211007075335-d3039528d8ac // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
)
