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
	managePWMPin rpio.Pin
	isOpenKey    bool
	lock         bool = false
)

const (
	DebugLogPrefix        = "[DEBUG]"
	PwmPin                = 13
	VID            uint16 = 0x054C // SONY
	PID            uint16 = 0x06C1 // RC-S380
	Debug                 = true
)

func main() {
	log.Printf("%s /////// START OPEN KEY PROCESS ///////\n", DebugLogPrefix)

	initialize()

	initializeRestApiServer()

	for {
		// sudoしないと動かないので注意
		idm, err := pasori.GetID(VID, PID)
		if err != nil {
		}
		if Contains(userData, idm) {
			if isOpenKey {
				CloseKey()
			} else {
				OpenKey()
			}
		}
		if isOpenKey {
			CloseKey()
		} else {
			OpenKey()
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func initialize() {
	log.Printf("%s -: Initializing -----\n", DebugLogPrefix)

	////////////////// SERVO

	fmt.Println("-: -: Servo setup...")
	err := rpio.Open()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	managePWMPin = rpio.Pin(PwmPin) // SEIGYO OUT PUT PIN
	managePWMPin.Mode(rpio.Pwm)
	managePWMPin.Freq(50 * 100)
	managePWMPin.DutyCycle(0, 100)
	fmt.Println("-: -: END Servo setup")

	////////////////// PASORI

	fmt.Println("-: -: IDM Read setup...")

	// 登録されているIDM読み取り処理
	ReadUserData()

	fmt.Println("-: -: END IDM Read setup")

}

func OpenKey() {
	lock = true
	for i := 1; i <= 60; i++ {
		managePWMPin.DutyCycle(uint32(i), 100)
		time.Sleep(10 * time.Millisecond)
	}
	isOpenKey = true
}

func CloseKey() {
	for i := 1; i <= 60; i++ {
		managePWMPin.DutyCycle(uint32(50-i), 100)
		time.Sleep(10 * time.Millisecond)
	}
	isOpenKey = false
}

func initializeRestApiServer() {
	router := gin.Default()
	router.GET("/users", getUsers)
	router.POST("/users", postUsers)
	go router.Run("localhost:8080")
	//go router.Run("localhost:8080")
}

func getUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, userData)
}

func postUsers(c *gin.Context) {
	fmt.Println("Start Register User")
	idm, err := pasori.GetID(VID, PID)
	if err != nil {
		panic(err)
	}
	fmt.Println(idm)
	id, _ := uuid.NewUUID()
	user := User{
		ID:          id.String(),
		IDM:         idm,
		Name:        c.Params.ByName("name"),
		LastLogging: "",
		StNum:       "",
	}
	userData = append(userData, user)
	fmt.Println("End Register User")
	c.IndentedJSON(http.StatusCreated, user)
}
