package five

import (
   "runtime"
   "fmt"
   "github.com/kraman/go-firmata"
   "os"
   "os/exec"
   "strings"
)

const grepForDevice = "ls /dev | grep -iE 'usb|acm'"
var portName string
var usedSerialDevices map[string]bool

func init() {
	if runtime.GOOS == "darwin" {
		portName = "cu"
	} else {
		portName = "tty"
	}
	usedSerialDevices = make(map[string]bool)
}

func Serial(b extendedBoard) (client *firmata.FirmataClient, err error) {
	fmt.Println("Connecting...")
	data, err := exec.Command("bash", "-c", grepForDevice).Output()
	output := string(data)
	possiblePorts := strings.Split(output, "\n")
	availablePorts := make([]string, 0, len(possiblePorts))

	for _, value := range possiblePorts {
		if !strings.Contains(value, portName) {
			continue
		}
		if _, ok := usedSerialDevices[value]; ok {
			continue
		}

		availablePorts = append(availablePorts, value)
	}

	if len(availablePorts) == 0 {
		fmt.Printf("No USB ports found\n")
		os.Exit(3)
	}

	fmt.Printf("Available Ports: %s\n", availablePorts)

	usb := availablePorts[0]
	fmt.Printf("Using /dev/%s\n", usb)

	usedSerialDevices[usb] = true

	client, err = firmata.NewClient("/dev/" + usb, 57600)
	return

}