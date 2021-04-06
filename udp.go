package main

import (
    "encoding/hex"
    "net"
    "log"
)



func udpListen(server string) {
    /* Lets prepare a address at any address at port 10001*/   
    ServerAddr,err := net.ResolveUDPAddr("udp", server)
    CheckError(err)
 
    /* Now listen at selected port */
    log.Printf("Listening on %s \n", server)
    ServerConn, err := net.ListenUDP("udp", ServerAddr)
    CheckError(err)
    defer ServerConn.Close()
 
    buf := make([]byte, 65536)
 
    for {
        n,addr,err := ServerConn.ReadFromUDP(buf)
        log.Printf("Received len %v from addr %v\n:", n, addr)
        log.Printf("%s", hex.Dump(buf[0:n]))

        // inject it back into wireless interface
 
        if err != nil {
            log.Println("Error: ",err)
        } 
    }
}