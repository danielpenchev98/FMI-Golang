package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-client/internal/commands"
)

/*
type CredentialsPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type BasicResponse struct {
	Status int `json:"status"`
}*/

/*
func registerUser(hostURL string) {
	registrationCommand := flag.NewFlagSet("register", flag.ExitOnError)

	username := registrationCommand.String("usr", "", "username")
	password := registrationCommand.String("pass", "", "password")

	registrationCommand.Parse(os.Args[2:])

	if *username == "" && *password == "" {
		flag.PrintDefaults()
		os.Exit(1)
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

type LoginResponse struct {
	Status int    `json:"status"`
	Token  string `json:"token"`
}

func login(hostURL string) {
	loginCommand := flag.NewFlagSet("login", flag.ExitOnError)

	username := loginCommand.String("usr", "", "username")
	password := loginCommand.String("pass", "", "password")

	loginCommand.Parse(os.Args[2:])

	if *username == "" || *password == "" {
		flag.PrintDefaults()
		os.Exit(1)
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

func logout() {
	if _, err := os.Stat("/tmp/jwt"); err != nil {
		os.Remove("/tmp/jwt")
	}
	fmt.Println("Logout successfull")
}*/

/*
type GroupPayload struct {
	GroupName string `json:"group_name"`
}

func createGroup(hostURL, token string) {
	createGroupCommand := flag.NewFlagSet("create-group", flag.ExitOnError)
	groupName := createGroupCommand.String("grp", "", "Name of the group to be created")

	createGroupCommand.Parse(os.Args[2:])
	if *groupName == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	rqBody := GroupPayload{
		GroupName: *groupName,
	}

	restClient := restclient.NewRestClientImpl(token)
	url := hostURL + endpoints.CreateGroupAPIEndpoint
	err := restClient.Post(url, &rqBody, nil)

	if err != nil {
		fmt.Printf("Problem with the group creation request. %s", err.Error())
		return
	}

	fmt.Printf("Group %s was succesfully created", *groupName)
}

func deleteGroup(hostURL, token string) {
	deleteGroupCommand := flag.NewFlagSet("delete-group", flag.ExitOnError)
	groupName := deleteGroupCommand.String("grp", "", "Name of the group to be deleted")
	deleteGroupCommand.Parse(os.Args[2:])

	if *groupName == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	rqBody := GroupPayload{
		GroupName: *groupName,
	}

	restClient := restclient.NewRestClientImpl(token)
	url := hostURL + endpoints.DeleteGroupAPIEndpoint
	err := restClient.Delete(url, &rqBody, nil)

	if err != nil {
		fmt.Printf("Problem with the group creation request. %s", err.Error())
		return
	}

	fmt.Printf("Group %s was succesfully deleted", *groupName)
}

type MembershipRequest struct {
	GroupPayload
	Username string `json:"username"`
}

func addUserToGroup(hostURL, token string) {
	addMemberCommand := flag.NewFlagSet("add-member", flag.ExitOnError)
	username := addMemberCommand.String("usr", "", "Name of the user to be added to the group")
	groupName := addMemberCommand.String("grp", "", "Name of the group to be deleted")
	addMemberCommand.Parse(os.Args[2:])

	if *groupName == "" || *username == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	rqBody := MembershipRequest{
		Username: *username,
	}
	rqBody.GroupName = *groupName

	restClient := restclient.NewRestClientImpl(token)
	url := hostURL + endpoints.AddMemberAPIEndpoint
	err := restClient.Post(url, &rqBody, nil)

	if err != nil {
		fmt.Printf("Problem with the group creation request. %s", err.Error())
		return
	}

	fmt.Printf("User %s was successfully added to group %s\n", *username, *groupName)
}

func removeUserFromGroup(hostURL, token string) {
	removeMemberCommand := flag.NewFlagSet("remove-member", flag.ExitOnError)

	username := removeMemberCommand.String("usr", "", "Name of the user to be added to the group")
	groupName := removeMemberCommand.String("grp", "", "Name of the group to be deleted")

	removeMemberCommand.Parse(os.Args[2:])
	if *groupName == "" || *username == "" {
		removeMemberCommand.PrintDefaults()
		return
	}

	rqBody := MembershipRequest{
		Username: *username,
	}
	rqBody.GroupName = *groupName

	restClient := restclient.NewRestClientImpl(token)
	url := hostURL + endpoints.RemoveMemberAPIEndpoint
	err := restClient.Delete(url, &rqBody, nil)

	if err != nil {
		fmt.Printf("Problem with the group creation request. %s", err.Error())
		return
	}

	fmt.Printf("User %s was successfully removed from group %s\n", *username, *groupName)
}*/

/*
type FileUploadResponse struct {
	FileID uint `json:"file_id"`
}

func uploadFile(hostURL, token string) {
	uploadFileCommand := flag.NewFlagSet("upload-file", flag.ExitOnError)
	filePath := uploadFileCommand.String("filepath", "", "Path to the file")
	groupName := uploadFileCommand.String("grp", "", "Name of the group, in which the file will be uploaded")

	uploadFileCommand.Parse(os.Args[2:])
	if *groupName == "" || *filePath == "" {
		uploadFileCommand.PrintDefaults()
		return
	}

	successBody := FileUploadResponse{}

	restClient := restclient.NewRestClientImpl(token)
	url := fmt.Sprintf("%s%s?group_name=%s", hostURL, endpoints.UploadFileAPIEndpoint, *groupName)
	err := restClient.UploadFile(url, *filePath, &successBody)

	if err != nil {
		fmt.Printf("Problem with the file upload request. %s", err.Error())
		return
	}

	fmt.Printf("File was successfully uploaded in group %s.\n The id of the file is %d\n", *groupName, successBody.FileID)
}

func downloadFile(hostURL, token string) {
	downloadFileCommand := flag.NewFlagSet("download-file", flag.ExitOnError)
	fileID := downloadFileCommand.Int("fileid", -1, "File id")
	groupName := downloadFileCommand.String("grp", "", "Name of the group, owning the file")
	targetPath := downloadFileCommand.String("target", "", "Target destination of file")

	downloadFileCommand.Parse(os.Args[2:])

	if *fileID == -1 || *groupName == "" || *targetPath == "" {
		downloadFileCommand.PrintDefaults()
		return
	}

	restClient := restclient.NewRestClientImpl(token)
	url := fmt.Sprintf("%s%s?group_name=%s&file_id=%d", hostURL, endpoints.DownloadFileAPIEndpoint, *groupName, *fileID)
	err := restClient.DownloadFile(url, *targetPath)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("File was successfully download.\n")
}

type FileRequest struct {
	GroupPayload
	FileID uint `json:"file_id"`
}

func deleteFile(hostURL, token string) {
	deleteFileCommand := flag.NewFlagSet("delete-file", flag.ExitOnError)
	fileID := deleteFileCommand.Int("fileid", -1, "File id")
	groupName := deleteFileCommand.String("grp", "", "Name of the group, in which the file will be uploaded")

	deleteFileCommand.Parse(os.Args[2:])

	if *fileID == -1 || *groupName == "" {
		deleteFileCommand.PrintDefaults()
		return
	}

	reqBody := FileRequest{
		FileID: uint(*fileID),
	}
	reqBody.GroupName = *groupName

	restClient := restclient.NewRestClientImpl(token)
	url := hostURL + endpoints.DeleteFileAPIEndpoint
	err := restClient.Delete(url, &reqBody, nil)

	if err != nil {
		fmt.Printf("Problem with the file deletion request. %s", err.Error())
		return
	}

	fmt.Println("File was successfully deleted")
}

type FileInfo struct {
	ID         uint
	Name       string
	OwnerID    uint
	UploadedAt time.Time
}

type FilesInfoResponse struct {
	Status    int        `json:"status"`
	FilesInfo []FileInfo `json:"files"`
}

func showAllFilesInGroup(hostURL, token string) {
	getAllFilesCommand := flag.NewFlagSet("show-all-files", flag.ExitOnError)
	groupName := getAllFilesCommand.String("grp", "", "Name of the group")

	getAllFilesCommand.Parse(os.Args[2:])

	if *groupName == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	rqBody := GroupPayload{
		GroupName: *groupName,
	}
	successBody := FilesInfoResponse{}

	restClient := restclient.NewRestClientImpl(token)
	url := hostURL + endpoints.GetAllFilesAPIEndpoint
	err := restClient.Post(url, &rqBody, &successBody)

	if err != nil {
		fmt.Printf("Problem with the group creation request. %s", err.Error())
		return
	}

	tableRows := make([]table.Row, len(successBody.FilesInfo))
	for _, fileInfo := range successBody.FilesInfo {
		tableRows = append(tableRows, table.Row{fileInfo.ID, fileInfo.Name, fileInfo.UploadedAt, fileInfo.OwnerID})
	}
	printTable(table.Row{"ID", "Name", "UploadedAt", "OwnerID"}, tableRows)
}

type GroupInfo struct {
	ID      uint
	OwnerID uint
	Name    string
}

type GroupsInfoResponse struct {
	Status     uint        `json:"status"`
	GroupsInfo []GroupInfo `json:"groups"`
}

func showAllGroups(hostURL, token string) {
	successBody := GroupsInfoResponse{}

	restClient := restclient.NewRestClientImpl(token)
	url := hostURL + endpoints.GetAllGroupsAPIEndpoint
	err := restClient.Get(url, &successBody)

	if err != nil {
		fmt.Printf("Problem with the group creation request. %s", err.Error())
		return
	}

	tableRows := make([]table.Row, len(successBody.GroupsInfo))
	for _, groupInfo := range successBody.GroupsInfo {
		tableRows = append(tableRows, table.Row{groupInfo.ID, groupInfo.Name, groupInfo.OwnerID})
	}
	printTable(table.Row{"ID", "Name", "OwnerID"}, tableRows)
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type UsersInfo struct {
	Status    uint       `json:"status"`
	UsersInfo []UserInfo `json:"users"`
}

func showAllUsers(hostURL string, token string) {
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
	printTable(table.Row{"ID", "Username"}, tableRows)
}

func showAllMembers(hostURL, token string) {
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
	printTable(table.Row{"ID", "Username"}, tableRows)
}*/

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
