FROM golang:1.25-alpine AS build
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/egresados_crud .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /out/egresados_crud ./egresados_crud
COPY conf ./conf
COPY swagger ./swagger
EXPOSE 8080
ENTRYPOINT ["./egresados_crud"]
