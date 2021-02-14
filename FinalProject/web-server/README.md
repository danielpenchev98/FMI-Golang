# Web-Server

## Description
The main idea of this server is to allow people to share freely their files with everyone.
In order to keep the anonymity of the users and the interformation they share, the following limitations are imposed:
* Users are grouped in `groups`. One user can belong to many groups and one group can have many members
* File are uploaded, given a specific `group`. Only the members of the `group` can access/view the `group` files
* The only identification of the user is his `username` (also his `id`)
Also there are limitations in terms of implementation:
* Only the `owner` of the `group` and the `owner` of the file can delete it from the group
* When the `owner` deletes the group or deletes his account, there is no transition of ownership (yet). Instead all group recources are deleted (files, memberships, etc)
* The group resources arent deleted immediately. Instead, when the group is request to be deleted, the group swithces to `deactivated` state. And after a particular time period the rosources are erased. After this operation succeeds, the name of the `group` is available for usage.

## Configuration
The server uses the following external dependencies, which should be installed:
### Production
* `go` - preferably versions above `1.5.0`
* `github.com/dgrijalva/jwt-go` - used for validation/creation of JWTokens
* `github.com/gin-gonic/gin` - used for the implementation of the REST API
* `github.com/pkg/errors` - used for easier creation of errors
* `github.com/robfig/cron/v3` - used for the async job for deletion of group resources
* `golang.org/x/crypto` - used for encryption of user information
* `gorm.io/gorm` - used for mapping models (go structs) to sql tables
* `gorm.io/driver/postgres` - used for the communication with the `postgres` database
### Testing
* `github.com/DATA-DOG/go-sqlmock` - used for testing the request, sent to the database
* `github.com/golang/mock` - used for mocking external dependencies
* `github.com/onsi/ginkgo` - used as the main testing framework
* `github.com/onsi/gomega` - used for assertions

The following environment variables must be set:
### DB configuration
* `DB_NAME` - env variable, containing the name of the database
* `DB_USER` - env variable, containing the db username
* `DB_PASS` - env variable, containing the db password
* `DB_PORT` - env variable, containing the port on which the db server is running on
* `DB_HOST` - env variable, containing the domain of the db server
### Auth configuration
* `SECRET` - env variable, containing a value, used for the encryption/decryption of the token
* `ISSUER` - env variable, containing the name of authority, issuing the token
* `EXPIRATION` - env variable, containing the expiration time of the issued tokens (in hours)

## Server startup
In the `cmd` package run the following command:
```Golang
go run main.go
```