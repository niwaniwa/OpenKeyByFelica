package main

import (
	"fmt"
	"github.com/bamchoh/pasori"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stianeikeland/go-rpio/v4"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	managePWMPin    rpio.Pin
	manageMosPin    rpio.Pin
	manageSwPin     rpio.Pin
	isOpenKey       bool
	lock            bool   = false
	isRegister      bool   = false
	isCloseProgress bool   = false
	tempName        string = ""
)

const (
	DebugLogPrefix        = "[DEBUG]"
	PwmPin                = 13
	MosPin                = 17
	SwPin                 = 18
	VID            uint16 = 0x054C // SONY
	PID            uint16 = 0x06C1 // RC-S380
	Debug                 = true
)

func main() {
	log.Printf("%s /////// START OPEN KEY PROCESS ///////\n", DebugLogPrefix)

	initialize()

	initializeRestApiServer()

	go checkDoorState()

	for {
		// sudoしないと動かないので注意
		idm, err := pasori.GetID(VID, PID)
		if err != nil {
			continue
		}

		log.Println(idm)

		if isRegister {
			log.Println("Start Register User #" + tempName)
			id, _ := uuid.NewUUID()
			user := User{
				ID:          id.String(),
				IDM:         idm,
				Name:        tempName,
				LastLogging: "",
				StNum:       "",
			}
			SaveUserData(user)
			userData = append(userData, user)
			log.Println("End Register User")
			isRegister = false
			time.Sleep(5000 * time.Millisecond)
			continue
		}

		if Contains(userData, idm) {
			if isOpenKey {
				CloseKey()
			} else {
				OpenKey()
			}
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

func initialize() {
	log.Printf("%s -: Initializing -----\n", DebugLogPrefix)

	////////////////// SERVO

	_ = os.MkdirAll("data", 0755)

	fmt.Println("-: -: Servo setup...")
	err := rpio.Open()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	manageMosPin = rpio.Pin(MosPin) // MOS SEIGYO OUT PUT PIN
	manageMosPin.Output()
	manageMosPin.Low()
	managePWMPin = rpio.Pin(PwmPin) // SEIGYO OUT PUT PIN
	managePWMPin.Mode(rpio.Pwm)
	managePWMPin.Freq(50 * 100)
	managePWMPin.DutyCycle(0, 100)
	managePWMPin.Low()
	fmt.Println("-: -: END Servo setup")

	////////////////// SWITCH
	fmt.Println("-: -: switch setup...")

	manageSwPin = rpio.Pin(SwPin)
	manageSwPin.Input()
	manageSwPin.PullUp()
	//manageSwPin.Detect(rpio.FallEdge)

	////////////////// PASORI
	fmt.Println("-: -: IDM Read setup...")

	// 登録されているIDM読み取り処理
	ReadUserData()

	fmt.Println("-: -: END IDM Read setup")

}

func OpenKey() {
	lock = true
	manageMosPin.High()
	managePWMPin.High()

	time.Sleep(500 * time.Millisecond)

	for i := 1; i <= 60; i++ {
		managePWMPin.DutyCycle(uint32(i), 100)
		time.Sleep(10 * time.Millisecond)
	}

	go func() {
		time.Sleep(1000 * time.Millisecond)
		managePWMPin.Low()
		manageMosPin.Low()
	}()
	isOpenKey = true
}

func CloseKey() {
	manageMosPin.High()
	managePWMPin.High()
	time.Sleep(500 * time.Millisecond)
	for i := 1; i <= 60; i++ {
		managePWMPin.DutyCycle(uint32(50-i), 100)
		time.Sleep(10 * time.Millisecond)
	}
	go func() {
		time.Sleep(1000 * time.Millisecond)
		managePWMPin.Low()
		manageMosPin.Low()
	}()
	isOpenKey = false
}

func initializeRestApiServer() {
	router := gin.Default()
	router.GET("/user", getUser)
	router.POST("/user", postUser)
	go router.Run("localhost:8080")
	//go router.Run("localhost:8080")
}

func getUser(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, userData)
}

func postUser(c *gin.Context) {
	isRegister = true
	tempName = c.Params.ByName("name")
	c.IndentedJSON(http.StatusOK, tempName)
}

func checkDoorState() {
	for {
		if manageSwPin.Read() == 0 {
			if !isCloseProgress {
				isCloseProgress = true
				time.AfterFunc(5*time.Second, func() {
					isCloseProgress = false
					if manageSwPin.Read() == 0 {
						if isOpenKey {
							CloseKey()
							log.Println(":: Closed Door")
						}
					}
				})
			}
		}
		time.Sleep(1 * time.Second)
	}
}
