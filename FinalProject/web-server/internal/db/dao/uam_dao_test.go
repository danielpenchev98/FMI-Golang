package dao

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"regexp"
	"time"

	"example.com/user/web-server/internal/db/models"
	myerr "example.com/user/web-server/internal/error"
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Any struct{}

// Match satisfies sqlmock.Argument interface
func (a Any) Match(v driver.Value) bool {
	return true
}

var _ = Describe("UamDAO", func() {
	var (
		uamDao UamDAO
		mock   sqlmock.Sqlmock
	)

	const (
		username  = "username"
		password  = "password"
		groupName = "test-group"
		userID    = 1
		groupID   = 2
	)

	BeforeEach(func() {
		var (
			db  *sql.DB
			err error
		)

		db, mock, err = sqlmock.New()
		Expect(err).NotTo(HaveOccurred())

		gdb, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})
		Expect(err).NotTo(HaveOccurred())

		uamDao = NewUamDAOImpl(gdb)
	})

	AfterEach(func() {
		err := mock.ExpectationsWereMet()
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("CreateUser", func() {
		When("request if user exists fails", func() {
			BeforeEach(func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "users"`)).
					WithArgs(username).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()

			})

			It("propagates error", func() {
				err := uamDao.CreateUser(username, password)
				Expect(err).To(HaveOccurred())
				_, ok := err.(*myerr.ServerError)
				Expect(ok).To(Equal(true))
				Expect(mock.ExpectationsWereMet()).To(BeNil())
			})
		})

		When("request if user exists is succesful", func() {
			Context("and user already exists", func() {
				BeforeEach(func() {
					rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

					mock.ExpectBegin()
					mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "users"`)).
						WithArgs(username).
						WillReturnRows(rows)
					mock.ExpectRollback()
				})

				It("propagates error", func() {
					err := uamDao.CreateUser(username, password)
					Expect(err).To(HaveOccurred())
					_, ok := err.(*myerr.ClientError)
					Expect(ok).To(Equal(true))
					Expect(mock.ExpectationsWereMet()).To(BeNil())
				})
			})

			Context("and user doesnt exist", func() {
				var rows *sqlmock.Rows
				BeforeEach(func() {
					rows = sqlmock.NewRows([]string{"count"}).AddRow(0)
				})

				Context("and creation query is successful", func() {
					BeforeEach(func() {
						mock.ExpectBegin()
						mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "users"`)).
							WithArgs(username).
							WillReturnRows(rows)
						mock.ExpectQuery("INSERT INTO \"users\"").
							WithArgs(Any{}, Any{}, username, password). // driver.NamedValue - {Name: Ordinal:1 Value:2020-12-28 01:22:59.344298 +0200 EET}"
							WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
						mock.ExpectCommit()
					})

					It("succeeds", func() {
						err := uamDao.CreateUser(username, password)
						Expect(err).NotTo(HaveOccurred())
						Expect(mock.ExpectationsWereMet()).To(BeNil())
					})
				})

				Context("and creation query fails", func() {
					BeforeEach(func() {
						mock.ExpectBegin()
						mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "users"`)).
							WithArgs(username).
							WillReturnRows(rows)
						mock.ExpectQuery("INSERT INTO \"users\"").
							WithArgs(Any{}, Any{}, username, password). // driver.NamedValue - {Name: Ordinal:1 Value:2020-12-28 01:22:59.344298 +0200 EET}"
							WillReturnError(fmt.Errorf("some error"))
						mock.ExpectRollback()
					})

					It("propagates error", func() {
						err := uamDao.CreateUser(username, password)
						Expect(err).To(HaveOccurred())
						_, ok := err.(*myerr.ServerError)
						Expect(ok).To(Equal(true))
						Expect(mock.ExpectationsWereMet()).To(BeNil())
					})
				})
			})
		})

	})

	Context("Delete user", func() {
		When("request if user exists fails", func() {
			BeforeEach(func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "users"`)).
					WithArgs(userID).
					WillReturnError(fmt.Errorf("some error"))
			})

			It("propagates error", func() {
				err := uamDao.DeleteUser(uint(userID))
				Expect(err).To(HaveOccurred())
				_, ok := err.(*myerr.ServerError)
				Expect(ok).To(Equal(true))
				Expect(mock.ExpectationsWereMet()).To(BeNil())
			})
		})

		When("request if user exists is succesful", func() {
			Context("and user does not exist", func() {
				BeforeEach(func() {
					rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
					mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "users"`)).
						WithArgs(userID).
						WillReturnRows(rows)

				})

				It("propagates error", func() {
					err := uamDao.DeleteUser(uint(userID))
					Expect(err).To(HaveOccurred())
					_, ok := err.(*myerr.ItemNotFoundError)
					Expect(ok).To(Equal(true))
					Expect(mock.ExpectationsWereMet()).To(BeNil())
				})
			})

			Context("and user exists", func() {
				var existCountRows *sqlmock.Rows
				BeforeEach(func() {
					existCountRows = sqlmock.NewRows([]string{"count"}).AddRow(1)
				})

				Context("and deletion query fails", func() {
					BeforeEach(func() {
						mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "users"`)).
							WithArgs(userID).
							WillReturnRows(existCountRows)
						mock.ExpectExec("DELETE FROM \"users\"").
							WithArgs(userID). // driver.NamedValue - {Name: Ordinal:1 Value:2020-12-28 01:22:59.344298 +0200 EET}"
							WillReturnError(fmt.Errorf("some error"))
					})

					It("propagates error", func() {
						err := uamDao.DeleteUser(userID)
						Expect(err).To(HaveOccurred())
						_, ok := err.(*myerr.ServerError)
						Expect(ok).To(Equal(true))
						Expect(mock.ExpectationsWereMet()).To(BeNil())
					})
				})

				Context("and deletion query is successful", func() {
					BeforeEach(func() {
						mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "users"`)).
							WithArgs(userID).
							WillReturnRows(existCountRows)
						mock.ExpectExec("DELETE FROM \"users\"").
							WithArgs(userID).
							WillReturnResult(sqlmock.NewResult(0, 1))
					})

					It("succeeds", func() {
						err := uamDao.DeleteUser(uint(userID))
						Expect(err).NotTo(HaveOccurred())
						Expect(mock.ExpectationsWereMet()).To(BeNil())
					})
				})
			})
		})
	})

	Context("GetUser", func() {
		When("request to get userID fails", func() {
			BeforeEach(func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
					WithArgs(username).
					WillReturnError(fmt.Errorf("some error"))
			})

			It("propagates error", func() {
				_, err := uamDao.GetUser(username)
				Expect(err).To(HaveOccurred())
				_, ok := err.(*myerr.ServerError)
				Expect(ok).To(Equal(true))
				Expect(mock.ExpectationsWereMet()).To(BeNil())
			})
		})

		When("request to get userID is succesful", func() {
			Context("and user does not exist", func() {
				BeforeEach(func() {
					mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
						WithArgs(username).
						WillReturnError(gorm.ErrRecordNotFound)
				})

				It("propagates error", func() {
					_, err := uamDao.GetUser(username)
					Expect(err).To(HaveOccurred())
					_, ok := err.(*myerr.ItemNotFoundError)
					Expect(ok).To(Equal(true))
					Expect(mock.ExpectationsWereMet()).To(BeNil())
				})
			})

			Context("and user exists", func() {
				var mockTime time.Time

				BeforeEach(func() {
					mockTime = time.Now()
					rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username", "password"}).AddRow(1, mockTime, mockTime, username, password)
					mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
						WithArgs(username).
						WillReturnRows(rows)
				})

				It("succeds", func() {
					user, err := uamDao.GetUser(username)
					Expect(err).NotTo(HaveOccurred())

					Expect(user.ID).To(Equal(uint(1)))
					Expect(user.Username).To(Equal(username))
					Expect(user.Password).To(Equal(password))
					Expect(user.CreatedAt).To(Equal(mockTime))
					Expect(user.UpdatedAt).To(Equal(mockTime))
					Expect(mock.ExpectationsWereMet()).To(BeNil())
				})
			})
		})
	})

	Context("CreateGroup", func() {
		When("request if group exists fails", func() {
			BeforeEach(func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "groups"`)).
					WithArgs(groupName).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			})

			It("propagates error", func() {
				err := uamDao.CreateGroup(uint(userID), groupName)
				Expect(err).To(HaveOccurred())
				_, ok := err.(*myerr.ServerError)
				Expect(ok).To(Equal(true))
				Expect(mock.ExpectationsWereMet()).To(BeNil())
			})
		})

		When("request if group exists is succesful", func() {
			Context("and user already exists", func() {
				BeforeEach(func() {
					rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

					mock.ExpectBegin()
					mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "groups"`)).
						WithArgs(groupName).
						WillReturnRows(rows)
					mock.ExpectRollback()
				})

				It("propagates error", func() {
					err := uamDao.CreateGroup(uint(userID), groupName)
					Expect(err).To(HaveOccurred())
					_, ok := err.(*myerr.ClientError)
					Expect(ok).To(Equal(true))
					Expect(mock.ExpectationsWereMet()).To(BeNil())
				})
			})

			Context("and group doesnt exist", func() {
				var zeroCountRows *sqlmock.Rows

				BeforeEach(func() {
					zeroCountRows = sqlmock.NewRows([]string{"count"}).AddRow(0)
				})

				Context("and group creation query fails", func() {
					BeforeEach(func() {
						mock.ExpectBegin()
						mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "groups"`)).
							WithArgs(groupName).
							WillReturnRows(zeroCountRows)
						mock.ExpectQuery("INSERT INTO \"groups\"").
							WithArgs(Any{}, Any{}, groupName, userID). // driver.NamedValue - {Name: Ordinal:1 Value:2020-12-28 01:22:59.344298 +0200 EET}"
							WillReturnError(fmt.Errorf("some error"))
						mock.ExpectRollback()
					})

					It("propagates error", func() {
						err := uamDao.CreateGroup(uint(userID), groupName)
						Expect(err).To(HaveOccurred())
						_, ok := err.(*myerr.ServerError)
						Expect(ok).To(Equal(true))
						Expect(mock.ExpectationsWereMet()).To(BeNil())
					})
				})

				Context("and group creation query is successful", func() {
					var creationRows *sqlmock.Rows
					var group models.Group

					BeforeEach(func() {
						group = models.Group{
							Name:    groupName,
							OwnerID: uint(userID),
						}
						group.ID = uint(userID)
						creationRows = sqlmock.NewRows([]string{"id"}).AddRow(group.ID)
					})

					Context("and membership creation query fails", func() {
						BeforeEach(func() {
							mock.ExpectBegin()
							mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "groups"`)).
								WithArgs(groupName).
								WillReturnRows(zeroCountRows)
							mock.ExpectQuery("INSERT INTO \"groups\"").
								WithArgs(Any{}, Any{}, groupName, userID). // driver.NamedValue - {Name: Ordinal:1 Value:2020-12-28 01:22:59.344298 +0200 EET}"
								WillReturnRows(creationRows)
							mock.ExpectQuery("INSERT INTO \"memberships\"").
								WithArgs(Any{}, Any{}, group.ID, group.OwnerID). // driver.NamedValue - {Name: Ordinal:1 Value:2020-12-28 01:22:59.344298 +0200 EET}"
								WillReturnError(fmt.Errorf("some error"))
							mock.ExpectRollback()
						})

						It("succeeds", func() {
							err := uamDao.CreateGroup(uint(userID), groupName)
							Expect(err).To(HaveOccurred())
							_, ok := err.(*myerr.ServerError)
							Expect(ok).To(Equal(true))
							Expect(mock.ExpectationsWereMet()).To(BeNil())
						})
					})

					Context("and membership creation query is successful", func() {
						BeforeEach(func() {
							mock.ExpectBegin()
							mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "groups"`)).
								WithArgs(groupName).
								WillReturnRows(zeroCountRows)
							mock.ExpectQuery("INSERT INTO \"groups\"").
								WithArgs(Any{}, Any{}, groupName, userID). // driver.NamedValue - {Name: Ordinal:1 Value:2020-12-28 01:22:59.344298 +0200 EET}"
								WillReturnRows(creationRows)
							mock.ExpectQuery("INSERT INTO \"memberships\"").
								WithArgs(Any{}, Any{}, group.ID, group.OwnerID). // driver.NamedValue - {Name: Ordinal:1 Value:2020-12-28 01:22:59.344298 +0200 EET}"
								WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
							mock.ExpectCommit()
						})

						It("succeeds", func() {
							err := uamDao.CreateGroup(uint(userID), groupName)
							Expect(err).NotTo(HaveOccurred())
							Expect(mock.ExpectationsWereMet()).To(BeNil())
						})
					})
				})
			})
		})
	})

	Context("AddUserToGroup", func() {
		When("get group request fails", func() {
			Context("problem with the database", func() {
				BeforeEach(func() {
					mock.ExpectBegin()
					mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups"`)).
						WithArgs(groupName).
						WillReturnError(fmt.Errorf("some error"))
					mock.ExpectRollback()
				})

				It("propagates error", func() {
					err := uamDao.AddUserToGroup(uint(userID), username, groupName)
					Expect(err).To(HaveOccurred())
					_, ok := err.(*myerr.ServerError)
					Expect(ok).To(Equal(true))
					Expect(mock.ExpectationsWereMet()).To(BeNil())
				})
			})

			Context("and group doesnt exist", func() {
				BeforeEach(func() {
					mock.ExpectBegin()
					mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups"`)).
						WithArgs(groupName).
						WillReturnError(gorm.ErrRecordNotFound)
					mock.ExpectRollback()
				})

				It("propagates error", func() {
					err := uamDao.AddUserToGroup(uint(userID), username, groupName)
					Expect(err).To(HaveOccurred())
					_, ok := err.(*myerr.ItemNotFoundError)
					Expect(ok).To(Equal(true))
					Expect(mock.ExpectationsWereMet()).To(BeNil())
				})
			})
		})

		When("get group request is sucessfull", func() {
			Context("and you arent the owner of the group", func() {
				BeforeEach(func() {
					rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "owner_id"}).
						AddRow(1, time.Now(), time.Now(), groupName, userID+1)
					mock.ExpectBegin()
					mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups"`)).
						WithArgs(groupName).
						WillReturnRows(rows)
					mock.ExpectRollback()
				})

				It("propagates error", func() {
					err := uamDao.AddUserToGroup(uint(userID), username, groupName)
					Expect(err).To(HaveOccurred())
					_, ok := err.(*myerr.ClientError)
					Expect(ok).To(Equal(true))
					Expect(mock.ExpectationsWereMet()).To(BeNil())
				})
			})

			Context("and you are the owner of the group", func() {
				var groupRow *sqlmock.Rows
				BeforeEach(func() {
					groupRow = sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "owner_id"}).
						AddRow(groupID, time.Now(), time.Now(), groupName, userID)
				})

				Context("and user get request fails", func() {
					Context("problem with the database", func() {
						BeforeEach(func() {
							mock.ExpectBegin()
							mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups"`)).
								WithArgs(groupName).
								WillReturnRows(groupRow)
							mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
								WithArgs(username).
								WillReturnError(fmt.Errorf("some error"))
							mock.ExpectRollback()
						})

						It("propagates error", func() {
							err := uamDao.AddUserToGroup(uint(userID), username, groupName)
							Expect(err).To(HaveOccurred())
							_, ok := err.(*myerr.ServerError)
							Expect(ok).To(Equal(true))
							Expect(mock.ExpectationsWereMet()).To(BeNil())
						})
					})

					Context("user does not exist", func() {
						BeforeEach(func() {
							mock.ExpectBegin()
							mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups"`)).
								WithArgs(groupName).
								WillReturnRows(groupRow)
							mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
								WithArgs(username).
								WillReturnError(gorm.ErrRecordNotFound)
							mock.ExpectRollback()
						})

						It("propagates error", func() {
							err := uamDao.AddUserToGroup(uint(userID), username, groupName)
							Expect(err).To(HaveOccurred())
							_, ok := err.(*myerr.ItemNotFoundError)
							Expect(ok).To(Equal(true))
							Expect(mock.ExpectationsWereMet()).To(BeNil())
						})
					})

				})

				Context("and user request is successfull", func() {
					var userRows *sqlmock.Rows
					BeforeEach(func() {
						userRows = sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username", "password"}).
							AddRow(userID, time.Now(), time.Now(), username, password)
					})

					Context("and request if membership exsists fails", func() {
						Context("and problem with lookup in the db for the membership", func() {
							BeforeEach(func() {
								mock.ExpectBegin()
								mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups"`)).
									WithArgs(groupName).
									WillReturnRows(groupRow)
								mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
									WithArgs(username).
									WillReturnRows(userRows)
								mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "memberships"`)).
									WithArgs(groupID, userID).
									WillReturnError(fmt.Errorf("some error"))
								mock.ExpectRollback()
							})

							It("propagates error", func() {
								err := uamDao.AddUserToGroup(uint(userID), username, groupName)
								Expect(err).To(HaveOccurred())
								_, ok := err.(*myerr.ServerError)
								Expect(ok).To(Equal(true))
								Expect(mock.ExpectationsWereMet()).To(BeNil())
							})
						})
						Context("and memberships found", func() {
							BeforeEach(func() {
								mock.ExpectBegin()
								mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups"`)).
									WithArgs(groupName).
									WillReturnRows(groupRow)
								mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
									WithArgs(username).
									WillReturnRows(userRows)
								mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "memberships"`)).
									WithArgs(groupID, userID).
									WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
								mock.ExpectRollback()
							})

							It("propagates error", func() {
								err := uamDao.AddUserToGroup(uint(userID), username, groupName)
								Expect(err).To(HaveOccurred())
								_, ok := err.(*myerr.ClientError)
								Expect(ok).To(Equal(true))
								Expect(mock.ExpectationsWereMet()).To(BeNil())
							})
						})
					})

					Context("and request if membership exists is successful", func() {
						Context("and creation of membership fails", func() {
							BeforeEach(func() {
								mock.ExpectBegin()
								mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups"`)).
									WithArgs(groupName).
									WillReturnRows(groupRow)
								mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
									WithArgs(username).
									WillReturnRows(userRows)
								mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "memberships"`)).
									WithArgs(groupID, userID).
									WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
								mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "memberships"`)).
									WithArgs(Any{}, Any{}, groupID, userID).
									WillReturnError(fmt.Errorf("some error"))
								mock.ExpectRollback()
							})

							It("propagates error", func() {
								err := uamDao.AddUserToGroup(uint(userID), username, groupName)
								Expect(err).To(HaveOccurred())
								_, ok := err.(*myerr.ServerError)
								Expect(ok).To(Equal(true))
								Expect(mock.ExpectationsWereMet()).To(BeNil())
							})
						})
						Context("and creation succeeds", func() {
							BeforeEach(func() {
								mock.ExpectBegin()
								mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups"`)).
									WithArgs(groupName).
									WillReturnRows(groupRow)
								mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
									WithArgs(username).
									WillReturnRows(userRows)
								mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "memberships"`)).
									WithArgs(groupID, userID).
									WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
								mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "memberships"`)).
									WithArgs(Any{}, Any{}, groupID, userID).
									WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
								mock.ExpectCommit()
							})

							It("returns no error", func() {
								err := uamDao.AddUserToGroup(uint(userID), username, groupName)
								Expect(err).NotTo(HaveOccurred())
								Expect(mock.ExpectationsWereMet()).To(BeNil())
							})
						})
					})
				})
			})
		})

	})

	Context("RemoveUserFromGroup", func() {

		When("get group request fails", func() {
			Context("problem with the database", func() {
				BeforeEach(func() {
					mock.ExpectBegin()
					mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups"`)).
						WithArgs(groupName).
						WillReturnError(fmt.Errorf("some error"))
					mock.ExpectRollback()
				})

				It("propagates error", func() {
					err := uamDao.RemoveUserFromGroup(uint(userID), username, groupName)
					Expect(err).To(HaveOccurred())
					_, ok := err.(*myerr.ServerError)
					Expect(ok).To(Equal(true))
					Expect(mock.ExpectationsWereMet()).To(BeNil())
				})
			})

			Context("and group doesnt exist", func() {
				BeforeEach(func() {
					mock.ExpectBegin()
					mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups"`)).
						WithArgs(groupName).
						WillReturnError(gorm.ErrRecordNotFound)
					mock.ExpectRollback()
				})

				It("propagates error", func() {
					err := uamDao.RemoveUserFromGroup(uint(userID), username, groupName)
					Expect(err).To(HaveOccurred())
					_, ok := err.(*myerr.ItemNotFoundError)
					Expect(ok).To(Equal(true))
					Expect(mock.ExpectationsWereMet()).To(BeNil())
				})
			})
		})

		When("get group request succeeds", func() {
			var groupRow *sqlmock.Rows
			BeforeEach(func() {
				groupRow = sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "owner_id"}).
					AddRow(groupID, time.Now(), time.Now(), groupName, userID)
			})

			Context("and get user request fails", func() {
				Context("problem with the database", func() {
					BeforeEach(func() {
						mock.ExpectBegin()
						mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups"`)).
							WithArgs(groupName).
							WillReturnRows(groupRow)
						mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
							WithArgs(username).
							WillReturnError(fmt.Errorf("some error"))
						mock.ExpectRollback()
					})

					It("propagates error", func() {
						err := uamDao.RemoveUserFromGroup(uint(userID), username, groupName)
						Expect(err).To(HaveOccurred())
						_, ok := err.(*myerr.ServerError)
						Expect(ok).To(Equal(true))
						Expect(mock.ExpectationsWereMet()).To(BeNil())
					})
				})

				Context("user does not exist", func() {
					BeforeEach(func() {
						mock.ExpectBegin()
						mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups"`)).
							WithArgs(groupName).
							WillReturnRows(groupRow)
						mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
							WithArgs(username).
							WillReturnError(gorm.ErrRecordNotFound)
						mock.ExpectRollback()
					})

					It("propagates error", func() {
						err := uamDao.RemoveUserFromGroup(uint(userID), username, groupName)
						Expect(err).To(HaveOccurred())
						_, ok := err.(*myerr.ItemNotFoundError)
						Expect(ok).To(Equal(true))
						Expect(mock.ExpectationsWereMet()).To(BeNil())
					})
				})

			})

			Context("and get user request succeeds", func() {
				var userRows *sqlmock.Rows
				BeforeEach(func() {
					userRows = sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username", "password"}).
						AddRow(userID, time.Now(), time.Now(), username, password)
				})

				Context("and you arent the owner of the group or the targeted user", func() {
					BeforeEach(func() {
						mock.ExpectBegin()
						mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups"`)).
							WithArgs(groupName).
							WillReturnRows(groupRow)
						mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
							WithArgs(username).
							WillReturnRows(userRows)
						mock.ExpectRollback()
					})

					It("propagates error", func() {
						err := uamDao.RemoveUserFromGroup(uint(userID+1), username, groupName)
						Expect(err).To(HaveOccurred())
						_, ok := err.(*myerr.ClientError)
						Expect(ok).To(Equal(true))
						Expect(mock.ExpectationsWereMet()).To(BeNil())
					})
				})

				Context("and you are either the owner or the targeted users", func() {
					Context("and delete membership request fails", func() {
						BeforeEach(func() {
							mock.ExpectBegin()
							mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups"`)).
								WithArgs(groupName).
								WillReturnRows(groupRow)
							mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
								WithArgs(username).
								WillReturnRows(userRows)
							mock.ExpectExec("DELETE FROM \"memberships\"").
								WithArgs(userID, groupID).
								WillReturnError(fmt.Errorf("some error"))
							mock.ExpectRollback()
						})

						It("propagates error", func() {
							err := uamDao.RemoveUserFromGroup(uint(userID), username, groupName)
							Expect(err).To(HaveOccurred())
							_, ok := err.(*myerr.ServerError)
							Expect(ok).To(Equal(true))
							Expect(mock.ExpectationsWereMet()).To(BeNil())
						})
					})
					Context("and delete membership requests succeeds", func() {
						Context("and membership did not exist", func() {
							BeforeEach(func() {
								mock.ExpectBegin()
								mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups"`)).
									WithArgs(groupName).
									WillReturnRows(groupRow)
								mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
									WithArgs(username).
									WillReturnRows(userRows)
								mock.ExpectExec("DELETE FROM \"memberships\"").
									WithArgs(userID, groupID).
									WillReturnResult(sqlmock.NewResult(0, 0))
								mock.ExpectRollback()
							})

							It("propagates error", func() {
								err := uamDao.RemoveUserFromGroup(uint(userID), username, groupName)
								Expect(err).To(HaveOccurred())
								_, ok := err.(*myerr.ClientError)
								Expect(ok).To(Equal(true))
								Expect(mock.ExpectationsWereMet()).To(BeNil())
							})
						})

						Context("and existing membership successfully deleted", func() {
							BeforeEach(func() {
								mock.ExpectBegin()
								mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups"`)).
									WithArgs(groupName).
									WillReturnRows(groupRow)
								mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
									WithArgs(username).
									WillReturnRows(userRows)
								mock.ExpectExec("DELETE FROM \"memberships\"").
									WithArgs(userID, groupID).
									WillReturnResult(sqlmock.NewResult(0, 1))
								mock.ExpectCommit()
							})

							It("propagates error", func() {
								err := uamDao.RemoveUserFromGroup(uint(userID), username, groupName)
								Expect(err).NotTo(HaveOccurred())
								Expect(mock.ExpectationsWereMet()).To(BeNil())
							})
						})
					})
				})
			})
		})

	})
})
