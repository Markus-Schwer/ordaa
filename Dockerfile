FROM docker.io/node:18

WORKDIR /usr/dotinder/app

ADD package*.json .
ADD tsconfig.json .
ADD src/ ./src/

# TODO: change to npm ci for production use
RUN npm install

CMD [ "npm", "start" ]
