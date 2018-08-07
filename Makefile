test-nats-gen:
	@echo "+ $@"
	@go run main.go \
	    -gen.path=testservice/service/gen \
	    -cfg.handlers-path=testservice/service/handlers

test-nats-test: test-nats-gen
	@echo "+ $@"
	@go test -failfast testservice/service/nats/testing/nats_test.go

nats-run:
	@echo "+ $@"
	@docker run --rm -d --name nats -p 4222:4222 -p 6222:6222 -p 8222:8222 nats

nats-stop:
	@echo "+ $@"
	@docker stop nats

.PHONY: all \
		nats-run \
		nats-stop \
		test-nats-test
