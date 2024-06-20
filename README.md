# dumbbin

a [pastebin](https://pastebin.com)-like tool for sharing texts, written in [Golang](https://go.dev), using [chi](https://github.com/go-chi/chi/).

### installation

in order to run this program, you need [Go](https://go.dev/dl/) and [Git](https://git-scm.com/) on your system.

clone this repository using
```
git clone https://github.com/keratomalacian/dumbbin
```

and then run
```
go mod tidy
go build -o dumbbin cmd/main.go
```
to build the project.

### running
after building the project, there will be a new executable called __dumbbin__. run the executable using


```
./dumbbin
```

and now dumbbin should be running at __localhost:3232__. check it out yourself!

## usage

you can send __POST__ requests to the root of the page __(/)__ to create bins. the body of the request should contain the content of the bin.

for example, to create a new bin using cURL you can run

```
curl -H 'Content-Type: text/plain' \
     -d 'Content of the bin'  \
     -X POST \
     http://localhost:3232
```

in the command line. then the server will return a link for accessing the newly created bin.

## command line arguments

dumbbin can be modified using command line arguments. the available arguments are:

- --address: sets the address in which the server will run. eg: __localhost:8080__.
- --binpath: sets the directory where bins are stored. eg: __/home/user/bins__
- --logger: enables or disables the request logger.
- --requestsize: sets the maximum number of megabytes that will be read from the request body. eg: __10__.

- --ratelimit: enables or disables ratelimiting.
- --requestlimit: sets the limit of requests that can be done in the __timelimit__ time window.
- --timelimit: sets the time window for ratelimiting. eg: __1h30m__.

## license

created by __keratomalacian__. licensed under the __GNU General Public License V3.0__. see the __LICENSE__ file for more details.