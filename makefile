all:
	wails build -platform darwin/arm64,darwin/amd64,windows/amd64,windows/arm64,linux/amd64,linux/arm64

clean:
	rm -rf build/bin
