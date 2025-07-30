FROM golang:1.22 as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o feedback ./services/feedback

FROM gcr.io/distroless/base-debian12
COPY --from=build /app/feedback /feedback
CMD ["/feedback"]
