package main

import (
	"fmt"
	"github.com/bamchoh/pasori"
	"github.com/stianeikeland/go-rpio/v4"
	"os"
	"time"
)

var (
	VID uint16 = 0x054C // SONY
	PID uint16 = 0x06C1 // RC-S380
)

const (
	PWM_PIN = 13
)

func main() {
	fmt.Println("// GPIO initializing")

	err := rpio.Open()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	managePin := rpio.Pin(PWM_PIN) // SEIGYO OUT PUT PIN
	managePin.Mode(rpio.Pwm)
	managePin.Freq(50 * 100)
	managePin.DutyCycle(0, 100)

	fmt.Println("// Start Felica read process")

	idm, err := pasori.GetID(VID, PID)
	if err != nil {
		panic(err)
	}
	fmt.Println(idm)
	fmt.Println("End Process")

	rpio.StartPwm()
	for i := 1; i <= 100; i++ {
		//fmt.Println(i)
		managePin.DutyCycle(uint32(i), 100)
		time.Sleep(100 * time.Millisecond)
	}

	for i := 1; i <= 100; i++ {
		//fmt.Println(i)
		managePin.DutyCycle(uint32(100-i), 100)
		time.Sleep(100 * time.Millisecond)
	}

}
