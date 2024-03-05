FROM node:20-alpine3.19 as frontend-development

RUN mkdir /frontend
WORKDIR /frontend

COPY ./package.json ./
COPY ./package-lock.json ./
RUN --mount=type=cache,id=npmmod,target="/root/.npm" npm install --include=dev

CMD ["npm", "run", "dev"]
