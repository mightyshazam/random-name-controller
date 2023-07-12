docker_build(
    'controller',
    '.',
    dockerfile_contents="""
FROM golang:alpine as build
WORKDIR /go/random_string_controller
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o controller pkg/cmd/controller/main.go

FROM alpine
WORKDIR /app
COPY --from=build /go/random_string_controller/controller .
ENTRYPOINT ["./controller"]
    """,
    only=['pkg/', 'go.mod', 'go.sum'],
    extra_tag=['development']
)

k8s_yaml(kustomize('manifests/local'))