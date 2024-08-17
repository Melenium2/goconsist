# go-consit

The not concurrency safety implementation of consistency hashing algorithm +
virtual nodes.

More info about consistency hashing you can look at 
https://liuzhenglaichn.gitbook.io/system-design/advanced/consistent-hashing.

## Get started

Lets install the library.

```sh
go get github.com/Melenium2/goconsist
```

### Initialize the ring

```go
func main() {
  // Servers that should be distributed across the ring.
  servers := []netip.AddrPort{
	  netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 10),
	  netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 20),
	  netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 30),
  }

  config := goconsist.Config{
    // SectionFactor is a range of numbers included to single ring section.
    //
    // Example:
    //  Given a ring of 3 ranges:
    //  0 - 2, 3 - 5, 6 - 0.
    //  In this case, shard factor equals to 2.
    SectionFactor: 10,
    // SectionCount is a number of ranges located in the ring.
    //
    // Example:
    //  0 - 1, 2 - 3, 4 - 5, 6 - 0.
    //  In this case ranges count equals to 4.
    SectionCount: 30,
  }

  ring := goconsist.NewRing(config, servers...)

  // ...
}
```

### Add servers

You can initialize the ring without servers and adds the servers later.
Attention!! If no servers provided to the ring, each request to acquire 
the server address will return the empty netip.AddrPort{} structure.

```go
func main() {
  // Servers that should be distributed across the ring.
  servers := []netip.AddrPort{
    netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 10),
    netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 20),
	  netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 30),
  }

  config := goconsist.Config{
    // SectionFactor is a range of numbers included to single ring section.
    //
    // Example:
    //  Given a ring of 3 ranges:
    //  0 - 2, 3 - 5, 6 - 0.
    //  In this case, shard factor equals to 2.
    SectionFactor: 10,
    // SectionCount is a number of ranges located in the ring.
    //
    // Example:
    //  0 - 1, 2 - 3, 4 - 5, 6 - 0.
    //  In this case ranges count equals to 4.
    SectionCount: 30,
  }

  ring := goconsist.NewRing(config)

  // Adds servers after initializing.
  // Each call of AddServers trigger the "distribution" of
  // the servers across the ring.
  ring.AddServers(servers...)
}
```

### Remove servers

At any time you can delete any existing server in the ring. If the server 
does not exist, the function does nothing. 

```go
func main() {
  server1 := netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 10)
  notExistedServer := netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 1, 1, 1}), 555)

  // Servers that should be distributed across the ring.
  servers := []netip.AddrPort{
    server1,
    netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 20),
    netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 30),
  }

  ring := goconsist.NewRing(goconsist.Config{}, servers)

  // Removes server1 from the ring.
  // Each call of RemoveServer trigger the "distribution" of
  // the servers across the ring.
  ring.RemoveServer(server1)

  // Do noting if the server is not present.
  ring.RemoveServer(notExistedServer)
}
```

### Acquire server 

The function calculate a hash for specified key then search for section
where this key is included and return the server from this section.
The library uses the murmur3 hash algorithm.

```go
func main() {
  // Servers that should be distributed across the ring.
  servers := []netip.AddrPort{
	  netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 10),
	  netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 20),
	  netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 30),
  }

  ring := goconsist.NewRing(goconsist.Config{}, servers)

  // Got server netip.AddrPort{} structure.
  server := ring.Acquire([]byte("any-value-you-provide"))
}
```
