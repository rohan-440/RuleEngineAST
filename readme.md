# Problem Statement
AST Rule Engine Service 


# Assumption & implementation
1. All the functionality are supported & tested
2. Operands are assumed are always on the LHS
3. Rules are being store in sqlite db 
4. Supported Merge Rule Strategy is "AND" & "OR"
5. Tests are added in the code. JSON file reading is not required for the tests. 
6. This read me provides with sample curls to test out the all service endpoints 
7. Server is running on port 8080 
8. Go server created using gin framework 
9. We are using sqllite disk based storage for db (Note: data is retained when app is restarted)

# TODOS
1. Further rule_id can be used in the others endpoints. It is skipped as it is out of scope for now.

# Setup, Build and un
1. Install go "brew install go"
2. Cd into the folder
3. Run below command to build the binary 
```
make build
```
4. Run below command to run the RuleEngineAST http server
```
make run
```
5. Run below command to auto test the code. 
```
make test
```

# Below are helpful curls to test the endpoints

# get all rules

```
curl --location 'http://localhost:8080/rules'
```

# create a new rule

```
curl --location 'localhost:8080/rules' \
--header 'Content-Type: application/json' \
--data '{
    "rule" : "((age > 30 AND department == '\''Marketing'\'')) AND (salary > 20000 OR experience > 5)"
}'
```


# evaluate rules

1. Rule Matching Request
```
curl --location 'localhost:8080/rules/evaluate' \
--header 'Content-Type: application/json' \
--data '{
    "rule" : "((age > 30 AND department == '\''Marketing'\'')) AND (salary > 20000 OR experience > 5)",
    "data" : {
        "age":        "31",
		"department": "Marketing",
		"salary":     "51000",
		"experience": "6"
    }
}'
```

2. Non Rule Matching Request (due to department mismatch)
```
curl --location 'localhost:8080/rules/evaluate' \
--header 'Content-Type: application/json' \
--data '{
    "rule" : "((age > 30 AND department == '\''Marketing'\'')) AND (salary > 20000 OR experience > 5)",
    "data" : {
        "age":        "31",
		"department": "Sales",
		"salary":     "51000",
		"experience": "6"
    }
}'
```

# merge rules

```
curl --location 'localhost:8080/rules/merge' \
--header 'Content-Type: application/json' \
--data '{
    "first_rule" :"age > 30",
    "second_rule" :"department == '\''Sales'\''",
    "merge_strategy" : "AND"
}'
```

# ref for lib & other helpful methods for golang

https://gorm.io/docs/update.html

https://blog.logrocket.com/rest-api-golang-gin-gorm/

https://semaphoreci.com/community/tutorials/building-go-web-applications-and-microservices-using-gin