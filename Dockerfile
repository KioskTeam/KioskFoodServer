FROM kioskteam/food-server-base:1.3.1-godep
ADD . /go/src/github.com/KioskTeam/KioskFoodServer
WORKDIR /go/src/github.com/KioskTeam/KioskFoodServer
RUN godep go install
