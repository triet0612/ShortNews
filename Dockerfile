FROM golang:1.22.0 as api_stage
WORKDIR /app
COPY . .
RUN go build ./cmd/main.go

FROM node:lts as frontend_stage
WORKDIR /app
COPY ./frontend .
RUN npm run build

FROM debian:stable
RUN apt-get update -y
RUN apt-get install -y ca-certificates
WORKDIR /app
COPY --from=api_stage /app/main /app
COPY --from=frontend_stage /app/build /app/build

ENTRYPOINT [ "/app/main" ]
