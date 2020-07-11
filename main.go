package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
	"gobot.io/x/gobot/platforms/mqtt"
	"time"
)

func main() {

	var trigger bool
	var pomodoroCompleted bool
	var pomodoroCount int
	var message []byte

	firmataAdaptor := firmata.NewTCPAdaptor("192.168.3.108:3030")

	mqttAdaptor := mqtt.NewAdaptor("tcp://192.168.3.12:1883", "GOBOT")
	mqttAdaptor.SetAutoReconnect(true)

	button := gpio.NewButtonDriver(firmataAdaptor, "5")
	led := gpio.NewLedDriver(firmataAdaptor, "2")

	work := func() {

		led.Off()

		button.On(gpio.ButtonRelease, func(data interface{}) {
			led.On()
			fmt.Println("Button released")

			trigger = false

			start := time.Now()
			t := time.Now()

			fmt.Println("PomodoroCount:", pomodoroCount)
			if pomodoroCount != 0 && pomodoroCount < 4 && pomodoroCompleted {
			break5:
				for t.Sub(start) <= 6*time.Minute && !trigger {
					fmt.Println("Break elapsed time: ", t.Sub(start))
					if t.Sub(start) > 5*time.Minute {
						message = []byte("Bring your Ass back")
						mqttAdaptor.Publish("butt/pomodoro", message)
						fmt.Println("Breaking break5 at: ", t.Sub(start))
						break break5
					}
					t = time.Now()
				}
				fmt.Println("Outside break5 loop")
			} else if pomodoroCount != 0 && pomodoroCount >= 4 && pomodoroCompleted {
			break30:
				for t.Sub(start) <= 31*time.Minute && !trigger {
					fmt.Println("Break elapsed time: ", t.Sub(start))
					if t.Sub(start) > 30*time.Minute {
						message = []byte("Bring your Ass back")
						mqttAdaptor.Publish("butt/pomodoro", message)
						fmt.Println("Breaking break30 at: ", t.Sub(start))
						pomodoroCount = 0
						break break30
					}
					t = time.Now()
				}
				fmt.Println("Outside break30 loop")
			} else {
				fmt.Println("Pomodoro not completed")
				message = []byte("Bring your Ass back")
				mqttAdaptor.Publish("butt/pomodoro", message)
			}

		})

		button.On(gpio.ButtonPush, func(data interface{}) {
			led.Off()
			fmt.Println("Button pushed")

			trigger = true
			pomodoroCompleted = false

			start := time.Now()
			t := time.Now()

			message = []byte("Started")
			mqttAdaptor.Publish("butt/pomodoro", message)

		pomodoro25:
			for t.Sub(start) <= 26*time.Minute && trigger {
				fmt.Println("Pomodoro elapsed time: ", t.Sub(start))
				if t.Sub(start) > 25*time.Minute {
					message = []byte("Move your Ass")
					mqttAdaptor.Publish("butt/pomodoro", message)
					fmt.Println("Breaking pomodoro25 at: ", t.Sub(start))
					pomodoroCount++
					pomodoroCompleted = true
					break pomodoro25
				}
				t = time.Now()
			}

			fmt.Println("Outside pomodoro25 loop")

		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor, mqttAdaptor},
		[]gobot.Device{button, led},
		work,
	)

	// TODO Handle butt wiggle i.e. unwanted butt triggers
	// Experimented with separate goroutines for timers which can be started/stopped
	// with channels from button push and release but it results in unnecessary race conditions
	// and deadlocks; I guess the reason could be the way Gobot handles button.On() goroutines.

	// TODO Handle connection loss with the adaptors
	// Gobot API doesn't seem to offer mechanism to handle loss of connection
	// Issue raised - https://github.com/hybridgroup/gobot/issues/758
	/*	buttpomodoro.Every(1*time.Second, func() {
				firmataConnectionError := firmataAdaptor.Connect()
		        // Doesn't work, connection error is same even when connection is lost
				if !strings.Contains(firmataConnectionError.Error(), "client is already connected") {
					fmt.Println("Error connecting to Firmata adaptor")

					message = []byte("Error connecting to Firmata adaptor")
					mqttAdaptor.Publish("butt/pomodoro", message)
				}
				mqttConnectionError := mqttAdaptor.Connect()
				if mqttConnectionError != nil {
					fmt.Println("Error connecting to MQTT adaptor")
				}
	     })*/

	robot.Start()
}
