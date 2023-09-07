FROM golang:latest

RUN sudo hwclock --hctosys; apt update
RUN apt install -y vim
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install golang.org/x/tools/gopls@latest
