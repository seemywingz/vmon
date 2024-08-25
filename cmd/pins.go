package cmd

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/seemywingz/gotoolbox/gtb"
	"github.com/spf13/viper"
	"github.com/stianeikeland/go-rpio/v4"
)

type Pin struct {
	Name string
	On   bool
	Num  int
	Pin  *rpio.Pin
}

var pins = make(map[string]Pin)

func initPins() {
	pinConfigs := viper.Get("pins").([]interface{})

	for _, config := range pinConfigs {
		// Assert the config to the appropriate type (map[string]interface{})
		pinConfig := config.(map[string]interface{})

		// Extract the pin number, name, and on status
		pinNum := int(pinConfig["pin"].(float64)) // Viper may interpret numbers as float64
		name := pinConfig["name"].(string)
		on := pinConfig["on"].(bool)

		// Create a new GPIO pin for the switch
		pin, err := NewGPIOPin(pinNum, rpio.Output)
		if err != nil {
			gtb.EoE(err) // Handle error gracefully
			continue
		}

		// Create the Switch object and append it to the switches slice
		p := Pin{
			Name: name,
			On:   on,
			Num:  pinNum,
			Pin:  pin,
		}

		if p.On {
			p.Pin.High()
		} else {
			p.Pin.Low()
		}

		pins[name] = p
		log.Printf("Initialized switch: %s, on: %v, Pin: %d\n", p.Name, p.On, p.Num)
	}

}

func handlePin(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		// Parse the request body
		var req struct {
			Name string `json:"name"`
			On   bool   `json:"on"`
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
			log.Printf("Failed to parse JSON: %v", err)
			return
		}

		// Find the switch by name
		p, ok := pins[req.Name]
		if !ok {
			http.Error(w, "Switch not found", http.StatusNotFound)
			log.Printf("Switch not found: %s", req.Name)
			return
		}

		// Toggle the switch
		if req.On {
			p.Pin.High()
		} else {
			p.Pin.Low()
		}

		p.On = req.On
		pins[req.Name] = p

		// format the config to save to the config file
		var pinConfigs []interface{}
		for _, pin := range pins {
			pinConfigs = append(pinConfigs, map[string]interface{}{
				"name": pin.Name,
				"on":   pin.On,
				"pin":  pin.Num,
			})
		}

		// Update the config file
		viper.Set("pins", pinConfigs)
		viper.WriteConfig()

		// Write a response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Switch toggled"))
	}
}
