package auth

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JWT authentication", func() {
	var jwtWrapper JwtWrapper
	const (
		secretKey  = "secret"
		issuer     = "issuer"
		expiration = 24
	)

	BeforeEach(func() {
		jwtWrapper = JwtWrapper{
			SecretKey:       secretKey,
			Issuer:          issuer,
			ExpirationHours: expiration,
		}
	})

	Context("GenerateToken", func() {
		When("creating a token", func() {
			It("succeeds", func() {
				_, err := jwtWrapper.GenerateToken(1)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Context("ValidateToken", func() {
		When("validating legal token", func() {
			var token string
			const userID = 1
			BeforeEach(func() {
				token, _ = jwtWrapper.GenerateToken(userID)
			})

			It("classifies the token as legal", func() {
				claims, err := jwtWrapper.ValidateToken(token)
				Expect(err).NotTo(HaveOccurred())
				Expect(claims.UserID).To(Equal(uint(userID)))
				Expect(claims.Issuer).To(Equal(issuer))
			})
		})
	})
})
