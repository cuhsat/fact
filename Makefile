.PHONY: all clean

all: build tools

tools:
	"$(CURDIR)/scripts/eztools.sh"

build:
	"$(CURDIR)/scripts/gobuild.sh"

clean:
	"$(CURDIR)/scripts/goclean.sh"
