FROM public.ecr.aws/docker/library/golang:1.19 as build-image
WORKDIR /src

COPY . ./

RUN go mod download

WORKDIR /src/main
RUN go build -o /src/lambda-handler

FROM public.ecr.aws/lambda/provided:al2
COPY --from=build-image /src/lambda-handler .
ENTRYPOINT ./lambda-handler