Shipper 
==========

Shipper is a web UI for http://docker.io

[![Build Status](https://travis-ci.org/CodeRoamer/Shipper.svg?branch=master)](https://travis-ci.org/CodeRoamer/Shipper)
[![wercker status](https://app.wercker.com/status/ef5186c0ddc28e83397186c2ac549cda/s "wercker status")](https://app.wercker.com/project/bykey/ef5186c0ddc28e83397186c2ac549cda)

### Quick Start

> more wiki on the road, project is still underdeveloped.

1. `go get github.com/coderoamer/shipper` - download shipper into your *$GOPATH*
2. `cd $GOPATH/src/github.com/coderoamer/shipper`
3. `go build` - build shipper
4. `shipper dev` - download assets & put them into the right place
5. `shipper web` - to kick off your web app

### For Developers

Run `shipper test` to run unit test.

With flag `-b true` to run benchmark test, like this:
`shipper test -b true`

### Help

type `shipper` to get help