help:
	@echo -e "\nYou can run all tests using:\t\t\t\tmake test\n"
	@echo -e "You can benchmark the code using:\t\t\tmake benchmark\n"
	@echo -e "You can check test coverage using:\t\t\tmake coverage\n"
	@echo -e "You can generate a test coverage report using:\t\tmake report\n"
	@echo -e "You can clean everything up using:\t\t\tmake clean\n"

test:
	@go test -v -race -timeout 30s ./...

benchmark:
	@go test -bench=. -benchmem -v ./...

coverage:
	@go test -cover ./...

report:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

clean:
	rm ./coverage.out