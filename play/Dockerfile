FROM node:alpine3.13

RUN npm install -g serve

COPY public /public

ENTRYPOINT [ "serve", "--single", "public", "--listen", "8080"]
