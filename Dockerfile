FROM public.ecr.aws/docker/library/golang:1.21 as build-image
WORKDIR /src

COPY go.mod go.sum ./
COPY . ./

RUN go mod tidy

WORKDIR /src/cmd/bonded
RUN go build -o /src/lambda-handler

FROM public.ecr.aws/lambda/provided:al2
COPY --from=build-image /src/lambda-handler .
ENTRYPOINT ./lambda-handler
