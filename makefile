# Ward Jaeger, CS 403
GO_PKG:=src
EXE_NAME:=wixme.exe

all: build run

build:
	go build -o ${EXE_NAME} ${GO_PKG}/*.go

run:
	./${EXE_NAME} ${FILE}

clean:
	rm ${EXE_NAME}