ARG=
ENV=dev

deploy:
	node env.js ${ENV} > env.json
	apex deploy --env-file env.json --env ${ENV} ${ARG}
	rm env.json

test:
	{ java -Djava.library.path=./infrastructure/local/.dynamodb/DynamoDBLocal_lib -jar ./infrastructure/local/.dynamodb/DynamoDBLocal.jar -sharedDb & }; echo $$! > .dynamo.pid
	sleep 1

	cd infrastructure/local && terraform apply -auto-approve

	go test ./functions/...

	kill `cat .dynamo.pid`
	rm .dynamo.pid
	rm shared-local-instance.db
