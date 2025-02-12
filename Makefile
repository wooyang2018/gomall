.PHONY: all
all: help

#当执行 make 时，如果没有指定目标，它会调用 help 目标。
default: help

.PHONY: help
#这段代码的作用是生成一个帮助文档，列出 Makefile 中的所有目标及其注释，并按照章节进行分组。
#$(MAKEFILE_LIST) 是一个 Makefile 变量，包含当前 Makefile 的文件名
help: ## display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Initialize Project

.PHONY: init
init: ## Just copy `.env.example` to `.env` with one click, executed once.
	@scripts/copy_env.sh

.PHONY: gen-client
gen-client: ## gen client code of {svc}. example: make gen-client svc=product
	@cd rpc_gen && cwgo client --type RPC --service ${svc} --module github.com/cloudwego/biz-demo/gomall/rpc_gen  -I ../idl  --idl ../idl/${svc}.proto

.PHONY: gen-server
gen-server: ## gen service code of {svc}. example: make gen-server svc=product
	@cd app/${svc} && cwgo server --type RPC --service ${svc} --module github.com/cloudwego/biz-demo/gomall/app/${svc} -I ../../idl  --idl ../../idl/${svc}.proto

.PHONY: gen-frontend
gen-frontend: ## gen frontend code
	@cd app/frontend && cwgo server -I ../../idl --type HTTP --service frontend --module github.com/cloudwego/biz-demo/gomall/app/frontend --idl ../../idl/frontend/checkout_page.proto

##@ Development Env

.PHONY: watch-frontend
watch-frontend: ## run frontend with air
	@cd app/frontend && air

.PHONY: watch-svc
watch-svc: ## run {svc} server with air. example: make watch-svc svc=product
	@cd app/${svc} && air

.PHONY: run
run: ## run {svc} server. example: make run svc=product
	@scripts/run.sh ${svc}

#gofumpt：gofmt 的扩展工具，提供了更严格的代码格式化规则。
#-l：列出需要格式化的文件，但不实际修改文件内容。
#-w：将格式化后的内容写回源文件（即直接修改文件）。
#app：指定要格式化的目标目录或文件（这里是 app 目录）。
.PHONY: lint
lint: ## run `gofmt` for all go module
	@gofmt -l -w app
	@gofumpt -l -w  app

.PHONY: vet
vet: ## run `go vet` for all go module
	@scripts/vet.sh

.PHONY: lint-fix
lint-fix: ## run `golangci-lint` for all go module
	@scripts/fix.sh

##@ Build

.PHONY: build-frontend
build-frontend: ## build frontend image
	docker build -f ./deploy/Dockerfile.frontend -t frontend:${v} .

#--build-arg 允许你在构建镜像时传递一些变量值到 Dockerfile 中。
#这些变量可以在 Dockerfile 中使用 ARG 指令定义，并在构建过程中使用。
.PHONY: build-svc
build-svc: ## build {svc} image. example: make build-svc svc=product v=v0.1
	docker build -f ./deploy/Dockerfile.svc -t ${svc}:${v} --build-arg SVC=${svc} .

.PHONY: build-all
build-all: ## build all images
	docker build -f ./deploy/Dockerfile.frontend -t frontend:${v} .
	docker build -f ./deploy/Dockerfile.svc -t cart:${v} --build-arg SVC=cart .
	docker build -f ./deploy/Dockerfile.svc -t checkout:${v} --build-arg SVC=checkout .
	docker build -f ./deploy/Dockerfile.svc -t email:${v} --build-arg SVC=email .
	docker build -f ./deploy/Dockerfile.svc -t order:${v} --build-arg SVC=order .
	docker build -f ./deploy/Dockerfile.svc -t payment:${v} --build-arg SVC=payment .
	docker build -f ./deploy/Dockerfile.svc -t product:${v} --build-arg SVC=product .
	docker build -f ./deploy/Dockerfile.svc -t user:${v} --build-arg SVC=user .

.PHONY: run-frontend
run-frontend: ## run frontend with docker
	docker run -p 8080:8080 --network gomall_default frontend:${v}

.PHONY: run-svc
run-svc: ## run {svc} with docker. example: make run-svc svc=product
	docker run --network gomall_default ${svc}:${v}

.PHONY: env-start
env-start:  ## launch all middleware software as the docker
	@docker-compose up -d

.PHONY: env-stop
env-stop: ## stop all docker
	@docker-compose down

.PHONY: clean
clean: ## clern up all the tmp files
	@rm -r app/**/log/ app/**/tmp/

##@ Open Browser

.PHONY: open.gomall
open-gomall: ## open `gomall` website in the default browser
	@open "http://localhost:8080/"

.PHONY: open.consul
open-consul: ## open `consul ui` in the default browser
	@open "http://localhost:8500/ui/"

.PHONY: open.jaeger
open-jaeger: ## open `jaeger ui` in the default browser
	@open "http://localhost:16686/search"

.PHONY: open.prometheus
open-prometheus: ## open `prometheus ui` in the default browser
	@open "http://localhost:9090"