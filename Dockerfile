FROM golang:1.3.1
ADD . /go/src/github.com/KioskTeam/KioskFoodServer
WORKDIR /go/src/github.com/KioskTeam/KioskFoodServer
RUN go install
