CACHE_CONTAINER=benthos_cache

reset:
	reset

cache:
	@-docker stop ${CACHE_CONTAINER} && docker rm ${CACHE_CONTAINER} >/dev/null
	@-docker run --rm --name ${CACHE_CONTAINER} -d redis >/dev/null

run: reset cache
	@#-[ -f result.txt ] && rm result.txt
	@#sleep 2s
	@CACHE_CONTAINER=${CACHE_CONTAINER} benthos -c ./config.yaml

run-silent:
	@$(MAKE) run 1>/dev/null
