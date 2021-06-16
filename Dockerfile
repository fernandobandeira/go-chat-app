FROM golang:1.16.5-alpine AS build

WORKDIR /go/src/github.com/fernandomalmeida/go-chat-app
COPY . .

RUN go build -o go-chat-app .

FROM alpine:3.12
EXPOSE 5000
WORKDIR /go-chat-app

COPY --from=build /go/src/github.com/fernandomalmeida/go-chat-app/go-chat-app /go-chat-app/go-chat-app
COPY --from=build go/src/github.com/fernandomalmeida/go-chat-app/app.prod.env /go-chat-app/app.env
COPY --from=build go/src/github.com/fernandomalmeida/go-chat-app/chat/views /go-chat-app/chat/views

CMD ["/go-chat-app/go-chat-app"]