.PHONY: full-run

full-run: test run

run:
	docker compose up --build -d

stop:
	docker compose down

test:
	go test ./...

stress-test:
	ghz --insecure --async \
      --proto ./proto/video/video.proto \
      --call video.VideoService/RedirectVideo \
      -c 10 -n 100000 --rps 10000 \
      -d '{"video":"http://s1.origin-cluster/video/123/xcg2djHckad.m3u8"}' 0.0.0.0:50051