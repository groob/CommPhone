package commportal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type config struct {
	username string
	password string
	session  string
	loginURL string
}

func (c config) Phones() (Subscribers, error) {
	return phones(c.session, c.loginURL)
}

type Subscribers []Subscriber

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

func (s Subscribers) Find(name string) []Subscriber {
	phones := s.findByFull(name)
	if len(phones) >= 1 {
		return phones
	}
	return s.findByName(name)

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
		}
	}
	return phones
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func phones(session string, url string) (Subscribers, error) {
	type Phones struct {
		Subscribers []Subscriber `json:"Subscriber,omitempty"`
	}
	var list Phones
	jsnURL := fmt.Sprintf("%v/%v/data?version=8.1.04&callback=dataObjectManager.callback&data=Msrb_BusinessGroup_ChildrenList_Subscriber", url, session)

	resp, err := http.Get(jsnURL)
	if err != nil {
		return nil, err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r, _ := regexp.Compile("\t(.*)}")
	str := r.FindString(string(buf))
	jsnStr := bytes.NewBufferString(str)
	err = json.NewDecoder(jsnStr).Decode(&list)
	if err != nil {
		return nil, err
	}
	return list.Subscribers, nil
}

func (c *config) getSession() error {
	// Create the Login Form
	loginURL := c.loginURL + "/login"
	v := url.Values{}
	v.Set("version", "8.1.04")
	v.Add("ApplicationID", "MS_WebClient")
	v.Add("ContextInfo", "version=8.1.04")
	v.Add("DirectoryNumber", c.username)
	v.Add("Password", c.password)
	v.Encode()
	// Post form
	resp, err := http.PostForm(loginURL, v)
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
	c.session = strings.Join(strings.Split(string(buf), "="), "")
	return nil

}

func New(username, password, url string) (*config, error) {
	config := &config{
		username: username,
		password: password,
		loginURL: url,
	}
	return config, config.getSession()
}
