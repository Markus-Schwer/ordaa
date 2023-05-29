# .inder

This project was created to automate our orders from out preferred indian restaurant.

## Configuration
For configuration create a `.env` file and place this lines in there. If you run the app in docker, you can configure these values easily via docker environment variables.
``` bash
USER_ID=@your_username:matrix.org
ACCESS_TOKEN=your_access_token
BASE_URL=https://matrix.org
ROOM_ID=!your_room_id:matrix.org
```

## Commands

``` bash
# build and run the application
npm start

# check code format
npm run format:check

# format code
npm run format:write

# build docker image
docker compose build

# run the app in docker
docker compose up -d
```
