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
  s       = short honk
  l       = long honk
  <space> = gap between honks
  */
  honkTranslations["s"] = "I am altering my course to STARBOARD"
  honkTranslations["ss"] = "I am altering my course to PORT"
  honkTranslations["sss"] = "I am going ASTERN"
  honkTranslations["ssss s"] = "I am turning through 360 degrees to STARBOARD"
  honkTranslations["ssss ss"] = "I am turning through 360 degrees to PORT"
  honkTranslations["sssss"] = "I do not understand your intentions, *keep clear*, I doubt whether you are taking sufficient action to avoid a collision"
  honkTranslations["l"] = "I am about to get underway, enter the fairway or I am approaching a blind bend"
  honkTranslations["l s s"] = "I am unable to manoeuvre - not under command"
  honkTranslations["l l s"] = "I intend to overtake you on YOUR STARBOARD side"
  honkTranslations["l l s s"] = "I intend to overtake you on YOUR PORT side"
  honkTranslations["l s l s"] = "I agree to be overtaken"

  if translation, found := honkTranslations[honks]; found {
    return translation
  } else {
    return "I don't recognise those honks"
  }
}
