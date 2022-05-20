# Build the executable.
FROM golang:1.17.7-bullseye AS build
ARG GITHUB_USERNAME
ARG GITHUB_TOKEN

# Github login for private repositories.
RUN echo "machine github.com login ${GITHUB_USERNAME} password ${GITHUB_TOKEN}" > ~/.netrc
RUN echo "machine api.github.com login ${GITHUB_USERNAME} password ${GITHUB_TOKEN}" >> ~/.netrc

WORKDIR /source
COPY . .
RUN go env -w GOOS="linux"
RUN go env -w GOARCH="amd64"
RUN go env -w GOPRIVATE="github.com/deliveryhero"
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.1
RUN swag init
RUN go build

# Run the executable in a different runtime environment.
FROM golang:1.17.7-bullseye as runtime
WORKDIR /app
COPY --from=build /source/go-clean-architecture /app
COPY --from=build /source/.env* /app

# Setting local timezone.
ENV TZ=Europe/Istanbul

EXPOSE 80
CMD ["./go-clean-architecture"]