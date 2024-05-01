FROM node:lts as frontend_stage
WORKDIR /app
COPY ./frontend .
RUN npm install
RUN npm run build

FROM golang:1.22.0 as api_stage
WORKDIR /app
COPY . .
COPY --from=frontend_stage /app/build /app/frontend/build
RUN go build ./cmd/main.go

FROM debian:stable
RUN apt-get update -y
RUN apt-get install -y ca-certificates
WORKDIR /app
COPY --from=api_stage /app/main /app

ENTRYPOINT [ "/app/main" ]
