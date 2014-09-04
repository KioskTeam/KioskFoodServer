FROM kioskteam/food-server-base:1.3.1-godep-goose
ADD . /go/src/github.com/KioskTeam/KioskFoodServer
WORKDIR /go/src/github.com/KioskTeam/KioskFoodServer
RUN go install
