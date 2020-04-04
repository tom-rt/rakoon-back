PID         = /tmp/rakoon-api-gateway.pid
GO_FILES    =  $(*.go)
APP         = ./api-gateway

serve:  restart
	@fswatch -o . | xargs -n1 -I{} make restart || make kill

kill:
	@kill `cat $(PID)` || true

before:
	@echo "## RELOAD ##"

$(APP): $(GO_FILES)
	@go build $? -o $@

restart: kill before $(APP)
	@./api-gateway & echo $$! > $(PID)