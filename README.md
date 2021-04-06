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
# iface wlan0, remote host 1.2.3.4:1234
./wifinet -i wlan0 -r 1.2.3.4:1234
```

## Progress

- [x] sniff 802.11 packets
- [ ] channel selection (lock and filter on only one channel)
- [x] client/server UDP communication
- [x] send/recv frames
- [ ] reinject frames
- [ ] channel hopping
- [ ] tested duplex communication
- [ ] \[optional\] symmetric key encryption?
