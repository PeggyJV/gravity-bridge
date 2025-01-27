FROM node:23.6-alpine3.21

RUN apk update
RUN apk add --no-cache git python3 make build-base

COPY package.json package.json
COPY yarn.lock yarn.lock

RUN yarn install --production=false

COPY . .

ENV ARCHIVE_NODE_URL=""
EXPOSE 8545

RUN yarn run compile

CMD yarn start
