all: install
install: go.sum
		GO111MODULE=on go install -tags "$(build_tags)" ./cmd/pdu
go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify
