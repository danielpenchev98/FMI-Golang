package dao

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

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
				mock.ExpectQuery("SELECT count").
					WithArgs(username).
					WillReturnError(fmt.Errorf("some error"))
			})

			It("propagates error", func() {
				_, err := uamDao.CreateUser(username, password)
				_, ok := err.(*myerr.ServerError)
				Expect(ok).To(Equal(true))
			})
		})

		When("request if user exists is succesful", func() {
			Context("and user already exists", func() {
				BeforeEach(func() {
					username, password = "test-user", "test-pass"
					rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
					mock.ExpectQuery("SELECT count").
						WithArgs(username).
						WillReturnRows(rows)
				})

				It("propagates error", func() {
					_, err := uamDao.CreateUser(username, password)
					_, ok := err.(*myerr.ClientError)
					Expect(ok).To(Equal(true))
				})
			})

			Context("and user doesnt exist", func() {
				BeforeEach(func() {
					username, password = "test-user", "test-pass"
					rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
					mock.ExpectQuery("SELECT count").
						WithArgs(username).
						WillReturnRows(rows)
				})

				Context("and creation query is successful", func() {
					BeforeEach(func() {
						mock.ExpectQuery("INSERT INTO").
							WithArgs(Any{}, Any{}, Any{}, username, password). // driver.NamedValue - {Name: Ordinal:1 Value:2020-12-28 01:22:59.344298 +0200 EET}"
							WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
					})

					It("succeeds", func() {
						id, err := uamDao.CreateUser(username, password)
						Expect(err).NotTo(HaveOccurred())
						Expect(id).To(Equal(uint(1)))
					})
				})

				Context("and creation query fails", func() {
					BeforeEach(func() {
						mock.ExpectQuery("INSERT INTO").
							WithArgs(Any{}, Any{}, Any{}, username, password). // driver.NamedValue - {Name: Ordinal:1 Value:2020-12-28 01:22:59.344298 +0200 EET}"
							WillReturnError(fmt.Errorf("some error"))
					})

					It("propagates error", func() {
						_, err := uamDao.CreateUser(username, password)
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
				mock.ExpectQuery("SELECT count").
					WithArgs(id).
					WillReturnError(fmt.Errorf("some error"))
			})

			It("propagates error", func() {
				err := uamDao.DeleteUser(id)
				_, ok := err.(*myerr.ServerError)
				Expect(ok).To(Equal(true))
			})
		})

		When("request if user exists is succesful", func() {
			Context("and user does not exist", func() {
				BeforeEach(func() {
					rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
					mock.ExpectQuery("SELECT count").
						WithArgs(id).
						WillReturnRows(rows)
				})

				It("propagates error", func() {
					err := uamDao.DeleteUser(id)
					_, ok := err.(*myerr.ItemNotFoundError)
					Expect(ok).To(Equal(true))
				})
			})

			Context("and user exists", func() {
				BeforeEach(func() {
					rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
					mock.ExpectQuery("SELECT count").
						WithArgs(id).
						WillReturnRows(rows)
				})
				Context("and deletion query is successful", func() {
					BeforeEach(func() {
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
						mock.ExpectExec("DELETE FROM \"users\"").
							WithArgs(id). // driver.NamedValue - {Name: Ordinal:1 Value:2020-12-28 01:22:59.344298 +0200 EET}"
							WillReturnError(fmt.Errorf("some error"))
					})

					It("propagates error", func() {
						err := uamDao.DeleteUser(id)
						_, ok := err.(*myerr.ServerError)
						Expect(ok).To(Equal(true))
					})
				})
			})
		})
	})

	Context("CheckIfUserExists",func(){
		const (
			username = "username"
			password = "password"
		)

		When("request if user exists fails", func() {
			BeforeEach(func() {
				mock.ExpectQuery("SELECT count").
					WithArgs(username,password).
					WillReturnError(fmt.Errorf("some error"))
			})

			It("propagates error", func() {
				_, err := uamDao.CheckIfUserExists(username,password)
				_, ok := err.(*myerr.ServerError)
				Expect(ok).To(Equal(true))
			})
		})

		When("request if user exists is succesful", func() {
			Context("and user does not exist", func() {
				BeforeEach(func() {
					rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
					mock.ExpectQuery("SELECT count").
						WithArgs(username,password).
						WillReturnRows(rows)
				})

				It("succeds", func() {
					result, err := uamDao.CheckIfUserExists(username,password)
					Expect(err).NotTo(HaveOccurred())
					Expect(result).To(Equal(false))
				})
			})

			Context("and user exists", func() {
				BeforeEach(func() {
					rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
					mock.ExpectQuery("SELECT count").
						WithArgs(username,password).
						WillReturnRows(rows)
				})

				It("succeds", func() {
					result, err := uamDao.CheckIfUserExists(username,password)
					Expect(err).NotTo(HaveOccurred())
					Expect(result).To(Equal(true))
				})
			})
		})
	})
})
