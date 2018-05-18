package main

import "fmt"

func main() {
  fmt.Println(translate("s"))
  fmt.Println(translate("sss"))
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
    return "I don't recognise those honks"
  }
}
