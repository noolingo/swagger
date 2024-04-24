GO_SWAGGER_VER := 0.25.0

# Never edit indents
PROTO				:= $(shell go get github.com/noolingo/proto)
PROTO_VER           := $(shell go list -m -f '{{ .Version }}' github.com/noolingo/proto;)

GOPATH ?= $(shell go env GOPATH)

SWAGGER_PATH        := $(GOPATH)/pkg/mod/github.com/noolingo/proto@$(PROTO_VER)/codegen/swagger
INFO_PATH           := ${SWAGGER_PATH}/common/info.swagger.json
COMMON_PATH         := $(SWAGGER_PATH)/common/common.swagger.json


all:
	go mod download; \
	wget -qO swagger-bin https://github.com/go-swagger/go-swagger/releases/download/v${GO_SWAGGER_VER}/swagger_linux_amd64; \
	chmod +x swagger-bin; \
	./swagger-bin mixin \
	    $(INFO_PATH) \
        $(COMMON_PATH) \
		$(SWAGGER_PATH)/noolingo/*.json \
		-o ui/swagger.json; \
		./swagger-bin validate ui/swagger.json;

#go run ./cmd/enums/... --path=./ui/swagger.json; \

# refact:
# 	go mod edit -module github.com/noolingo/swagger
# 	-- rename all imported module
# 	find . -type f -name '*.go' \
#   	-exec sed -i -e 's,github.com/MelnikovNA/noolingoswagger,github.com/noolingo/swagger,g' {} \;