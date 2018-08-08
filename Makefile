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

test-http-gen:
	@echo "+ $@"
	@go run main.go \
	    -gen.path=testservice/http/gen \
	    -cfg.handlers-path=testservice/handlers \
	    -cfg.client-tmpl=.templates/http/client.gotmpl \
	    -cfg.server-tmpl=.templates/http/server.gotmpl

test-http-run: test-http-gen
	@echo "+ $@"
	@go run testservice/http/main.go

test-http-test:
	@echo "+ $@"
	@go test -tags=teste2e -failfast ./testservice/http/testing/...

.PHONY: all \
		test-nats-gen \
		test-nats-run \
		test-nats-test \
		nats-run \
		nats-stop \
		test-http-gen \
		test-http-run \
		test-http-test
