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

		uamDao = &UamDAOImpl{dbConn: gdb}
	})

	AfterEach(func() {
		err := mock.ExpectationsWereMet()
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("CreateUser", func() {
		var (
			username string
			password string
		)

		When("request if user exists fails", func() {
			BeforeEach(func() {
				username = "test-user"
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
			})
		})

		When("request if user exists is succesful", func() {
			Context("and user already exists", func() {
				BeforeEach(func() {
					username, password = "test-user", "test-pass"
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
				})
			})

			Context("and user doesnt exist", func() {
				var rows *sqlmock.Rows
				BeforeEach(func() {
					username, password = "test-user", "test-pass"
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
					})
				})
			})
		})

	})

	Context("Delete user", func() {
		var id uint

		BeforeEach(func() {
			id = 1
		})

		When("request if user exists fails", func() {
			BeforeEach(func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "users"`)).
					WithArgs(id).
					WillReturnError(fmt.Errorf("some error"))
			})

			It("propagates error", func() {
				err := uamDao.DeleteUser(id)
				Expect(err).To(HaveOccurred())
				_, ok := err.(*myerr.ServerError)
				Expect(ok).To(Equal(true))
			})
		})

		When("request if user exists is succesful", func() {
			Context("and user does not exist", func() {
				BeforeEach(func() {
					rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
					mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "users"`)).
						WithArgs(id).
						WillReturnRows(rows)

				})

				It("propagates error", func() {
					err := uamDao.DeleteUser(id)
					Expect(err).To(HaveOccurred())
					_, ok := err.(*myerr.ItemNotFoundError)
					Expect(ok).To(Equal(true))
				})
			})

			Context("and user exists", func() {
				var existCountRows *sqlmock.Rows
				BeforeEach(func() {
					existCountRows = sqlmock.NewRows([]string{"count"}).AddRow(1)
				})
				Context("and deletion query is successful", func() {
					BeforeEach(func() {
						mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "users"`)).
							WithArgs(id).
							WillReturnRows(existCountRows)
						mock.ExpectExec("DELETE FROM \"users\"").
							WithArgs(id).
							WillReturnResult(sqlmock.NewResult(0, 1))
					})

					It("succeeds", func() {
						err := uamDao.DeleteUser(id)
						Expect(err).NotTo(HaveOccurred())
					})
				})

				Context("and deletion query fails", func() {
					BeforeEach(func() {
						mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "users"`)).
							WithArgs(id).
							WillReturnRows(existCountRows)
						mock.ExpectExec("DELETE FROM \"users\"").
							WithArgs(id). // driver.NamedValue - {Name: Ordinal:1 Value:2020-12-28 01:22:59.344298 +0200 EET}"
							WillReturnError(fmt.Errorf("some error"))
					})

					It("propagates error", func() {
						err := uamDao.DeleteUser(id)
						Expect(err).To(HaveOccurred())
						_, ok := err.(*myerr.ServerError)
						Expect(ok).To(Equal(true))
					})
				})
			})
		})
	})

	Context("GetUser", func() {
		const (
			username = "username"
			password = "password"
		)

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
				})
			})
		})
	})

	Context("CreateGroup", func() {
		const (
			groupName = "test-group"
			id        = 1
		)

		When("request if group exists fails", func() {
			BeforeEach(func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "groups"`)).
					WithArgs(groupName).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			})

			It("propagates error", func() {
				err := uamDao.CreateGroup(uint(id), groupName)
				Expect(err).To(HaveOccurred())
				_, ok := err.(*myerr.ServerError)
				Expect(ok).To(Equal(true))
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
					err := uamDao.CreateGroup(uint(id), groupName)
					Expect(err).To(HaveOccurred())
					_, ok := err.(*myerr.ClientError)
					Expect(ok).To(Equal(true))
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
							WithArgs(Any{}, Any{}, groupName, uint(id)). // driver.NamedValue - {Name: Ordinal:1 Value:2020-12-28 01:22:59.344298 +0200 EET}"
							WillReturnError(fmt.Errorf("some error"))
						mock.ExpectRollback()
					})

					It("propagates error", func() {
						err := uamDao.CreateGroup(uint(id), groupName)
						Expect(err).To(HaveOccurred())
						_, ok := err.(*myerr.ServerError)
						Expect(ok).To(Equal(true))
					})
				})

				Context("and group creation query is successful", func() {
					var creationRows *sqlmock.Rows
					var group models.Group

					BeforeEach(func() {
						creationRows = sqlmock.NewRows([]string{"id"}).AddRow(1)
						group = models.Group{
							Name:    groupName,
							OwnerID: uint(id),
						}
						group.ID = uint(id)
					})

					Context("and membership creation query fails", func() {
						BeforeEach(func() {
							mock.ExpectBegin()
							mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "groups"`)).
								WithArgs(groupName).
								WillReturnRows(zeroCountRows)
							mock.ExpectQuery("INSERT INTO \"groups\"").
								WithArgs(Any{}, Any{}, groupName, uint(id)). // driver.NamedValue - {Name: Ordinal:1 Value:2020-12-28 01:22:59.344298 +0200 EET}"
								WillReturnRows(creationRows)
							mock.ExpectQuery("INSERT INTO \"memberships\"").
								WithArgs(Any{}, Any{}, group.ID, group.OwnerID). // driver.NamedValue - {Name: Ordinal:1 Value:2020-12-28 01:22:59.344298 +0200 EET}"
								WillReturnError(fmt.Errorf("some error"))
							mock.ExpectRollback()
						})

						It("succeeds", func() {
							err := uamDao.CreateGroup(uint(id), groupName)
							Expect(err).To(HaveOccurred())
							_, ok := err.(*myerr.ServerError)
							Expect(ok).To(Equal(true))
						})
					})

					Context("and membership creation query is successful", func() {
						BeforeEach(func() {
							mock.ExpectBegin()
							mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "groups"`)).
								WithArgs(groupName).
								WillReturnRows(zeroCountRows)
							mock.ExpectQuery("INSERT INTO \"groups\"").
								WithArgs(Any{}, Any{}, groupName, uint(id)). // driver.NamedValue - {Name: Ordinal:1 Value:2020-12-28 01:22:59.344298 +0200 EET}"
								WillReturnRows(creationRows)
							mock.ExpectQuery("INSERT INTO \"memberships\"").
								WithArgs(Any{}, Any{}, group.ID, group.OwnerID). // driver.NamedValue - {Name: Ordinal:1 Value:2020-12-28 01:22:59.344298 +0200 EET}"
								WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
							mock.ExpectCommit()
						})

						It("succeeds", func() {
							err := uamDao.CreateGroup(uint(id), groupName)
							Expect(err).NotTo(HaveOccurred())
						})
					})
				})
			})
		})
	})
})
