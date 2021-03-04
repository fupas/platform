.PHONY: all
all: static

.PHONY: static
static:
	cd examples/static && gcloud app deploy . --quiet --project=fupas-platform
