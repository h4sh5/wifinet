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
    handle *pcap.Handle
)


func handlePacket(packet gopacket.Packet) {
    fmt.Println(packet)
}

func startCapture(iface string, filter string) {
    ihandle, err := pcap.NewInactiveHandle(iface)
    if err != nil {
      log.Fatal(err)
    }
    defer ihandle.CleanUp()

    // put into monitor mode
    err = ihandle.SetRFMon(true)
    if err != nil {
      log.Fatal(err)
    }

    if err = ihandle.SetSnapLen(65536); err != nil {
        log.Fatal(err)
    }

    readTimeout := 500 * time.Millisecond
    if err = ihandle.SetTimeout(readTimeout); err != nil {
        log.Fatal(err)
    }

    // activate handle
    handle, err = ihandle.Activate()
    if err != nil {
      log.Fatal(err)
    }

    // set bpf and start capture
    if err := handle.SetBPFFilter(filter); err != nil {
        panic(err)
    } 


    fmt.Printf("starting capture on %s with bpf '%s'\n", iface, 
        filter)

    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

    for packet := range packetSource.Packets() {
        log.Printf("new packet: %v", packet)
        handlePacket(packet)  // Do something with a packet here.
    }



}


func main() {
 
    var iface string    
    var bpfFilter string
 
    // flags declaration using flag package
    flag.StringVar(&iface, "i", "", "Specify interface")
    flag.StringVar(&bpfFilter, "f", "", "Specify BPF filter")
    // TODO: add socket (IP and port)
    


    flag.Parse()
    
    if iface == "" {
        log.Printf("Need to specify wireless interface e.g. wlan0\n");
        log.Printf("usage: ")
        flag.PrintDefaults()
        
        os.Exit(1);
    }

    startCapture(iface, bpfFilter)


}