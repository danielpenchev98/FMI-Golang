package auth_test

import (
	"os"
	"strconv"

	"example.com/user/web-server/internal/auth"
	myerr "example.com/user/web-server/internal/error"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth module", func() {
	var jwtCreator auth.JwtCreator
	const (
		secret     = "secret"
		issuer     = "issuer"
		expiration = 24
	)

	BeforeEach(func() {
		os.Setenv("secret", secret)
		os.Setenv("issuer", issuer)
		os.Setenv("expiration", strconv.Itoa(expiration))
	})

	Context("NewJwtCreatorImpl", func() {
		When("when secret env variable is missing", func() {
			BeforeEach(func() {
				os.Unsetenv("secret")
			})

			AfterEach(func() {
				os.Setenv("secret", secret)
			})

			It("returns error", func() {
				_, err := auth.NewJwtCreatorImpl()
				Expect(err).To(HaveOccurred())
				_, ok := err.(*myerr.ServerError)
				Expect(ok).To(Equal(true))
			})
		})

		When("when issuer env variable is missing", func() {
			BeforeEach(func() {
				os.Unsetenv("issuer")
			})

			AfterEach(func() {
				os.Setenv("issuer", issuer)
			})

			It("returns error", func() {
				_, err := auth.NewJwtCreatorImpl()
				Expect(err).To(HaveOccurred())
				_, ok := err.(*myerr.ServerError)
				Expect(ok).To(Equal(true))
			})
		})

		When("when expiration env variable is missing", func() {
			BeforeEach(func() {
				os.Unsetenv("expiration")
			})

			AfterEach(func() {
				os.Setenv("expiration", strconv.Itoa(expiration))
			})

			It("returns error", func() {
				_, err := auth.NewJwtCreatorImpl()
				Expect(err).To(HaveOccurred())
				_, ok := err.(*myerr.ServerError)
				Expect(ok).To(Equal(true))
			})
		})
	})
	Context("JwtCreator", func() {
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
			When("validating legal token", func() {
				var token string
				const userID = 1
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
		})
	})

})
