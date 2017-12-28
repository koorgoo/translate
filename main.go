package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/koorgoo/translate/lib/translate"
)

var (
	key        = flag.String("key", "", "Yandex API key (defaults to $YANDEXTRANSLATEAPIKEY)")
	lang       = flag.String("to", translate.RU, "destination language (two- or three-letter code)")
	hint       = flag.String("from", "", "source language (two- or three-letter code); auto-detected by default")
	detect     = flag.Bool("lang", false, "detect source language")
	directions = flag.Bool("ls", false, "list translation directions")
)

func main() {
	flag.Parse()
	text := strings.Join(flag.Args(), " ")

	if len(text) == 0 && !*directions {
		flag.Usage()
		os.Exit(1)
	}
	if *key == "" {
		*key = os.Getenv("YANDEXTRANSLATEAPIKEY")
	}
	if *key == "" {
		log.Fatal("$YANDEXTRANSLATEAPIKEY not set")
	}

	c, err := translate.New(translate.Config{Key: *key})
	if err != nil {
		log.Fatal(err)
	}

	if *directions {
		v, _, err := c.GetLangs()
		if err != nil {
			log.Fatal(err)
		}
		for _, s := range v {
			fmt.Println(s)
		}
		return
	}

	if *detect {
		v, err := c.Detect(text)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(v)
		return
	}

	v, err := c.Translate(text, *lang, translate.From(*hint))
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range v {
		fmt.Println(s)
	}
}
