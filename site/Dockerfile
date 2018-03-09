FROM node:alpine

EXPOSE 5018

ADD ./ /app

RUN cd /app \
  && yarn install && yarn build \
  && rm -rf node_modules/ \
  && yarn install --production \
  && yarn cache clean \
  && yarn autoclean --force

HEALTHCHECK --interval=10s --timeout=3s \
  CMD node /app/check.js || exit 1

CMD [ "node", "/app/app.js" ]
