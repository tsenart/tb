# Token Bucket (tb) [![Build Status](https://drone.io/github.com/tsenart/tb/status.png)](https://drone.io/github.com/tsenart/tb/latest) [![GoDoc](https://godoc.org/github.com/tsenart/tb?status.png)](https://godoc.org/github.com/tsenart/tb)

This package provides a generic lock-free implementation of the "Token bucket"
algorithm where the handling of non-conformity is left to the user.
Read more about it in this [Wikipedia page](http://en.wikipedia.org/wiki/Token_bucket).
![Image](http://sardes.inrialpes.fr/~krakowia/MW-Book/Chapters/QoS/Chapters/QoS/Figs/bucket.gif)


## Install
```shell
$ go get github.com/tsenart/tb
```

## Usage
Read up the [docs](https://godoc.org/github.com/tsenart/tb) and have a look at some [examples](examples/).

## Licence
```
The MIT License (MIT)

Copyright (c) 2014 Tomás Senart

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```
