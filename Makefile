#Fixes a bug in OSX Make with exporting PATH environment variables
#See: http://stackoverflow.com/questions/11745634/setting-path-variable-from-inside-makefile-does-not-work-for-make-3-81
export SHELL := $(shell echo $$SHELL)

#get the path to this Makefile, its the last in this list
#MAKEFILE_LIST is the list of Makefiles that are executed.
TOP := $(dir $(lastword $(MAKEFILE_LIST)))
ROOT = $(realpath $(TOP))

#set the bin directory so that it's in our path for convenience
PATH := $(ROOT)/bin:$(PATH)
DEFAULT_APP := myapp
export PATH

export GO15VENDOREXPERIMENT := 1
export GOPATH := $(ROOT)/vendor:$(ROOT)
.DEFAULT_GOAL := build

#Update all imports, and remove any that aren't necessary, for all go files we can find.
imports:
	 find $(ROOT) -name '*.go' -exec goimports -w {} \;

build:
	goapp build $(DEFAULT_APP)

serve:
	foreman start

dev_server:
	goapp serve $(DEFAULT_APP)
	#foreman start

deploy:
	goapp deploy $(DEFAULT_APP)

info:
	echo $(GOPATH)

deps:
	goapp get "gopkg.in/dommmel/go-shopify.v2"
	goapp get "google.golang.org/appengine"
	goapp get "github.com/julienschmidt/httprouter"