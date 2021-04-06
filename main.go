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
)

var (
    activePcapHandle *pcap.Handle
)


func handleWirelessPacket(packet gopacket.Packet) {
    // send to network socket 
    fmt.Println(packet)
}

/* A Simple function to verify error */
func CheckError(err error) {
    if err  != nil {
        log.Fatal(err)
    }
}

func startCapture(iface string, filter string) {
    inactiveHandle, err := pcap.NewInactiveHandle(iface)
    CheckError(err)
    defer inactiveHandle.CleanUp()

    // put into monitor mode
    err = inactiveHandle.SetRFMon(true)
    CheckError(err)

    err = inactiveHandle.SetSnapLen(65536)
    CheckError(err)

    readTimeout := 500 * time.Millisecond
    inactiveHandle.SetTimeout(readTimeout)
    CheckError(err)

    // activate handle
    activePcapHandle, err = inactiveHandle.Activate()
    CheckError(err)

    // set bpf and start capture
    activePcapHandle.SetBPFFilter(filter)
    CheckError(err)


    fmt.Printf("starting capture on %s with bpf '%s'\n", iface, 
        filter)

    packetSource := gopacket.NewPacketSource(activePcapHandle, activePcapHandle.LinkType())

    for packet := range packetSource.Packets() {
        log.Printf("new packet: %v", packet)
        handleWirelessPacket(packet)  // Do something with a packet here.
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
    

    flag.Parse()
    
    if iface == "" {
        log.Printf("Need to specify wireless interface e.g. wlan0\n");
        log.Printf("usage: ")
        flag.PrintDefaults()
        
        os.Exit(1);
    }

    // start udp server
    go udp_listen(localSrv)

    startCapture(iface, bpfFilter)


}