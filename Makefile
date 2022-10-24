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
GORUN=$(GOCMD) run -ldflags="$(ldflags)" 
GOBUILD=$(GOCMD) build -ldflags="$(ldflags)" 


############################
###        DEMO APP      ###
############################
.PHONY: demo-getAccount demo-transfer demo-createAsset
demo.%: DEMO=$*
demo.%: 
	$(GORUN) ./cmd/demo/$(DEMO)
