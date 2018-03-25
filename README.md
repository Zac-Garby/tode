# tode

_NOTE: Since tode isn't finished yet, some of this README is hypothetical as of yet._

tode stands for "The Online Database of Equations". It's essentially a website where you can browse a list of equations, as well as adding your own to the database. It will be hosted on [tode.uk](http://tode.uk), but for now you'll have to host it yourself if you want to have a go with it.

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

Unless you encountered any errors, you should now have two servers running: a Redis server and a tode server. Open a browser and go to [localhost:7000](http://localhost:7000) to use the web frontend, or access the HTTP API by sending HTTP requests via curl:

```
curl http://localhost:7000/api/all/users | python -m json.tool
```

### Using the scripts

I've written some useful bash scripts for starting the servers and using the API. Instead of using up two terminal windows and manually starting the redis server and the tode server, you can alternatively execute `start`:

```
$ ./start
listening on http://localhost:7000
```

Then, in another terminal window, you can use `run` to test the API and automatically pretty-print it using Python's json.tool module.

```
$ ./run user/id/0
{
	"name": ...
	...
	...
}
```
