ARG = ''

deploy:
	node env.js > env.json
	apex deploy --env-file env.json ${ARG}
	rm env.json
