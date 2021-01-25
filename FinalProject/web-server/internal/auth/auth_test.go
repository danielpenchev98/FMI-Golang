package auth_test

import (
	"os"
	"strconv"
	"time"

	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/internal/auth"
	myerr "github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/internal/error"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth module", func() {

	const (
		secret     = "secret"
		issuer     = "issuer"
		expiration = 24
	)

	BeforeEach(func() {
		os.Clearenv()
	})

	Context("NewJwtCreatorImpl", func() {
		When("Creating new Jwt creator", func() {
			Context("when secret env variable is missing", func() {
				It("returns error", func() {
					_, err := auth.NewJwtCreatorImpl()
					Expect(err).To(HaveOccurred())
					_, ok := err.(*myerr.ServerError)
					Expect(ok).To(Equal(true))
				})
			})

			Context("when secret env variable exists", func() {
				BeforeEach(func() {
					os.Setenv("secret", secret)
				})
				AfterEach(func() {
					os.Unsetenv("secret")
				})

				Context("when issuer env variable is missing", func() {
					It("returns error", func() {
						_, err := auth.NewJwtCreatorImpl()
						Expect(err).To(HaveOccurred())
						_, ok := err.(*myerr.ServerError)
						Expect(ok).To(Equal(true))
					})
				})

				Context("when issuer env variable exists", func() {
					BeforeEach(func() {
						os.Setenv("issuer", issuer)
					})

					AfterEach(func() {
						os.Unsetenv("issuer")
					})

					Context("when expiration env variable is missing", func() {
						It("returns error", func() {
							_, err := auth.NewJwtCreatorImpl()
							Expect(err).To(HaveOccurred())
							_, ok := err.(*myerr.ServerError)
							Expect(ok).To(Equal(true))
						})
					})

					Context("when expiration env vairable exitst", func() {
						Context("and expiration variable is in illegal format", func() {
							BeforeEach(func() {
								os.Setenv("expiration", "wrong-format")
							})

							AfterEach(func() {
								os.Unsetenv("expiration")
							})

							It("returns error", func() {
								_, err := auth.NewJwtCreatorImpl()
								Expect(err).To(HaveOccurred())
								_, ok := err.(*myerr.ServerError)
								Expect(ok).To(Equal(true))
							})
						})

						Context("and expiration variable is in legal format", func() {
							BeforeEach(func() {
								os.Setenv("expiration", strconv.Itoa(expiration))
							})

							AfterEach(func() {
								os.Unsetenv("expiration")
							})

							It("succeeds", func() {
								actualResult, err := auth.NewJwtCreatorImpl()
								Expect(err).NotTo(HaveOccurred())
								expectedResult := &auth.JwtCreatorImpl{
									Secret:          secret,
									Issuer:          issuer,
									ExpirationHours: expiration,
								}
								Expect(actualResult).To(Equal(expectedResult))
							})
						})
					})
				})
			})
		})
	})
	Context("JwtCreator", func() {
		var jwtCreator auth.JwtCreator
		BeforeEach(func() {
			jwtCreator = &auth.JwtCreatorImpl{
				Secret:          secret,
				Issuer:          issuer,
				ExpirationHours: expiration,
			}
		})

		Context("GenerateToken", func() {
			When("creating a token", func() {
				It("succeeds", func() {
					_, err := jwtCreator.GenerateToken(1)
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

		Context("ValidateToken", func() {
			var token string
			const userID = 1

			When("validating legal token", func() {
				BeforeEach(func() {
					token, _ = jwtCreator.GenerateToken(userID)
				})

				It("classifies the token as legal", func() {
					claims, err := jwtCreator.ValidateToken(token)
					Expect(err).NotTo(HaveOccurred())
					Expect(claims.UserID).To(Equal(uint(userID)))
					Expect(claims.Issuer).To(Equal(issuer))
				})
			})

			When("validating expired token", func() {

				BeforeEach(func() {
					jwtCreator = &auth.JwtCreatorImpl{
						Secret:          secret,
						Issuer:          issuer,
						ExpirationHours: 0,
					}

					token, _ = jwtCreator.GenerateToken(userID)
				})

				It("returns error", func() {
					time.Sleep(1 * time.Second)
					_, err := jwtCreator.ValidateToken(token)
					Expect(err).To(HaveOccurred())
					_, ok := err.(*myerr.ClientError)
					Expect(ok).To(Equal(true))
				})
			})

			When("validating wrong type of token", func() {
				It("returns error", func() {
					_, err := jwtCreator.ValidateToken("wrong-type")
					Expect(err).To(HaveOccurred())
					_, ok := err.(*myerr.ClientError)
					Expect(ok).To(Equal(true))
				})
			})
		})
	})

})
