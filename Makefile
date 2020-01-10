update:
	bazel run //:gazelle

update-deps:
	bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_library_dependencies

push: push-bot push-object-proxy \

push-bot:
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //:image.tar
	docker load -i bazel-bin/image.tar
	docker tag bazel:image quay.io/f110/discord-bot:latest
	docker push quay.io/f110/discord-bot:latest
	docker rmi bazel:image

push-object-proxy:
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //:proxy_image.tar
	docker load -i bazel-bin/proxy_image.tar
	docker tag bazel:proxy_image quay.io/f110/object-proxy:latest
	docker push quay.io/f110/object-proxy:latest
	docker rmi bazel:proxy_image

.PHONY: push-bot push-object-proxy update update-deps