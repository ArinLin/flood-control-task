redis-up:
	docker run --name redis-vk -p 6379:6379 -d redis

redis-down:
	docker rm -f redis-vk

ratelimiter-start:
	make redis-up && go run . && make redis-down