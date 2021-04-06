package main

import (
    // "encoding/hex"
    "net"
    "log"
    "github.com/google/gopacket"
    "github.com/google/gopacket/pcap"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcapgo"
    "os"
    "time"
)

var (
    pktCount int
)

func getFramesFromServer(server string, activePcapHandle *pcap.Handle, dontInject bool) {

    ServerAddr,err := net.ResolveUDPAddr("udp", server)
    CheckError(err)
 
    /* Now listen at selected port */
    log.Printf("Listening on %s \n", server)
    ServerConn, err := net.ListenUDP("udp", ServerAddr)
    CheckError(err)
    defer ServerConn.Close()
 
    buf := make([]byte, 65536)
    pktCount = 0

    var pcapw *pcapgo.Writer
    if pcapOutfile != "" {
        f, err := os.Create(pcapOutfile) 
        CheckError(err)
        defer f.Close()
        pcapw = pcapgo.NewWriter(f)
        err = pcapw.WriteFileHeader(65536, layers.LinkTypeIEEE802_11)
        CheckError(err)

    }

    

    
    for {
        n,addr,err := ServerConn.ReadFromUDP(buf)
        pktCount++
        log.Printf("Received len %v from addr %v [%d pkts]\n:", n, addr, pktCount)

        if (pcapOutfile != "") {
            pcapw.WritePacket(gopacket.CaptureInfo{
                Timestamp: time.Now(),
                CaptureLength: n,
                Length: n,
                InterfaceIndex: 0}, buf[0:n])
        }


        if (!dontInject) {
            // inject packets into interface
            werr := activePcapHandle.WritePacketData(buf[0:n])
            if werr != nil {
                log.Println("WritePacketData Error: ",werr)
            } 


        }
        // TODO: add some decoding summary?
        // TODO: log to pcap file
        // log.Printf("%s", hex.Dump(buf[0:n]))

        // inject it back into wireless interface
 
        if err != nil {
            log.Println("Error: ",err)
        } 
    }
}