build:
	@cd cmd/config-manager-server; \
	go build -v

sync:
	govendor sync -v

generate:
	swagger generate server -P models.Customer --exclude-main --skip-validation -f swagger/swagger.yml -r LICENSE
