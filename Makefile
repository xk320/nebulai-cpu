BINARY_NAME = nebulai-cpu
PLATFORMS = windows/amd64 windows/386 darwin/amd64 darwin/arm64 linux/amd64 linux/arm64

build: ## 构建当前系统版本
	go build -o $(BINARY_NAME) main.go

release: clean ## 构建所有平台版本
	$(foreach platform, $(PLATFORMS), \
		$(eval OS = $(word 1,$(subst /, ,$(platform)))) \
		$(eval ARCH = $(word 2,$(subst /, ,$(platform)))) \
		GOOS=$(OS) GOARCH=$(ARCH) go build \
		-o bin/$(BINARY_NAME)-$(OS)-$(ARCH)$(if $(findstring windows,$(OS)),.exe,) main.go; \
	)

clean: ## 清理构建文件
	rm -f $(BINARY_NAME)
	rm -rf bin/*

.PHONY: build release clean