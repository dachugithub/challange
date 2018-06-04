# challange
to build: go build (all dependencies are vendored)
to test agains external psql db : go test -v

To test standalone (ie. CI):
docker-compose -f docker-compose-ci.yml up
NOTE:
I am using docker-machine. Please adjust docker-compose-ci.yml APP_DB_HOST: "192.168.99.100"
to your docker networking model.

to build final image:
docker build -t final-api .

NOTE:
app assumes following env vars to be set and will fail if they will be not:
POSTGRES_USER:
POSTGRES_PASSWORD: 
POSTGRES_DB: 
APP_DB_HOST: 
APP_SERVICE_PORT: 

Diagram is created with draw.io. After downloading it can be opened there :
https://www.draw.io/
