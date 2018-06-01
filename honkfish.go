package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type SlackResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

/*
	Translation map from honks to the boat's behaviour
	s = short honk
	l = long honk
	p = pause between honks
*/
var dictionary = map[string]string{
	"honk" 																	: "I am altering my course to STARBOARD",
	"honk honk" 														: "I am altering my course to PORT",
	"honk honk honk" 												: "I am going ASTERN",
	"honk honk honk honk pause honk" 				: "I am turning through 360 degrees to STARBOARD",
	"honk honk honk honk pause honk honk" 	: "I am turning through 360 degrees to PORT",
	"honk honk honk honk honk" 							: "I do not understand your intentions, *keep clear*, I doubt whether you are taking sufficient action to avoid a collision",
	"HONK" 																	: "I am about to get underway, enter the fairway or I am approaching a blind bend",
	"HONK pause honk pause honk" 						: "I am unable to manoeuvre - not under command",
	"HONK pause HONK pause honk" 						: "I intend to overtake you on YOUR STARBOARD side",
	"HONK pause HONK pause honk pause honk" : "I intend to overtake you on YOUR PORT side",
	"HONK pause honk pause HONK pause honk" : "I agree to be overtaken",
}

func main() {
	http.HandleFunc("/", requestHandler)
	port, customPort := os.LookupEnv("PORT")
	if !customPort {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userInput := r.Form.Get("text")

	switch userInput {
	case "help":
		sendResponse(w, usageText())
	case "usage":
		sendResponse(w, usageText())
	default:
		sendResponse(w, translateHonks(userInput))
	}
}

func translateHonks(honks string) SlackResponse {
	translation, found := dictionary[honks]

	if !found {
		translation = "Failed to convert honks. Perhaps you misheard?"
	}

	response := SlackResponse{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("Translation: %s", translation),
	}

	return response
}

func usageText() SlackResponse {
	response := SlackResponse{
		ResponseType: "ephemeral",
		Text:         "Usage:\n`/honkfish honk pause HONK`\nhonk = short honk\nHONK = long honk\npause = a gap between honks",
	}

	return response
}

// Handles JSON marshalling and sending response to client
func sendResponse(w http.ResponseWriter, response SlackResponse) {
	j, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
