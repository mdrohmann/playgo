# PlayGo

This repository includes some projects, in which I play with some of the
concepts of the Go language.

## GoTcp

GoTcp is a class client/server duo that communicates via TCP.  In order to use
the concurrency features (go-routines), the client is opening a 1000
connections in parallel.  And the server is synchronizing some state (a
counter) between all connections.

Synchronization on the server size is implemented in two different ways:

- with channel synchronization (TODO: make it an option to use buffers)
- lock-free via the `sync/atomic` package.

I tried to benchmark the `sync/atomic`, and channel synchronization
with and without buffers, but noticed that this has to be done without the TCP
communication.  If I increase the number of parallel connections too much, some
TCP packets have to be re-sent causing delays that exceed the packet
synchronization by orders of magnitude. My first impression, however, is that
there will be hardly any difference for such a simple example.

TODO: Make a useful benchmark

## CgoSpy

CgoSpy is an application that links against the `opencv` (C++) library with
`Cgo`:

It attempts to take a picture from the main available camera device on the
machine it runs on.  The resultant image matrix is wrapped into a Go structure
that implements the `image.Image` interface.  This way, I can write out /
encode the image as a PNG file with the Go standard library:  Go interfaces
rock!

In my experience `Cgo` was very straightforward, and I appreciated that `go
build` calls it automatically and correctly.  The only part that required some
research was wrapping a C array in a GoLang slice structure with:

```go
slice := (*[1 << 28]C.CInt)(unsafe.Pointer(intArray))[:sizeLen:sizeLen]
```

This casts the C array `intArray` into a giant array of size `2^28` of type
`C.CInt` and then returns a slice to that array of correct length and capacity
`sizeLen`. The argument after the second colon in the slice argument limits the
capacity, which for simple slices is set to the capacity of the underlying
array (minus the start position).  In this case, the capacity would be set to
`2^28` which does not reflect reality.

It took a while to find the information [about the trick itself](https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices), and the [capacity limitation of array slices](https://golang.org/ref/spec#Slice_expressions).

## Conclusion

All in all, I was pleased how easy it is to interface with C libraries, and to
write efficient parallel and scalable programs without having to care much
about dead-locks or synchronization.


