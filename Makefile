.PHONY:
all: metering

metering: metering.go
	go build -o build/ metering.go

.PHONY:
kong: metering kong-restart

.PHONY:
kong-restart:
	docker-compose restart kong

.PHONY:
kong-up: metering
	docker-compose up -d

.PHONY:
kong-down:
	docker-compose down
	docker volume rm kong-plugin_kong_data

.PHONY:
kong-init: setup-mockbin setup-mike enable-key-auth setup-metering

.PHONY:
setup-mockbin:
	./scripts/admin.sh POST /services -d '{"name": "mockbin", "url": "http://mockbin.org"}'
	./scripts/admin.sh POST /services/mockbin/routes -d '{"name": "mock", "paths": ["/mock"]}'

.PHONY:
setup-mike:
	./scripts/admin.sh POST /consumers -d '{"username": "mike"}'
	./scripts/admin.sh POST /consumers/mike/key-auth -d '{"key": "super-secret-key"}'

.PHONY:
enable-key-auth:
	./scripts/admin.sh POST /plugins -d '{"name": "key-auth", "config": {"key_names": ["x-api-key"], "key_in_query": false}}'

.PHONY:
setup-metering:
	./scripts/admin.sh POST /plugins -d @metering.json

.PHONY:
test-with-mike:
	./scripts/test.sh GET /mock/requests -H 'x-api-key: super-secret-key' -i
