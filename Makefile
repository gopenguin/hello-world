.PHONY: all
all: hello-world timer http

.PHONY: FORCE
FORCE:

hello-world: FORCE
	go build -o hello-world ./cmd/hello-world/

timer: FORCE
	go build -o timer ./cmd/timer/

http: FORCE
	go build -o http ./cmd/http/

.PHONY: clean
clean:
	rm -f hello-world timer http