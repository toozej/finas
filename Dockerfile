# setup project and deps
FROM golang:1.20-bullseye AS init

WORKDIR /go/finas/

COPY go.mod* go.sum* ./
RUN go mod download

COPY . ./

FROM init as vet
RUN go vet ./...

# run tests
FROM init as test
RUN go test -coverprofile c.out -v ./...

# build binary
FROM init as build
ARG LDFLAGS

RUN CGO_ENABLED=0 go build -ldflags="${LDFLAGS}" ./cmd/finas/

# runtime image
FROM scratch
# Copy our static executable.
COPY --from=build /go/finas/finas /go/bin/finas
# Run the binary.
ENTRYPOINT ["/go/bin/finas"]
