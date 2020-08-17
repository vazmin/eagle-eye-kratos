FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
#RUN go build -o main ./service/organization/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main ./service/licensing/cmd
FROM scratch
COPY --from=builder /build/service/licensing/configs /app/configs
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main", "-conf", "./configs"]