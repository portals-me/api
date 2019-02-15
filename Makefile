ARG=
ENV=dev

deploy:
	node env.js ${ENV} > env.json
	apex deploy --env-file env.json --env ${ENV} ${ARG}
	rm env.json

install:
	mkdir -p ./infrastructure/local/.dynamodb
	cd ./infrastructure/local/.dynamodb; \
	wget https://s3-ap-northeast-1.amazonaws.com/dynamodb-local-tokyo/dynamodb_local_latest.tar.gz; \
	tar -xf ./dynamodb_local_latest.tar.gz

test:
	$(MAKE) startTest && $(MAKE) runTest && $(MAKE) endTest || $(MAKE) endTest

endTest:
	kill `cat .dynamo.pid`
	rm .dynamo.pid
	rm shared-local-instance.db

startTest:
	{ java -Djava.library.path=./infrastructure/local/.dynamodb/DynamoDBLocal_lib -jar ./infrastructure/local/.dynamodb/DynamoDBLocal.jar -sharedDb & }; echo $$! > .dynamo.pid
	sleep 1

	cd infrastructure/local && terraform apply -auto-approve
	sleep 2

runTest:
	export EntityTable=portals-me-test-entities; \
	export SortIndex=DataTable; \
	export FeedTable=portals-me-test-feeds; \
	go test ./functions/...
