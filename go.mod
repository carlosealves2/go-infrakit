module github.com/carlosealves2/go-infrakit

go 1.24.2

require (
	github.com/phuslu/log v0.0.0
	github.com/redis/go-redis/v9 v9.0.0
	go.opentelemetry.io/otel v0.0.0
)

replace github.com/phuslu/log => ./third_party/phuslu/log

replace github.com/redis/go-redis/v9 => ./third_party/redis

replace go.opentelemetry.io/otel => ./third_party/otel
