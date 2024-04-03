package main

import (
	"flag"
	"log"

	"github.com/chrisbsmith/goldfinger/bonds"
	"github.com/chrisbsmith/goldfinger/config"
)

func main() {

	// String that contains the configured configuration path
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "./config.yaml", "path to config file")

	// Actually parse the flags
	flag.Parse()

	// configPath := "config.yaml"
	c, err := config.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}
	Bonds, err := bonds.LoadBonds(c)
	if err != nil {
		panic(err)
	}

	v, e := Bonds.GetTotalValue()
	if e != nil {
		log.Printf("Error getting total value: %s", e.Error())
	}
	log.Printf("Total Value = $%.2f\n", v)

	i, e := Bonds.GetTotalInterest()
	if e != nil {
		log.Printf("Error getting total interest: %s", e.Error())
	}
	log.Printf("Total Interest = $%.2f\n", i)

	p, e := Bonds.GetTotalPurchasedPrice()
	if e != nil {
		log.Printf("Error getting total purchase price: %s", e.Error())
	}
	log.Printf("Total Original Purchase Price = $%.2f\n", p)

	bs := Bonds.FindUnmaturedBonds()
	log.Printf("Found %d bonds still to mature\n", len(bs))
	if len(bs) > 0 {
		for _, bss := range bs {
			log.Printf("Bond %s will mature on %s\n", bss.Serial, bss.FinalMaturity)
		}
	}
}
