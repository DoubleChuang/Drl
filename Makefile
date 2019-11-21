DEBUG_MODE=0

BIN_DIR = bin/

MAIN_FILES := $(wildcard ./cmd/*/main.go)
CMDS := $(patsubst ./cmd/%/main,%,$(basename $(MAIN_FILES)))
PKG_LIST := $(shell go list ./... )

LDFLAG=
GCFLAG=
ifeq ($(DEBUG_MODE), 0)
	LDFLAG:=-ldflags "-s -w" 
else
	GCFLAG:=-gcflags="-N -l"
endif	
all:
	make target

target:$(CMDS)

$(CMDS):
	@go build $(GCFLAG)  $(LDFLAG) -o $(BIN_DIR)$@ ./cmd/$@
	
clean:
	@rm -rf $(BIN_DIR)
.PHONY: all target clean $(CMDS)
