KioskFoodServer
===============

To back your food enjoyment!


Dependencies
------------

* [Docker](https://www.docker.com/)
* [Fig](http://www.fig.sh/)

It also requires that a container with the name `images` to exist.
And also `images` container should have a `/images` volume to serve images from
there.(e.g. look here: https://github.com/KioskTeam/LobbyImages )


How to run
----------

```
$ docker run --name images kioskteam/lobbyimages:1.0.0 true
$ git clone https://github.com/KioskTeam/KioskFoodServer.git
$ cd KioskFoodServer
$ fig up
```

Also note that DB migrations live in a separate docker image. So, if you want
to execute migrations, after you ensured that postgres service started
successfully from previous command, you can do it like this, for example:

```
$ docker run -t -i --rm --link postgres_1:postgres_1 kioskteam/foodmigrations
```
