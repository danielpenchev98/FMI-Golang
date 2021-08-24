package main

const (
	//apiVersionPath - version of the api endpoint
	apiVersionPath = "/v1"
	//publicAPIPath - publicly accessible api path
	publicAPIPath = apiVersionPath + "/public"
	//protectedAPIPath - protected api path
	protectedAPIPath       = apiVersionPath + "/protected"
	healthCheckAPIEndpoint = publicAPIPath + "/healtcheck"
	//LoginAPIEndpoint - api endpoint for user login
	loginAPIEndpoint = publicAPIPath + "/user/login"
	//RegisterAPIEndpoint - api endpoint for user registration
	registerAPIEndpoint = publicAPIPath + "/user/registration"
	//CreateGroupAPIEndpoint - api endpoint for group creation
	createGroupAPIEndpoint = protectedAPIPath + "/group/creation"
	//DeleteGroupAPIEndpoint - api endpoint for group deletion
	deleteGroupAPIEndpoint = protectedAPIPath + "/group/deletion"
	//AddMemberAPIEndpoint - api endpoint for adding an user to a group
	addMemberAPIEndpoint = protectedAPIPath + "/group/invitation"
	//RemoveMemberAPIEndpoint - api endpoint for removing an user from a group
	removeMemberAPIEndpoint = protectedAPIPath + "/group/membership/revocation"
	//UploadFileAPIEndpoint - api endpoint for uploading a file for a specific group
	uploadFileAPIEndpoint = protectedAPIPath + "/group/file/upload"
	//DownloadFileAPIEndpoint - api endpoint for downloading a file from a specific group
	downloadFileAPIEndpoint = protectedAPIPath + "/group/file/download"
	//DeleteFileAPIEndpoint - api endpoint for deleting file, given a group
	deleteFileAPIEndpoint = protectedAPIPath + "/group/file/deletion"
	//GetAllFilesAPIEndpoint - api endpoint for fetching all files, uploaded for a specific group
	getAllFilesAPIEndpoint = protectedAPIPath + "/group/:groupname/files"
	//GetAllGroupsAPIEndpoint - api endpoint for fetching all existing groups
	getAllGroupsAPIEndpoint = protectedAPIPath + "/groups"
	//GetAllUsersAPIEndpoint - api endpoint for fetching all users
	getAllUsersAPIEndpoint = protectedAPIPath + "/users"
	//GetAllMembersAPIEndpoint - api endpoint for fetching all members of a group
	getAllMembersAPIEndpoint = protectedAPIPath + "/group/:groupname/users"
)
