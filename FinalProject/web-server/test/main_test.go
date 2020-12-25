package main_test

import (
	. "github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/main"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validator", func() {
	Describe("username validation", func() {
		When("username is invalid", func() {
			Context("username less than 7 symbols", func() {
				It("returns error", func() {
					err := ValidateUsername("example")
					Expect(err).To(HaveOccurred())
				})
			})

			Context("username greather than 20 symbols", func() {
				It("returns error", func() {
					err := ValidateUsername("somerandomlygeneratedusername")
					Expect(err).To(HaveOccurred())
				})
			})

			Context("username begins with number", func() {
				It("returns error", func() {
					err := ValidateUsername("1someusername")
					Expect(err).To(HaveOccurred())
				})
			})

			Context("username begins with non char symbol", func() {
				It("returns error", func() {
					err := ValidateUsername("9someusername")
					Expect(err).To(HaveOccurred())
				})
			})

			Context("username should contains a special character(except _ and -)", func() {
				It("returns error", func() {
					err := ValidateUsername("nonumberusername!")
					Expect(err).To(HaveOccurred())
				})
			})
		})
		When("username is valid", func() {
			It("succeeds", func() {
				err := ValidateUsername("valid-username")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
	Describe("password validation", func() {
		When("password is invalid", func() {
			Context("password less than 10 symbols", func() {
				It("returns error", func() {
					err := ValidatePassword("1random~")
					Expect(err).To(HaveOccurred())
				})
			})

			Context("password doesnt contain special symbol", func() {
				It("returns error", func() {
					err := ValidatePassword("1someusername")
					Expect(err).To(HaveOccurred())
				})
			})

			Context("password doesnt contain digit symbol", func() {
				It("returns error", func() {
					err := ValidatePassword("password~")
					Expect(err).To(HaveOccurred())
				})
			})
		})
		When("password is valid", func() {
			When("password is valid", func() {
				It("succeeds", func() {
					err := ValidatePassword("1validpassword~")
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})
	})
})
