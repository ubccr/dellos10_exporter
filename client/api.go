package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/ini.v1"
)

var (
	config *ini.File
)

type OS10Command interface {
	GetCmdEndpoint() string
}

type Handler struct {
	commands []OS10Command
	target   string
	username string
	password string
}

func LoadConfig(path string) error {
	var err error
	config, err = ini.Load(path)
	if err != nil {
		return err
	}

	return nil
}

func GetHandle(target string) (*Handler, error) {
	section, err := config.GetSection("connection:" + target)
	if err != nil {
		return nil, fmt.Errorf("Connection not found in config: %s", target)
	}

	h := &Handler{}
	k, err := section.GetKey("host")
	if err != nil {
		return nil, fmt.Errorf("Connection missing host key: %s", target)
	}
	h.target = k.MustString("")

	k, err = section.GetKey("username")
	if err != nil {
		return nil, fmt.Errorf("Connection missing username key: %s", target)
	}
	h.username = k.MustString("")

	k, err = section.GetKey("password")
	if err != nil {
		return nil, fmt.Errorf("Connection missing password key: %s", target)
	}
	h.password = k.MustString("")

	return h, nil
}

func (h *Handler) AddCommand(cmd OS10Command) {
	h.commands = append(h.commands, cmd)
}

func (h *Handler) callCommand(cmd OS10Command) error {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/restconf%s", h.target, cmd.GetCmdEndpoint()), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(h.username, h.password)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("DellOS10 restconf api call failed with HTTP status code: %d", res.StatusCode)
	}

	rawJson, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(rawJson, &cmd)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) Call() error {
	for _, cmd := range h.commands {
		err := h.callCommand(cmd)
		if err != nil {
			return err
		}
	}

	return nil
}
