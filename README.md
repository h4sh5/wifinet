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

needs to run on both ends.  sudo / root perms required because of packet sniffing things (unless you do setcap on the binary correctly)

```
go build
# iface wlan0, remote host 1.2.3.4:1234
sudo ./wifinet -i wlan0 -r 1.2.3.4:1234
```

## Progress

- [x] sniff 802.11 packets
- [x] channel selection (lock and filter on only one channel)
- [x] client/server UDP communication
- [x] send/recv frames
- [x] reinject frames
- [x] save to pcap
- [ ] MAC address filtering
- [ ] channel hopping
- [ ] tested duplex communication
- [ ] \[optional\] symmetric key encryption?
