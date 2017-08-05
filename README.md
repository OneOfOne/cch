# cch (closable channels) [![GoDoc](https://godoc.org/github.com/OneOfOne/cch?status.svg)](https://godoc.org/github.com/OneOfOne/cch) [![Build Status](https://travis-ci.org/OneOfOne/cch.svg?branch=master)](https://travis-ci.org/OneOfOne/cch)

cch is a simple wrapper over `chan interface{}` to allow multiple writers/readers with safe closing without panicing.

## Install

	go get github.com/OneOfOne/cch


## FAQ

### Why?
* Because quiet often I need to have multiple writers and readers and issue [https://github.com/golang/go/issues/15411](#15411) got rejected.

### How's the performance?
* The only overhead is a mutex.RLock on Send/Recv.

## Usage

```go
import (
	"github.com/OneOfOne/cch"
)

const NumWorkers = 100

func main() {
	ch := cch.NewBuffered(NumWorkers)
	for i := 0; i < NumWorkers; i++ {
		go worker(ch)
	}

	for i := 0; i < NumManagers; i++ {
		go manager(ch)
	}


	// wait for some condition
	ch.Close()
}
```

## TODO:

* Better documentation.
* More tests.
* Better select/multi-impl.

## License

Apache v2.0 (see [LICENSE](https://github.com/OneOfOne/cch/blob/master/LICENSE) file).

Copyright 2016-2017 Ahmed <[OneOfOne](https://github.com/OneOfOne/)> W.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
