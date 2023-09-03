docker build -t strava_db ./database
docker build -t strava_api ./api
docker-compose up -d