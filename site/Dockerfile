FROM node:alpine

ENV NODE_ENV=production

ADD ./ /app

RUN cd /app \
  && yarn install && yarn build \
  && rm -rf node_modules/ \
  && yarn install --production \
  && yarn cache clean

CMD [ "node", "/app/app.js" ]
