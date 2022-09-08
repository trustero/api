SUB_DIRS := go docs
MODULE := api
include ./make/globals.mk

fmt:
	$(call outmsg, $(MODULE))
	@find . -name *.go -exec ./make/license_header.sh {} \;
	@find . -name *.proto -exec ./make/license_header.sh {} \;
