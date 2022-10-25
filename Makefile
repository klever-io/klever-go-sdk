# KKKKKKKKK    KKKKKKKLLLLLLLLLLL             EEEEEEEEEEEEEEEEEEEEEEVVVVVVVV           VVVVVVVVEEEEEEEEEEEEEEEEEEEEEERRRRRRRRRRRRRRRRR   
# K:::::::K    K:::::KL:::::::::L             E::::::::::::::::::::EV::::::V           V::::::VE::::::::::::::::::::ER::::::::::::::::R  
# K:::::::K    K:::::KL:::::::::L             E::::::::::::::::::::EV::::::V           V::::::VE::::::::::::::::::::ER::::::RRRRRR:::::R 
# K:::::::K   K::::::KLL:::::::LL             EE::::::EEEEEEEEE::::EV::::::V           V::::::VEE::::::EEEEEEEEE::::ERR:::::R     R:::::R
# KK::::::K  K:::::KKK  L:::::L                 E:::::E       EEEEEE V:::::V           V:::::V   E:::::E       EEEEEE  R::::R     R:::::R
#   K:::::K K:::::K     L:::::L                 E:::::E               V:::::V         V:::::V    E:::::E               R::::R     R:::::R
#   K::::::K:::::K      L:::::L                 E::::::EEEEEEEEEE      V:::::V       V:::::V     E::::::EEEEEEEEEE     R::::RRRRRR:::::R 
#   K:::::::::::K       L:::::L                 E:::::::::::::::E       V:::::V     V:::::V      E:::::::::::::::E     R:::::::::::::RR  
#   K:::::::::::K       L:::::L                 E:::::::::::::::E        V:::::V   V:::::V       E:::::::::::::::E     R::::RRRRRR:::::R 
#   K::::::K:::::K      L:::::L                 E::::::EEEEEEEEEE         V:::::V V:::::V        E::::::EEEEEEEEEE     R::::R     R:::::R
#   K:::::K K:::::K     L:::::L                 E:::::E                    V:::::V:::::V         E:::::E               R::::R     R:::::R
# KK::::::K  K:::::KKK  L:::::L         LLLLLL  E:::::E       EEEEEE        V:::::::::V          E:::::E       EEEEEE  R::::R     R:::::R
# K:::::::K   K::::::KLL:::::::LLLLLLLLL:::::LEE::::::EEEEEEEE:::::E         V:::::::V         EE::::::EEEEEEEE:::::ERR:::::R     R:::::R
# K:::::::K    K:::::KL::::::::::::::::::::::LE::::::::::::::::::::E          V:::::V          E::::::::::::::::::::ER::::::R     R:::::R
# K:::::::K    K:::::KL::::::::::::::::::::::LE::::::::::::::::::::E           V:::V           E::::::::::::::::::::ER::::::R     R:::::R
# KKKKKKKKK    KKKKKKKLLLLLLLLLLLLLLLLLLLLLLLLEEEEEEEEEEEEEEEEEEEEEE            VVV            EEEEEEEEEEEEEEEEEEEEEERRRRRRRR     RRRRRRR


# LOAD builder info
ifndef VERSION
SHELL := /bin/bash
VERSION := $(shell git describe --always --long --dirty --tag)
endif
ldflags := -X 'main.appVersion=${VERSION}'

GOCMD=go
GOIMPORT=goimports
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GORUN=$(GOCMD) run -ldflags="$(ldflags)" 
GOBUILD=$(GOCMD) build -ldflags="$(ldflags)" 


############################
###        DEMO APP      ###
############################
demo.%: DEMO=$*
demo.%: 
	$(GORUN) ./cmd/demo/$(DEMO)

############################
###    TESTS/COVERAGE    ###
############################
.PHONY: dependencies vet test test_coverage

dependencies:
	$(GOINSTALL) github.com/nikolaydubina/go-cover-treemap@latest
	$(GOINSTALL) golang.org/x/tools/cmd/goimports@latest

vet:
	$(GOCMD) vet ./...

imports:
	$(eval CHECKFILES = $(shell find . -type f -name '*.go' -not -name '*.pb.go' -not -name '*_setter.go' -not -path './vendor/*'))
	$(GOCMD) mod tidy
	$(GOIMPORT) -d -w $(CHECKFILES)

test:
	$(GOCMD) clean -testcache
	$(GOTEST) ./...

test_coverage:
	$(eval PACKAGES = $(shell go list ./... | grep -v '/proto/'))
	$(GOTEST) -coverprofile cover.out $(PACKAGES)
	go-cover-treemap -coverprofile cover.out > out.svg

test_coveragehtml:
	$(eval PACKAGES = $(shell go list ./... | grep -v '/proto/'))
	$(GOTEST) -coverprofile cover.out $(PACKAGES)
	$(GOCMD) tool cover -html=cover.out
