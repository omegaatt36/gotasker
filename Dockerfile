FROM golang:1.22-alpine as build

WORKDIR /go/src/app

COPY ["go.mod", "go.sum", "./"]

RUN ["go", "mod", "download"]

COPY . .

ENV APP_NAME=gotasker

RUN ["go", "build", "-o", "build/${APP_NAME}"]

FROM gcr.io/distroless/static-debian12 as prod

WORKDIR /home/app/

COPY --from=build /go/src/app/build/${APP_NAME} ./

CMD ["./${APP_NAME}"]
