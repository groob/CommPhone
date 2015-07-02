package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type Config struct {
	Username string
	Password string
	Session  string
}

type Subscriber struct {
	AccountLocation struct {
		_ string `json:"_"`
	} `json:"AccountLocation,omitempty"`
	AdditionalNumber struct {
		_ bool `json:"_"`
	} `json:"AdditionalNumber,omitempty"`
	AttendantLineType struct {
		_ string `json:"_"`
	} `json:"AttendantLineType,omitempty"`
	BusinessGroupAdministrationAccountType struct {
		_ string `json:"_"`
	} `json:"BusinessGroupAdministrationAccountType,omitempty"`
	BusinessGroupAdministrationAdministrationDepartment struct {
		_ string `json:"_"`
	} `json:"BusinessGroupAdministrationAdministrationDepartment,omitempty"`
	Department struct {
		Department string `json:"_"`
	} `json:"Department,omitempty"`
	DirectoryNumber struct {
		Line string `json:"_"`
	} `json:"DirectoryNumber,omitempty"`
	IntercomDialingCode struct {
		Extension string `json:"_"`
	} `json:"IntercomDialingCode,omitempty"`
	Name struct {
		Name string `json:"_"`
	} `json:"Name,omitempty"`
	SpokenNameRecorded struct {
		SpNameRecorded bool `json:"_"`
	} `json:"SpokenNameRecorded,omitempty"`
	SubscriberType struct {
		_ string `json:"_"`
	} `json:"SubscriberType,omitempty"`
	SubscriberUID struct {
		_ string `json:"_"`
	} `json:"SubscriberUID,omitempty"`
}

type Subscribers []Subscriber

func phones() (Subscribers, error) {
	type Phones struct {
		Subscribers []Subscriber `json:"Subscriber,omitempty"`
	}
	var list Phones
	resp, err := http.Get(fmt.Sprintf("http://optimumlightpathvoice.com/%v/bg/data?version=8.1.04&callback=dataObjectManager.callback&data=Msrb_BusinessGroup_ChildrenList_Subscriber", config.Session))
	if err != nil {
		log.Fatal(err)
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	r, _ := regexp.Compile("\t(.*)}")
	str := r.FindString(string(buf))
	jsnStr := bytes.NewBufferString(str)
	err = json.NewDecoder(jsnStr).Decode(&list)
	if err != nil {
		log.Fatal(err)
	}
	return list.Subscribers, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (s Subscribers) printAll() {
	for _, phone := range s {
		fmt.Printf("%v %v\n", phone.Name.Name, phone.IntercomDialingCode.Extension)
	}
}

func (s Subscribers) findByFull(name string) []Subscriber {
	var phones []Subscriber
	for _, phone := range s {
		if phone.Name.Name == name {
			phones = append(phones, phone)
		}
	}
	return phones
}

func (s Subscribers) findByName(name string) []Subscriber {
	var phones []Subscriber
	for _, phone := range s {
		if stringInSlice(name, strings.Split(phone.Name.Name, " ")) {
			phones = append(phones, phone)
			//fmt.Printf("%v %v\n", phone.Name.Name, phone.IntercomDialingCode.Extension)
		}
	}
	return phones
}

func (s Subscribers) printExtension(name string) {
	switch len(strings.Split(name, " ")) {
	case 1:
		phones := s.findByName(name)
		for _, phone := range phones {
			fmt.Printf("%v %v\n", phone.Name.Name, phone.IntercomDialingCode.Extension)
		}
	default:
		phones := s.findByFull(name)
		for _, phone := range phones {
			fmt.Printf("%v %v\n", phone.Name.Name, phone.IntercomDialingCode.Extension)
		}

	}
}

func (s Subscribers) printLine(name string) {
	switch len(strings.Split(name, " ")) {
	case 1:
		phones := s.findByName(name)
		for _, phone := range phones {
			fmt.Printf("%v %v\n", phone.Name.Name, phone.DirectoryNumber.Line)
		}
	default:
		phones := s.findByFull(name)
		for _, phone := range phones {
			fmt.Printf("%v %v\n", phone.Name.Name, phone.DirectoryNumber.Line)
		}

	}
}

var (
	line     bool
	all      bool
	config   *Config
	Version  = "unreleased"
	fVersion = flag.Bool("version", false, "display the version")
)

func init() {
	config = &Config{
		Username: os.Getenv("COMMPORTAL_USER"),
		Password: os.Getenv("COMMPORTAL_PASS"),
	}
	err := config.getSession()
	if err != nil {
		log.Fatal(err)
	}
	flag.BoolVar(&line, "line", false, "Print Direct Line instead of Extension")
	flag.BoolVar(&all, "all", false, "Print all Extensions")
	flag.Parse()
}

func (c *Config) getSession() error {
	// Create the Login Form
	v := url.Values{}
	v.Set("version", "8.1.04")
	v.Add("ApplicationID", "MS_WebClient")
	v.Add("ContextInfo", "version=8.1.04")
	v.Add("DirectoryNumber", config.Username)
	v.Add("Password", config.Password)
	v.Encode()
	// Post form
	resp, err := http.PostForm("http://optimumlightpathvoice.com/login", v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// resp should return something like session=SessionID
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// Set the session in the config
	c.Session = strings.Join(strings.Split(string(buf), "="), "")
	return nil

}

func main() {
	if *fVersion {
		fmt.Printf("CommPortal Web Scraper - Version %s\n", Version)
		return
	}
	phones, err := phones()
	if err != nil {
		log.Fatal(err)
	}
	arg := flag.Args()
	// print direct line or extension
	if line {
		phones.printLine(strings.Join(arg, " "))
	} else {
		phones.printExtension(strings.Join(arg, " "))
	}
	// print all extensions
	if all {
		phones.printAll()

	}
}
