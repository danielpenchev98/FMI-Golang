package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-client/internal/endpoints"
	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-client/internal/restclient"
	"github.com/jedib0t/go-pretty/v6/table"
)

type MembershipRequest struct {
	GroupPayload
	Username string `json:"username"`
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

func CreateGroup(hostURL, token string) {
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

func DeleteGroup(hostURL, token string) {
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

func AddMember(hostURL, token string) {
	addMemberCommand := flag.NewFlagSet("add-member", flag.ExitOnError)
	username := addMemberCommand.String("usr", "", "Name of the user to be added to the group")
	groupName := addMemberCommand.String("grp", "", "Name of the group to be deleted")
	addMemberCommand.Parse(os.Args[2:])

	if *groupName == "" || *username == "" {
		addMemberCommand.PrintDefaults()
		return
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

func RemoveMember(hostURL, token string) {
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
}

func ShowAllGroups(hostURL, token string) {
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
	PrintTable(table.Row{"ID", "Name", "OwnerID"}, tableRows)
}
