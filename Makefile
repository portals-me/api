deploy:
	node env.js > env.json
	apex deploy --env-file env.json
	rm env.json
