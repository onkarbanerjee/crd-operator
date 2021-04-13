FROM alpine:3.11.3 as builder


RUN apk add --no-cache git make musl-dev go=1.13.13-r0

COPY ./ /

WORKDIR /

RUN ls -al && pwd && go mod download

RUN go build -o operator && pwd && ls -al

################

FROM alpine:3.11.3

COPY --from=builder /operator /operator

ENTRYPOINT [ "./operator" ]