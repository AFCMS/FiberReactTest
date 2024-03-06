# syntax=docker/dockerfile:1

FROM node:20-alpine3.19 as frontend-development

RUN mkdir /frontend
WORKDIR /frontend

CMD npm install --include=dev && npm run dev
