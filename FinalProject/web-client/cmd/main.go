package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-client/internal/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Print("No arguments were supplied")
		os.Exit(1)
	}

	command := os.Args[1]
	if command == "target" {
		targetServer()
	} else {
		commandsWithHostURL(command)
	}

}

func commandsWithHostURL(command string) {
	hostURL, err := retrieveHostURL()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	switch command {
	case "login":
		commands.Login(hostURL)
	case "register":
		commands.RegisterUser(hostURL)
	default:
		commandsWithAuth(command, hostURL)
	}
}

func commandsWithAuth(command, hostURL string) {
	token, err := retrieveJwtToken()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	switch command {
	case "logout":
		commands.Logout()
	case "create-group":
		commands.CreateGroup(hostURL, token)
	case "delete-group":
		commands.DeleteGroup(hostURL, token)
	case "add-member":
		commands.AddMember(hostURL, token)
	case "remove-member":
		commands.RemoveMember(hostURL, token)
	case "upload-file":
		commands.UploadFile(hostURL, token)
	case "download-file":
		commands.DownloadFile(hostURL, token)
	case "delete-file":
		commands.DeleteFile(hostURL, token)
	case "show-all-files":
		commands.ShowAllFilesInGroup(hostURL, token)
	case "show-all-groups":
		commands.ShowAllGroups(hostURL, token)
	case "show-all-users":
		commands.ShowAllUsers(hostURL, token)
	case "show-all-members":
		commands.ShowAllMembers(hostURL, token)
	default:
		fmt.Printf("Invalid command [%s]\n", command)
	}
}

func targetServer() {
	targetServerCommand := flag.NewFlagSet("target", flag.ExitOnError)
	hostURL := targetServerCommand.String("host", "", "Host url of the server")

	targetServerCommand.Parse(os.Args[2:])

	ioutil.WriteFile("/tmp/host", []byte(*hostURL), 0644)
	fmt.Printf("Server with host url [%s] was successfully targeted\n", *hostURL)
}

func retrieveJwtToken() (string, error) {
	if _, err := os.Stat("/tmp/jwt"); os.IsNotExist(err) {
		return "", errors.New("Problem with authentication, please login again")
	}

	token, err := ioutil.ReadFile("/tmp/jwt")
	if string(token) == "" || err != nil {
		return "", errors.New("Problem with authentication, please login again")
	}

	return string(token), nil
}

func retrieveHostURL() (string, error) {
	if _, err := os.Stat("/tmp/host"); os.IsNotExist(err) {
		return "", errors.New("Problem with host url, please target the server again")
	}

	url, err := ioutil.ReadFile("/tmp/host")
	if string(url) == "" || err != nil {
		return "", errors.New("Problem with host url, please target the server again")
	}

	return string(url), nil
}
