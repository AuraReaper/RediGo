run:
	docker run --rm -it \
		-v $(PWD):/app \
		-w /app \
		-p 7379:7379 \
		golang:1.25 \
		go run .

bash:
	docker run --rm -it \
		-v $(PWD):/app \
		-w /app \
		golang:1.25 \
		bash
