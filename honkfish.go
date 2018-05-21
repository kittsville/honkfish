package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	rawHonks := r.Form.Get("text")

	if rawHonks == "help" || rawHonks == "usage" {
		fmt.Fprintf(w, "Usage:\n`/honkfish honk pause HONK`\nhonk = short honk\nHONK = long honk\npause = a gap between honks")
		return
	}

	honks := strings.Replace(rawHonks, "honk", "s", -1)
	honks = strings.Replace(honks, "HONK", "l", -1)
	honks = strings.Replace(honks, "pause", "p", -1)
	honks = strings.Replace(honks, " ", "", -1)
	translation := translate(honks)
	fmt.Fprintf(w, "Boat: %s\nTranslation: %s", rawHonks, translation)
}

func main() {
	http.HandleFunc("/", handler)
	port, customPort := os.LookupEnv("PORT")
	if !customPort {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func translate(honks string) string {
	honkTranslations := make(map[string]string)

	/*
	  Translation map from honks to the boat's behaviour
	  s = short honk
	  l = long honk
	  p = pause between honks
	*/
	honkTranslations["s"] = "I am altering my course to STARBOARD"
	honkTranslations["ss"] = "I am altering my course to PORT"
	honkTranslations["sss"] = "I am going ASTERN"
	honkTranslations["ssssps"] = "I am turning through 360 degrees to STARBOARD"
	honkTranslations["sssspss"] = "I am turning through 360 degrees to PORT"
	honkTranslations["sssss"] = "I do not understand your intentions, *keep clear*, I doubt whether you are taking sufficient action to avoid a collision"
	honkTranslations["l"] = "I am about to get underway, enter the fairway or I am approaching a blind bend"
	honkTranslations["lpsps"] = "I am unable to manoeuvre - not under command"
	honkTranslations["lplps"] = "I intend to overtake you on YOUR STARBOARD side"
	honkTranslations["lplpsps"] = "I intend to overtake you on YOUR PORT side"
	honkTranslations["lpsplps"] = "I agree to be overtaken"

	if translation, found := honkTranslations[honks]; found {
		return translation
	} else {
		return "Failed to convert honks. Perhaps you misheard?"
	}
}
