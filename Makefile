update:
	bazel run //:gazelle

update-deps:
	bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_library_dependencies

.PHONY: update update-deps