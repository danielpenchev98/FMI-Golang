package commands

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-client/internal/endpoints"
	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-client/internal/restclient"
	"github.com/jedib0t/go-pretty/v6/table"
)

type LoginResponse struct {
	Status int    `json:"status"`
	Token  string `json:"token"`
}

type CredentialsPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type UsersInfo struct {
	Status    uint       `json:"status"`
	UsersInfo []UserInfo `json:"users"`
}

func RegisterUser(hostURL string) {
	registrationCommand := flag.NewFlagSet("register", flag.ExitOnError)

	username := registrationCommand.String("usr", "", "username")
	password := registrationCommand.String("pass", "", "password")

	registrationCommand.Parse(os.Args[2:])

	if *username == "" && *password == "" {
		registrationCommand.PrintDefaults()
		return
	}

	rqBody := CredentialsPayload{
		Username: *username,
		Password: *password,
	}

	successBody := BasicResponse{}
	restClient := restclient.NewRestClientImpl("")
	url := hostURL + endpoints.RegisterAPIEndpoint
	err := restClient.Post(url, &rqBody, &successBody)

	if err != nil {
		fmt.Printf("Problem with the registration request. %s", err.Error())
		return
	}

	fmt.Println("User successfully created")
	fmt.Println(successBody.Status)
}

func Login(hostURL string) {
	loginCommand := flag.NewFlagSet("login", flag.ExitOnError)

	username := loginCommand.String("usr", "", "username")
	password := loginCommand.String("pass", "", "password")

	loginCommand.Parse(os.Args[2:])

	if *username == "" || *password == "" {
		loginCommand.PrintDefaults()
		return
	}

	rqBody := CredentialsPayload{
		Username: *username,
		Password: *password,
	}

	successBody := LoginResponse{}

	restClient := restclient.NewRestClientImpl("")
	url := hostURL + endpoints.LoginAPIEndpoint
	err := restClient.Post(url, &rqBody, &successBody)

	if err != nil {
		fmt.Printf("Problem with the login request. %s", err.Error())
		return
	}

	ioutil.WriteFile("/tmp/jwt", []byte(successBody.Token), 0644)
	fmt.Println("Login is successful")
}

func Logout() {
	if _, err := os.Stat("/tmp/jwt"); err != nil {
		os.Remove("/tmp/jwt")
	}
	fmt.Println("Logout successfull")
}

func ShowAllUsers(hostURL string, token string) {
	successBody := UsersInfo{}

	restClient := restclient.NewRestClientImpl(token)
	url := hostURL + endpoints.GetAllUsersAPIEndpoint
	err := restClient.Get(url, &successBody)

	if err != nil {
		fmt.Printf("Problem with the group creation request. %s", err.Error())
		return
	}

	tableRows := make([]table.Row, len(successBody.UsersInfo))
	for _, userInfo := range successBody.UsersInfo {
		tableRows = append(tableRows, table.Row{userInfo.ID, userInfo.Username})
	}
	PrintTable(table.Row{"ID", "Username"}, tableRows)
}

func ShowAllMembers(hostURL, token string) {
	getAllMembers := flag.NewFlagSet("show-all-members", flag.ExitOnError)
	groupName := getAllMembers.String("grp", "", "Name of the group")

	getAllMembers.Parse(os.Args[2:])

	if *groupName == "" {
		getAllMembers.PrintDefaults()
		os.Exit(1)
	}
	successBody := UsersInfo{}

	restClient := restclient.NewRestClientImpl(token)
	url := fmt.Sprintf("%s%s?group_name=%s", hostURL, endpoints.GetAllMembersAPIEndpoint, *groupName)
	err := restClient.Get(url, &successBody)

	if err != nil {
		fmt.Printf("Problem with the group creation request. %s", err.Error())
		return
	}

	tableRows := make([]table.Row, len(successBody.UsersInfo))
	for _, userInfo := range successBody.UsersInfo {
		tableRows = append(tableRows, table.Row{userInfo.ID, userInfo.Username})
	}
	PrintTable(table.Row{"ID", "Username"}, tableRows)
}
