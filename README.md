# wifinet

tunneling wifi packets over the internet


CONTRIBUTIONS WELCOME! If you like the idea, please help out by sending Pull Requests :)

## Proposed Architecture

```
             \    /                                                                802.11
              \  /                                                                  frames   
            |------|     -------          custom        ((  )          -------       \__/
802.11 <~~> | wifi |<==> wifinet => BPF + filters <=> (       )  <===> wifinet <~~>  |  | 
frames <~~> | card |     -------                    ( internet )       -------       |__|
            |------|                                                  

```


## Usage

needs to run on both ends

```
go build
./wifinet -i <iface> -f <bpf filter> [other options]
```

## Progress

- [x] sniff 802.11 packets
- [ ] channel selection / hopping
- [ ] client/server UDP communication
- [ ] tunnel packet bytes
- [ ] reinject frame
- [ ] duplex communications with multiple threads (go routines?)
- [ ] \[optional\] symmetric key encryption?
