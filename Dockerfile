FROM mkhoeini/golang:1.3.1with_goose_godep
ADD . /go/src/github.com/KioskTeam/KioskFoodServer
WORKDIR /go/src/github.com/KioskTeam/KioskFoodServer
RUN go install
