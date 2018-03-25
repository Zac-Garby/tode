# tode

_NOTE: Since tode isn't finished yet, some of this README is hypothetical as of yet._

tode stands for "The Online Database of Equations". It's essentially a website where you can browse a list of equations, as well as adding your own to the database. It will be hosted on [tode.uk](tode.uk), but for now you'll have to host it yourself if you want to have a go with it.

## Getting started

I assume you have Go and Redis installed.

```
go get -u github.com/gorilla/mux
go get -u github.com/go-redis/redis
git clone github.com/Zac-Garby/tode
cd tode
redis-server
```

Then, in another terminal (still in the `tode` directory):

```
go install
tode
```

Unless you encountered any errors, you should now have two servers running: a Redis server and a tode server. Open a browser and go to [localhost:7000](localhost:7000) to use the web frontend, or access the HTTP API by sending HTTP requests via curl:

```
curl http://localhost:7000/api/all/users | python -m json.tool
```
