package main
// channel switching
// mostly inspired and copied from https://github.com/bettercap/bettercap/blob/master/network/
// especially net_wifi.go, net_linux.go, net_darwin.go etc.

// windows does not support wifi channel hopping.

import (
	"os/exec"
	"fmt"
	"log"
	"strconv"
	"runtime"
	"github.com/evilsocket/islazy/str"
	"strings"
)

// from bettercap core
func ExecCmd(executable string, args []string) (string, error) {
	path, err := exec.LookPath(executable)
	if err != nil {
		return "", err
	}

	raw, err := exec.Command(path, args...).CombinedOutput()
	if err != nil {
		return "", err
	} else {
		return str.Trim(string(raw)), nil
	}
}


const airPortPath = "/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport"


func PathHasBinary(executable string) bool {
	if path, err := exec.LookPath(executable); err != nil || path == "" {
		return false
	}
	return true
}


func SetInterfaceChannelDarwin(iface string, channel int) error {

	log.Printf("SetInterfaceChannelDarwin(%s,%d)", iface, channel)
	_, err := ExecCmd(airPortPath, []string{iface, fmt.Sprintf("-c%d", channel)})
	if err != nil {
		return err
	}

	return nil
}



func SetInterfaceChannelLinux(iface string, channel int) error {

	if PathHasBinary("iw") {
		log.Printf("SetInterfaceChannel(%s, %d) iw based", iface, channel)
		out, err := ExecCmd("iw", []string{"dev", iface, "set", "channel", fmt.Sprintf("%d", channel)})
		if err != nil {
			return err
		} else if out != "" {
			return fmt.Errorf("Unexpected output while setting interface %s to channel %d: %s", iface, channel, out)
		}
	} else if PathHasBinary("iwconfig") {
		log.Printf("SetInterfaceChannel(%s, %d) iwconfig based", iface, channel)

		out, err := ExecCmd("iwconfig", []string{iface, "channel", fmt.Sprintf("%d", channel)})
		if err != nil {
			return err
		} else if out != "" {
			return fmt.Errorf("Unexpected output while setting interface %s to channel %d: %s", iface, channel, out)
		}
	} else {
		return fmt.Errorf("no iw or iwconfig binaries found in $PATH")
	}

	return nil
}


func SetInterfaceChannel (iface string, channel int) {
	/* wrapper around the OS-dependent set channel functions */
	os := runtime.GOOS
    switch os {
    case "windows":
        log.Println("Channel hopping is not supported on Windows")
    case "darwin":
        SetInterfaceChannelDarwin(iface, channel)
    case "linux":
        SetInterfaceChannelLinux(iface, channel)
    default:
        log.Println("Unknown OS when trying to set interface")
    }


}

func SetupChannelHopping(iface string, channels string) {
	/* only single channel for now */
	chans := strings.Split(channels, ",")
	channel,err := strconv.Atoi(chans[0])
	CheckError(err)
	SetInterfaceChannel(iface, channel)
}

