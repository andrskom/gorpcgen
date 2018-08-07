test-nats-gen:
	@echo "+ $@"
	@go run main.go \
	    -gen.path=testservice/nats/gen \
	    -cfg.handlers-path=testservice/handlers

test-nats-run: test-nats-gen
	@echo "+ $@"
	@go run testservice/nats/main.go

test-nats-test:
	@echo "+ $@"
	@go test -tags=teste2e -failfast ./testservice/nats/testing/...

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
