package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-client/internal/restclient"
)

type CredentialsPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type BasicResponse struct {
	Status int `json:"status"`
}

func registerUser() {
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
	err := restClient.Post("http://localhost:8080/v1/public/user/registration", &rqBody, &successBody)

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

func login() {
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
	err := restClient.Post("http://localhost:8080/v1/public/user/login", &rqBody, &successBody)

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
}

type GroupPayload struct {
	GroupName string `json:"group_name"`
}

func createGroup(token string) {
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
	err := restClient.Post("http://localhost:8080/v1/protected/group/creation", &rqBody, nil)

	if err != nil {
		fmt.Printf("Problem with the group creation request. %s", err.Error())
		return
	}

	fmt.Printf("Group %s was succesfully created", *groupName)
}

func deleteGroup(token string) {
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
	err := restClient.Delete("http://localhost:8080/v1/protected/group/deletion", &rqBody, nil)

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

func addUserToGroup(token string) {
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
	err := restClient.Post("http://localhost:8080/v1/protected/group/invitation", &rqBody, nil)

	if err != nil {
		fmt.Printf("Problem with the group creation request. %s", err.Error())
		return
	}

	fmt.Printf("User %s was successfully added to group %s\n", *username, *groupName)
}

func removeUserFromGroup(token string) {
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
	err := restClient.Delete("http://localhost:8080/v1/protected/group/membership/revocation", &rqBody, nil)

	if err != nil {
		fmt.Printf("Problem with the group creation request. %s", err.Error())
		return
	}

	fmt.Printf("User %s was successfully removed from group %s\n", *username, *groupName)
}

type FileUploadResponse struct {
	FileID uint `json:"file_id"`
}

func uploadFile(token string) {
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
	url := fmt.Sprintf("http://localhost:8080/v1/protected/group/file/upload?group_name=%s", *groupName)
	err := restClient.UploadFile(url, *filePath, &successBody)

	if err != nil {
		fmt.Printf("Problem with the file upload request. %s", err.Error())
		return
	}

	fmt.Printf("File was successfully uploaded in group %s.\n The id of the file is %d\n", *groupName, successBody.FileID)
}

func downloadFile(token string) {
	downloadFileCommand := flag.NewFlagSet("download-file", flag.ExitOnError)
	fileID := downloadFileCommand.Int("fileid", -1, "File id")
	groupName := downloadFileCommand.String("grp", "", "Name of the group, owning the file")
	targetPath := downloadFileCommand.String("target", "", "Target destination of file")

	downloadFileCommand.Parse(os.Args[2:])

	if *fileID == -1 || *groupName == "" || *targetPath == "" {
		downloadFileCommand.PrintDefaults()
		return
	}

	req := FileRequest{
		FileID: uint(*fileID),
	}
	req.GroupName = *groupName

	restClient := restclient.NewRestClientImpl(token)
	err := restClient.DownloadFile("http://localhost:8080/v1/protected/group/file/download", *targetPath, &req)

	if err != nil {
		fmt.Printf("Problem with the file download request. %s", err.Error())
		return
	}

	fmt.Println("File was successfully download.\n")
}

type FileRequest struct {
	GroupPayload
	FileID uint `json:"file_id"`
}

func deleteFile(token string) {
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
	err := restClient.Delete("http://localhost:8080/v1/protected/group/file/delete", &reqBody, nil)

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

func showAllFilesInGroup(token string) {
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
	err := restClient.Post("http://localhost:8080/v1/protected/group/membership/reinvocation", &rqBody, &successBody)

	if err != nil {
		fmt.Printf("Problem with the group creation request. %s", err.Error())
		return
	}

	fmt.Println("---------Files----------")
	for _, fileInfo := range successBody.FilesInfo {
		fmt.Printf("ID: %d\nName: %s\nUploadedAt: %v\nOwnerID: %d\n", fileInfo.ID, fileInfo.Name, fileInfo.UploadedAt, fileInfo.OwnerID)
		fmt.Println("---------------")
	}

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

func showAllGroups(token string) {
	successBody := GroupsInfoResponse{}

	restClient := restclient.NewRestClientImpl(token)
	err := restClient.Get("http://localhost:8080/v1/protected/groups", &successBody)

	if err != nil {
		fmt.Printf("Problem with the group creation request. %s", err.Error())
		return
	}

	fmt.Println("---------Groups-----------")
	for _, groupInfo := range successBody.GroupsInfo {
		fmt.Printf("ID: %d\nName: %s\nOwnerID: %d\n", groupInfo.ID, groupInfo.Name, groupInfo.OwnerID)
		fmt.Println("------------")
	}
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type UsersInfo struct {
	Status    uint       `json:"status"`
	UsersInfo []UserInfo `json:"users"`
}

func showAllUsers(token string) {
	successBody := UsersInfo{}

	restClient := restclient.NewRestClientImpl(token)
	err := restClient.Get("http://localhost:8080/v1/protected/users", &successBody)

	if err != nil {
		fmt.Printf("Problem with the group creation request. %s", err.Error())
		return
	}

	fmt.Println("---------Users-----------")
	for _, userInfo := range successBody.UsersInfo {
		fmt.Printf("ID: %d\nUsername: %s\n", userInfo.ID, userInfo.Username)
		fmt.Println("------------")
	}
}

func showAllMembers(token string) {
	getAllMembers := flag.NewFlagSet("show-all-members", flag.ExitOnError)
	groupName := getAllMembers.String("grp", "", "Name of the group")

	getAllMembers.Parse(os.Args[2:])

	if *groupName == "" {
		getAllMembers.PrintDefaults()
		os.Exit(1)
	}
	successBody := UsersInfo{}

	restClient := restclient.NewRestClientImpl(token)
	url := fmt.Sprintf("http://localhost:8080/v1/protected/group/users?group_name=%s", *groupName)
	err := restClient.Get(url, &successBody)

	if err != nil {
		fmt.Printf("Problem with the group creation request. %s", err.Error())
		return
	}

	fmt.Println("---------Users-----------")
	for _, userInfo := range successBody.UsersInfo {
		fmt.Printf("ID: %d\nName: %s\n", userInfo.ID, userInfo.Username)
		fmt.Println("------------")
	}
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

func commandsWithAuth(command string) {
	token, err := retrieveJwtToken()
	if err != nil {
		fmt.Println("Couldnt create a group. Reason: %s", err.Error())
		os.Exit(1)
	}

	switch command {
	case "logout":
		logout()
	case "create-group":
		createGroup(token)
	case "delete-group":
		deleteGroup(token)
	case "add-member":
		addUserToGroup(token)
	case "remove-member":
		removeUserFromGroup(token)
	case "upload-file":
		uploadFile(token)
	case "download-file":
		downloadFile(token)
	case "delete-file":
		deleteFile(token)
	case "show-all-files":
		showAllFilesInGroup(token)
	case "show-all-groups":
		showAllGroups(token)
	case "show-all-users":
		showAllUsers(token)
	case "show-all-members":
		showAllMembers(token)
	default:
		fmt.Printf("Invalid command [%s]\n", command)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Print("No arguments were supplied")
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "login":
		login()
	case "register":
		registerUser()
	default:
		commandsWithAuth(command)
	}

}
