.PHONY:
all: metering

metering: metering.go
	go build -o build/ metering.go
