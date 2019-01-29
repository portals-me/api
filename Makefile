ARG =
ENV = dev

deploy:
	node env.js ${ENV} > env.json
	apex deploy --env-file env.json --env ${ENV} ${ARG}
	rm env.json
