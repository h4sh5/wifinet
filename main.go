package main

//     wifinet - tunnel wifi packets over the internet
//     Copyright (C) 2021 Haoxi Tan

//     This program is free software: you can redistribute it and/or modify
//     it under the terms of the GNU General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.

//     This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU General Public License for more details.

//     You should have received a copy of the GNU General Public License
//     along with this program.  If not, see <https://www.gnu.org/licenses/>.

 
import (
    "flag"
    "log"
    "fmt"
    "os"
    "time"
    "github.com/google/gopacket"
    "github.com/google/gopacket/pcap"
    "net"
)

var (
    activePcapHandle *pcap.Handle
    remoteConn net.Conn
    dontInject bool
    dontSniff bool
    pcapOutfile string
)

func sendFrame(conn net.Conn, data []byte) {
    /* send the wireless frame to the remote connection */
    conn.Write(data)
}


func handleWirelessPacket(packet gopacket.Packet) {
    // send to network socket 
    // fmt.Println(packet)
    if remoteConn != nil {
        sendFrame(remoteConn, packet.Data())
    }
}

/* A Simple function to verify error */
func CheckError(err error) {
    if err  != nil {
        log.Fatal(err)
    }
}

func setupWirelessIface(iface string, filter string, channels string) {

    inactiveHandle, err := pcap.NewInactiveHandle(iface)

    CheckError(err)
    defer inactiveHandle.CleanUp()

    // put into monitor modea
    err = inactiveHandle.SetRFMon(true)
    if err != nil {
        log.Printf("SetRFMon failed: %v\n", err)
        log.Printf("You might have to set this interface into monitor mode manually, such as using\nsudo iwconfig %s mode mon\n", iface)
        log.Println("If its already in monitor mode, don't worry.")
    }
    

    err = inactiveHandle.SetSnapLen(65536)
    CheckError(err)

    readTimeout := 500 * time.Millisecond
    inactiveHandle.SetTimeout(readTimeout)
    CheckError(err)

    SetupChannelHopping(iface, channels)

    // activate handle
    activePcapHandle, err = inactiveHandle.Activate()
    CheckError(err)

    // set bpf and start capture
    activePcapHandle.SetBPFFilter(filter)
    CheckError(err)

}

func captureSendLoop() {

    packetSource := gopacket.NewPacketSource(activePcapHandle, activePcapHandle.LinkType())

    for packet := range packetSource.Packets() {
        // log.Printf("new packet: %v", packet)
        go handleWirelessPacket(packet)  // Do something with a packet here.
    }

}


func main() {
 
    var iface string    
    var bpfFilter string
    var channels string
    var remoteSrv string
    var localSrv string
 
    // flags declaration using flag package
    flag.StringVar(&iface, "i", "", "Specify interface")
    flag.StringVar(&bpfFilter, "f", "", "Specify BPF filter")
    // TODO: multi channel support
    flag.StringVar(&channels, "c", "13", "Wireless channels to listen on, comma separated (only single channel supported for now. Default 13")
    flag.StringVar(&localSrv, "l", ":4141", "udp server:port to LISTEN on")
    flag.StringVar(&remoteSrv, "r", "", "host:port to connect to (via UDP)")
    flag.BoolVar(&dontInject, "ni", false, "do not inject packets (listen and forward only)")
    flag.BoolVar(&dontSniff, "ns", false, "do not sniff packets on wireless device")
    flag.StringVar(&pcapOutfile, "o", "", "pcapng file to log incoming packets") // TODO bi-directional packet logging

    

    flag.Parse()

    log.Printf("dontInject: %v, dontSniff: %v\n", dontInject, dontSniff)
    
    if iface == "" {
        log.Printf("Need to specify wireless interface e.g. wlan0\n");
        log.Printf("usage: ")
        flag.PrintDefaults()
        
        os.Exit(1);
    }

    setupWirelessIface(iface, bpfFilter, channels)

    if (dontSniff) {
        getFramesFromServer(localSrv, activePcapHandle, dontInject)
    } else {
        go getFramesFromServer(localSrv, activePcapHandle, dontInject)
    }
    

    // remote connection
    if remoteSrv != "" {
        serverAddr,err := net.ResolveUDPAddr("udp", remoteSrv)
        CheckError(err)
        remoteConn, err = net.DialUDP("udp", nil, serverAddr)
        CheckError(err)
    }

    if (!dontSniff) {
        fmt.Printf("starting capture on %s with bpf '%s'\n", iface, bpfFilter)
        captureSendLoop()    
    }
    

}