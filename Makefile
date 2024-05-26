.PHONY: bootstrap diff deploy destroy gen-file build-migration-container run-migration-container docker-network openapi gen-models

ENV=dev
PROJECT=boilerplate-dev
BootstrapBucketName=weigkbweaksnfwelansfwekbak

diff:
	cd cdk && \
	cdk diff \
	-c ABN=$(BootstrapBucketName) -c ENV=$(ENV) -c P=$(PROJECT) \
	--profile sasaki

deploy:
	cd cdk && \
	cdk deploy \
	-c ABN=$(BootstrapBucketName) -c ENV=$(ENV) -c P=$(PROJECT) \
	--profile sasaki

destroy: 
	cd cdk && \
	cdk destroy \
	-c ABN=$(BootstrapBucketName) -c ENV=$(ENV) -c P=$(PROJECT) \
	--profile sasaki

docker-network:
	docker network create -d bridge boilerplate

migration-file:
	migrate create -ext sql -dir migration/migrations $(F)

# build-migration-container:
# 	docker build -f docker/migration/Dockerfile.dev -t migration .

build-migration-container: # apple silicon
	docker build -f docker/migration/Dockerfile.dev -t migration --build-arg ImagePrefix=arm64v8/ .

run-migration-container: build-migration-container
	docker run --net container:boilerplate_db migration:latest

openapi:
	oapi-codegen -old-config-style -generate "spec" -package restapi ./server/openapi.yml > ./server/restapi/spec.gen.go
	oapi-codegen -old-config-style -generate "types" -package restapi ./server/openapi.yml > ./server/restapi/types.gen.go
	oapi-codegen -old-config-style -generate "chi-server" -package restapi ./server/openapi.yml > ./server/restapi/server.gen.go

models:
	sqlboiler psql