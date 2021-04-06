# wifinet

tunneling wifi packets over the internet

## Proposed Architecture

```
             \    /                                                                802.11
              \  /                                                                 frames   
            |------|     -------\         custom       ((  )            -------      \__/
802.11 <~~> | wifi | ==> wifinet => BPF + filters <=> (      ) )  <===> wifinet <~~> |  | 
frames <~~> | card |     -------/                    ( internet )       -------      |__|
            |------|                                                  

```


## Usage

needs to run on both ends

```
go build
./wifinet -i <iface> -f <bpf filter> [other options]
```

