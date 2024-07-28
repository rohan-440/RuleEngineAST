build:
	go build -o RuleEngineAST

run:
	go run RuleEngineAST

clean:
	rm RuleEngineAST

test:
	go test ./...