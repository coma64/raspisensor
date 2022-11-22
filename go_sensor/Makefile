arm_name = raspisensor-arm
remote_bin_path = /home/pi/.local/bin

all: build-arm execute-on-raspi

make-build-folder:
	mkdir -p build

build-amd64: make-build-folder
	go build -o build/raspisensor ./main.go

build-arm: make-build-folder
	GOOS=linux GOARCH=arm GOARM=5 go build -o build/$(arm_name) ./main.go

execute-on-raspi: copy-to-raspi
	ssh raspi $(remote_bin_path)/$(arm_name)

copy-to-raspi:
	scp build/$(arm_name) raspi:$(remote_bin_path)
