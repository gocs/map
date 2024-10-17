FROM node:18-alpine AS static

WORKDIR /usr/src/app

COPY ./static/package.json .

RUN npm install

COPY ./static .

CMD [ "npm", "run", "build" ]

FROM golang:1.23-alpine3.19

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
COPY --from=static /usr/src/app/dist dist
RUN go build -v -o /usr/local/bin/app

CMD ["app"]