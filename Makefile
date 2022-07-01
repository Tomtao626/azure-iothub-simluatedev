SRC = $(wildcard *.go)
GITVISION = `git describe`"-"`git rev-parse --short HEAD`
all:simluatedev
simluatedev:$(SRC)
	go build -o azuretask -x -ldflags  " -w -s -X main.BuildVersion=$(GITVISION)" $^
arm:

	GOARM=7 GOARCH=arm64 GOOS=linux go build -x -ldflags "-w -s -X main.BuildVersion=$(GITVISION)" -o simluatedev $^


clean:
	rm -rvf simluatedev


