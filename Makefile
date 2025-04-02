build-alpine:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s"   -o ./bin/binary ./cmd/app
swagger:
	swag init -g cmd/app/main.go --parseDependency --parseInternal