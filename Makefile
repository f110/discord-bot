update:
	bazel run //:gazelle

update-deps:
	bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_library_dependencies

run:
	bazel run //cmd/discord-bot -- -c $(CURDIR)/config_debug.yaml

.PHONY: update update-deps run