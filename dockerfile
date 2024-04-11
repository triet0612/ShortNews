FROM golang:1.22.0 as api_stage
WORKDIR /app
COPY . .
RUN go build ./cmd/main.go

FROM node:lts as frontend_stage
WORKDIR /app
COPY ./frontend .
RUN npm run build

FROM ubuntu:noble

WORKDIR /app

RUN apt-get update -y
RUN apt-get install -y ca-certificates

RUN apt-get install -y libespeak-ng1
RUN apt install -y python3-minimal
RUN apt install -y python3-pip

RUN python3 -m pip install --break-system-packages -U mycroft-mimic3-tts[all]

RUN mimic3-download en_US/cmu-arctic_low &&\
    mimic3-download vi_VN/vais1000_low

RUN apt install -y curl && \
    apt install -y pciutils

RUN curl -fsSL https://ollama.com/install.sh | sh

COPY ./bin/llm ./llm

COPY --from=api_stage /app/main /app
COPY --from=frontend_stage /app/build /app/build

ENTRYPOINT [ "/app/main" ]

EXPOSE 8000
