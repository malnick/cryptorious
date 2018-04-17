.PHONY: all build test clean
THIS_FILE := $(lastword $(MAKEFILE_LIST))

os?="darwin"

all: clean test build-all

build-all:
	@$(MAKE) -f $(THIS_FILE) build os=linux
	@$(MAKE) -f $(THIS_FILE) build os=darwin
	@$(MAKE) -f $(THIS_FILE) build os=freebsd
	@$(MAKE) -f $(THIS_FILE) build os=windows

build:
	$(call i,building cryptorious for $(os))
	@bash -c "scripts/build.sh $(os)"

test:
	bash -c "./scripts/test.sh libraries unit"

docker-test:
	bash -c "./scripts/docker-test.sh"

clean:
	rm -rf ./release

# Helper Functions
define i 
	@tput setaf 6 && echo "[INFO] ==> $(1)"
	@tput sgr0
endef

