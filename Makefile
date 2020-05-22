GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get
NAME=larder
DIR=./cmd/larder

all: build
build:
	#$(GO_GET) github.com/theycallmemac/larder
	$(GO_BUILD) -o $(NAME) -v $(DIR)
install:
	$(GO_GET) github.com/theycallmemac/larder
	$(GO_BUILD) -o /bin$(NAME) -v $(DIR)
clean:
	$(GO_CLEAN) -v ./...
	rm $(NAME)
