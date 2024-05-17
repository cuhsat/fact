.PHONY: all clean

all: build

build:
	"$(CURDIR)/scripts/gobuild.sh"

clean:
	"$(CURDIR)/scripts/goclean.sh"
