FROM golang:1.19-bullseye

ENV GOPATH=/

RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh  \
    && go install github.com/swaggo/swag/cmd/swag@latest \
    && sh install.sh  \
    && cp ./bin/air /bin/air
COPY go.mod go.sum /app/
WORKDIR /app

RUN go mod download
COPY . .
CMD air