package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/groob/CommPhone/commportal"
)

var (
	line     bool
	all      bool
	Version  = "unreleased"
	fVersion = flag.Bool("version", false, "display the version")
)

func init() {
	flag.BoolVar(&line, "line", false, "Print direct line instead of Extension")
	flag.BoolVar(&all, "all", false, "Print all extensions")
	flag.Parse()
}

func printAll(s commportal.Subscribers) {
	for _, phone := range s {
		fmt.Printf("%v %v\n", phone.Name.Name, phone.IntercomDialingCode.Extension)
	}
}

func printPhone(name string, s commportal.Subscribers) {
	phones := s.Find(name)
	for _, phone := range phones {
		if line {
			fmt.Printf("%v %v\n", phone.Name.Name, phone.DirectoryNumber.Line)
		} else {
			fmt.Printf("%v %v\n", phone.Name.Name, phone.IntercomDialingCode.Extension)
		}
	}
}

func main() {
	// Print version
	if *fVersion {
		fmt.Printf("CommPortal Web Scraper - Version %s\n", Version)
		return
	}
	username := os.Getenv("COMMPORTAL_USER")
	password := os.Getenv("COMMPORTAL_PASS")
	loginURL := os.Getenv("COMMPORTAL_URL")
	portal, err := commportal.New(username, password, loginURL)
	if err != nil {
		log.Fatal(err)
	}
	phones, err := portal.Phones()
	if err != nil {
		log.Fatal(err)
	}

	// print all extensions
	if all {
		printAll(phones)
	}

	arg := flag.Args()
	printPhone(strings.Join(arg, " "), phones)
}
