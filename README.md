# Globally Unique ID Generator

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/jefurry/goid) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/jefurry/goid/master/LICENSE) [![Build Status](https://travis-ci.org/rs/xid.svg?branch=master)](https://travis-ci.org/jefurry/goid) [![Coverage](http://gocover.io/_badge/github.com/jefurry/goid)](http://gocover.io/github.com/jefurry/goid)

Package goid is a globally unique id generator library, ready to be used safely directly in your server code.

Xid is using 128bit and hashids algorithm to generate globally unique ids to make it shorter when transported as a string:
https://github.com/speps/go-hashids

- 8-byte value representing the milliseconds since the Unix epoch,
- 4-byte counter, starting with a random value,
- 6-bit server room id,
- 8-bit cluster id,
- 9-bit machine id,
- 3-bit work id, and
- 6-bit op id.

The string representation is using hashids to generates (28bytes).

Features:

- Size: 16 bytes (128 bits)
- Lock-free (i.e.: unlike UUIDv1 and v2)

Best used with [xlog](https://github.com/rs/xlog)'s
[RequestIDHandler](https://godoc.org/github.com/rs/xlog#RequestIDHandler).

References:

- https://github.com/rs/xid

## Install

    go get github.com/jefurry/goid

## Usage

```go
id := goid.New()

s, _ := id.Encode(0, 0, 0, 0, 0)
err := id.FromString(s)
```

Get `goid` embedded info:

```go
id.Timestamp()
id.Counter()
id.ServerRoomID()
id.ClusterID()
id.MachineID()
id.WorkID()
id.OpID()
```

## Licenses

All source code is licensed under the [BSD License](https://raw.github.com/jefurry/goid/master/LICENSE).
