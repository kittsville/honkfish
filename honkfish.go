package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/arbovm/levenshtein"
	"github.com/bradfitz/slice"
)

type slackResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

type alphabetically []string

func (s alphabetically) Len() int           { return len(s) }
func (s alphabetically) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s alphabetically) Less(i, j int) bool { return len(s[i]) < len(s[j]) }

var version = "1.5-BIRTHDAY-EDITION.1"

/*
	Translation map from honks to the boat's behaviour
*/
var optionsWithDesc = map[string]string{
	"honk":                                "I am altering my course to STARBOARD",
	"honk honk":                           "I am altering my course to PORT",
	"honk honk honk":                      "I am going ASTERN",
	"honk honk honk honk pause honk":      "I am turning through 180 degrees to STARBOARD",
	"honk honk honk honk pause honk honk": "I am turning through 180 degrees to PORT",
	"honk honk honk honk honk":            "I do not understand your intentions, *keep clear*, I doubt whether you are taking sufficient action to avoid a collision",
	"HONK": "I am about to get underway, enter the fairway or I am approaching a blind bend",
	"HONK pause honk pause honk":            "I am unable to manoeuvre - not under command",
	"HONK pause HONK pause honk":            "I intend to overtake you on YOUR STARBOARD side",
	"HONK pause HONK pause honk pause honk": "I intend to overtake you on YOUR PORT side",
	"HONK pause honk pause HONK pause honk": "I agree to be overtaken",
	"pause pause pause pause HONK HONK":     "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
	"HONK HONK":                             "Happy Birthday Dave!",
}

// A slice of the keys from `optionsWithDesc` so we can do operations with just
// those.
var options = []string{}

// Init is called before main, used for initialisers in Go.
func init() {
	for k := range optionsWithDesc {
		options = append(options, k)
	}
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
	case "list":
		sendResponse(w, honksList())
	case "all":
		sendResponse(w, honksList())
	case "version":
		sendResponse(w, botVersion())
	default:
		sendResponse(w, translateHonks(userInput))
	}
}

func translateHonks(honks string) slackResponse {
	translation, found := optionsWithDesc[honks]

	if !found {
		// If couldn't find the honk that was asked for exactly, let's use
		// Levenstien distance to suggest the closest two matches.
		optionsSortedByDistance := options[:]
		slice.Sort(optionsSortedByDistance, func(i, j int) bool {
			return levenshtein.Distance(options[i], honks) < levenshtein.Distance(options[j], honks)
		})

		// There will always be at least two in here, so it's safe for us to just
		// assume the indexes exist.
		translation = fmt.Sprintf("Failed to convert honks. Perhaps you misheard? Did you mean '%s' or '%s'?", optionsSortedByDistance[0], optionsSortedByDistance[1])
	}

	response := slackResponse{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("Translation: %s", translation),
	}

	return response
}

func usageText() slackResponse {
	response := slackResponse{
		ResponseType: "ephemeral",
		Text:         "Usage:\n`/honkfish honk pause HONK`\nhonk = short honk\nHONK = long honk\npause = a gap between honks",
	}

	return response
}

func honksList() slackResponse {
	formattedHonks := "Honks:"

	var honks []string

	for key := range optionsWithDesc {
		honks = append(honks, key)
	}

	sort.Sort(alphabetically(honks))

	for _, honk := range honks {
		if honk != "pause pause pause pause HONK HONK" {
			formattedHonks += fmt.Sprintf("\n_%s_ -> %s", honk, optionsWithDesc[honk])
		}
	}

	response := slackResponse{
		ResponseType: "ephemeral",
		Text:         formattedHonks,
	}

	return response
}

func botVersion() slackResponse {
	response := slackResponse{
		ResponseType: "ephemeral",
		Text:         "I am Honkfish Version " + version,
	}

	return response
}

// Handles JSON marshalling and sending response to client
func sendResponse(w http.ResponseWriter, response slackResponse) {
	j, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
