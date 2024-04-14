.PHONY: lint
lint:
	black . --line-length 120
	# isort . --line-length=120 --ca --profile=black --honor-noqa --skip bazel-bin --skip bazel-out --skip bazel-chico --skip bazel-testlogs