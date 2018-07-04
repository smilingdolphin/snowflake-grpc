BUILD_TAG = '1.10-alpine'
NAME = "snowflake"
PACKAGE = "github.com/smilingdolphin/snowflake-grpc"
MAIN = "$(PACKAGE)/entry"

CL_RED  = "\033[0;31m"
CL_BLUE = "\033[0;34m"
CL_GREEN = "\033[0;32m"
CL_ORANGE = "\033[0;33m"
CL_NONE = "\033[0m"

define color_out
	@echo $(1)$(2)$(CL_NONE)
endef

docker-build:
	$(call color_out,$(CL_BLUE),"Building binary in docker ...")
	@docker run --rm -v "$(PWD)":/go/src/$(PACKAGE) \
		-w /go/src/$(PACKAGE) \
		golang:$(BUILD_TAG) \
		go build -v -o $(NAME) $(MAIN)
	$(call color_out,$(CL_GREEN),"Building binary ok")

docker: docker-build
	$(call color_out,$(CL_BLUE),"Building docker image ...")
	@docker build -t $(NAME) .
	$(call color_out,$(CL_GREEN),"Building docker image ok")

build:
	@go build -v -o $(NAME) $(MAIN)

linux:
	@GOOS=linux GOARCH=amd64 go build -v -o $(NAME) $(MAIN)

.PHONY: all
all:
	build
