DOCKER=docker

PROJECT_NAME=gotasker

test:
	go test ./...

setup-redis:
	@if ! $(DOCKER) ps | /bin/grep ${PROJECT_NAME}-redis-local; then \
		$(DOCKER) run --name ${PROJECT_NAME}-redis-local \
			-p 6379:6379 \
			-v ${PROJECT_NAME}_data:/data \
			--restart always \
			-d redis:alpine;\
	fi

setup-swagger:
	@if ! $(DOCKER) ps | /bin/grep ${PROJECT_NAME}-swagger-local; then \
		$(DOCKER) run --name ${PROJECT_NAME}-swagger-local \
			-e PORT=9527 \
			-p 9527:9527 \
			-v ./doc/openapi:/config \
			-e SWAGGER_JSON=/config/api.yaml \
			--restart always \
			-d swaggerapi/swagger-ui:latest;\
	fi

remove:
	$(DOCKER) rm -f ${PROJECT_NAME}-redis-local ${PROJECT_NAME}-swagger-local
