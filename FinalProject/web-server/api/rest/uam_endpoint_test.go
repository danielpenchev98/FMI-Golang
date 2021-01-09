package rest_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"example.com/user/web-server/api/common/response"
	"example.com/user/web-server/api/rest"
	"example.com/user/web-server/internal/auth/auth_mocks"
	"example.com/user/web-server/internal/db/dao/uam_dao_mocks"
	"example.com/user/web-server/internal/db/models"
	myerr "example.com/user/web-server/internal/error"
	"example.com/user/web-server/internal/validator/validator_mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
)

func setupRouter(uamRest rest.UamEndpoint, userID uint) *gin.Engine {
	r := gin.Default()

	public := r.Group("/public")
	{
		public.POST("/registration", uamRest.CreateUser)
		public.POST("/login", uamRest.Login)
	}
	protected := r.Group("/protected").Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	}) // check how to do it
	{
		protected.DELETE("/deletion", uamRest.DeleteUser)
	}
	return r
}

func assertErrorResponse(recorder *httptest.ResponseRecorder, expStatusCode int, expMessage string) {
	Expect(recorder.Code).To(Equal(expStatusCode))
	body := response.ErrorResponse{}
	json.Unmarshal([]byte(recorder.Body.String()), &body)
	Expect(body.ErrorCode).To(Equal(expStatusCode))
	Expect(body.ErrorMsg).To(ContainSubstring(expMessage))
}

var _ = Describe("UamEndpoint", func() {
	var (
		router     *gin.Engine
		recorder   *httptest.ResponseRecorder
		jwtCreator *auth_mocks.MockJwtCreator
		uamDAO     *uam_dao_mocks.MockUamDAO
		validator  *validator_mocks.MockValidator
		req *http.Request
	)

	const userID = 1

	BeforeEach(func() {
		controller := gomock.NewController(GinkgoT())
		uamDAO = uam_dao_mocks.NewMockUamDAO(controller)
		jwtCreator = auth_mocks.NewMockJwtCreator(controller)
		validator = validator_mocks.NewMockValidator(controller)
		uamRest := rest.NewUamEndPointImpl(uamDAO, jwtCreator, validator)

		router = setupRouter(uamRest, userID)
		recorder = httptest.NewRecorder()
	})

	Context("CreateUser", func() {
		When("creation request is sent", func() {
			var reqBody *rest.RequestWithCredentials

			const (
				username = "username"
				password = "password"
			)

			BeforeEach(func() {
				reqBody = &rest.RequestWithCredentials{
					Username: username,
					Password: password,
				}
			})

			Context("with non-json body", func() {

				BeforeEach(func() {
					validator.EXPECT().
						ValidateUsername(username).
						Times(0)

					validator.EXPECT().
						ValidatePassword(password).
						Times(0)

					uamDAO.EXPECT().
						CreateUser(gomock.Any(), gomock.Any()).
						Times(0)

					req, _ = http.NewRequest("POST", "/public/registration", strings.NewReader("test"))
				})

				It("returns bad request", func() {
					router.ServeHTTP(recorder, req)
					assertErrorResponse(recorder, http.StatusBadRequest, "Invalid json body")
				})
			})

			Context("with json body", func() {

				BeforeEach(func() {
					jsonBody, _ := json.Marshal(*reqBody)
					req, _ = http.NewRequest("POST", "/public/registration", bytes.NewBuffer(jsonBody))
					req.Header.Set("Content-Type", "application/json")
				})

				Context("with invalid format of username", func() {
					BeforeEach(func() {
						validator.EXPECT().
							ValidateUsername(gomock.Any()).
							Return(myerr.NewClientError("test-error"))

						uamDAO.EXPECT().
							CreateUser(gomock.Any(), gomock.Any()).
							Times(0)
					})

					It("returns bad request", func() {
						router.ServeHTTP(recorder, req)
						assertErrorResponse(recorder, http.StatusBadRequest, "test-error")
					})
				})

				Context("with valid formatted username", func() {
					Context("and invalid formatted password", func() {
						BeforeEach(func() {
							gomock.InOrder(
								validator.EXPECT().
									ValidateUsername(username).
									Return(nil),

								validator.EXPECT().
									ValidatePassword(password).
									Return(myerr.NewClientError("test-error")),
							)

							uamDAO.EXPECT().
								CreateUser(gomock.Any(), gomock.Any()).
								Times(0)
						})

						It("returns bad request", func() {
							router.ServeHTTP(recorder, req)
							assertErrorResponse(recorder, http.StatusBadRequest, "test-error")
						})
					})

					Context("and valid formatted password", func() {
						Context("and CreateUser request fails", func() {
							BeforeEach(func() {
								gomock.InOrder(
									validator.EXPECT().
										ValidateUsername(username).
										Return(nil),

									validator.EXPECT().
										ValidatePassword(password).
										Return(nil),

									uamDAO.EXPECT().
										CreateUser(username, gomock.Any()). // Any() because bscrypt cannot be mocked easily -> Could be wrapped
										Return(myerr.NewServerError("test-error")),
								)
							})

							It("returns bad request", func() {
								router.ServeHTTP(recorder, req)
								assertErrorResponse(recorder, http.StatusInternalServerError, "Problem with the server, please try again later")
							})
						})

						Context("and CreateUser request succeeds", func() {
							Context("and user already exists", func() {
								BeforeEach(func() {
									gomock.InOrder(
										validator.EXPECT().
											ValidateUsername(username).
											Return(nil),

										validator.EXPECT().
											ValidatePassword(password).
											Return(nil),

										uamDAO.EXPECT().
											CreateUser(username, gomock.Any()). // Any() because bscrypt cannot be mocked easily -> Could be wrapped
											Return(myerr.NewClientError("test-error")),
									)
								})

								It("returns bad request", func() {
									router.ServeHTTP(recorder, req)
									assertErrorResponse(recorder, http.StatusBadRequest, "test-error")
								})
							})

							Context("and user doesnt already exist", func() {
								BeforeEach(func() {
									gomock.InOrder(
										validator.EXPECT().
											ValidateUsername(username).
											Return(nil),

										validator.EXPECT().
											ValidatePassword(password).
											Return(nil),

										uamDAO.EXPECT().
											CreateUser(username, gomock.Any()). // Any() because bscrypt cannot be mocked easily -> Could be wrapped
											Return(nil),
									)
								})

								It("returns successfull response", func() {
									router.ServeHTTP(recorder, req)

									Expect(recorder.Code).To(Equal(http.StatusCreated))
									body := response.BasicResponse{}
									json.Unmarshal([]byte(recorder.Body.String()), &body)
									Expect(body.Status).To(Equal(http.StatusCreated))
								})
							})
						})
					})
				})
			})
		})
	})

	Context("Login", func() {
		When("login request is sent", func() {
			var reqBody *rest.RequestWithCredentials

			const (
				username = "username"
				password = "password"
			)

			BeforeEach(func() {
				reqBody = &rest.RequestWithCredentials{
					Username: username,
					Password: password,
				}
			})

			Context("with non-json body", func() {
				BeforeEach(func() {
					uamDAO.EXPECT().
						GetUser(gomock.Any()).
						Times(0)

					jwtCreator.EXPECT().
						GenerateToken(gomock.Any()).
						Times(0)

					req, _ = http.NewRequest("POST", "/public/login", strings.NewReader("test"))
				})

				It("returns bad request", func() {
					router.ServeHTTP(recorder, req)
					assertErrorResponse(recorder, http.StatusBadRequest, "Invalid json body")
				})
			})

			Context("with json body", func() {
				BeforeEach(func() {
					jsonBody, _ := json.Marshal(*reqBody)
					req, _ = http.NewRequest("POST", "/public/login", bytes.NewBuffer(jsonBody))
					req.Header.Set("Content-Type", "application/json")
				})

				Context("and request to check if user exist fail", func() {
					BeforeEach(func() {
						uamDAO.EXPECT().
							GetUser(username).
							Return(models.User{}, myerr.NewServerError("test-error"))

						jwtCreator.EXPECT().
							GenerateToken(gomock.Any()).
							Times(0)
					})

					It("returns internal server error response", func() {
						router.ServeHTTP(recorder, req)
						assertErrorResponse(recorder, http.StatusInternalServerError, "Problem with the server, please try again later")
					})
				})

				Context("and request to check if user exist succeeds", func() {
					Context("and user doesnt exist", func() {
						BeforeEach(func() {
							uamDAO.EXPECT().
								GetUser(username).
								Return(models.User{}, myerr.NewClientError("test-error"))

							jwtCreator.EXPECT().
								GenerateToken(gomock.Any()).
								Times(0)
						})

						It("returns bad request error response", func() {
							router.ServeHTTP(recorder, req)
							assertErrorResponse(recorder, http.StatusBadRequest, "test-error")
						})
					})

					Context("and user exist", func() {
						var user models.User

						BeforeEach(func() {
							encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

							user = models.User{
								Username: username,
								Password: string(encryptedPass),
							}
							user.ID = 1
						})

						Context("and token generation fails", func() {
							BeforeEach(func() {

								gomock.InOrder(
									uamDAO.EXPECT().
										GetUser(user.Username).
										Return(user, nil),

									jwtCreator.EXPECT().
										GenerateToken(user.ID).
										Return("", errors.New("test-error")),
								)
							})

							It("returns internal server error response", func() {
								router.ServeHTTP(recorder, req)
								assertErrorResponse(recorder, http.StatusInternalServerError, "Problem with the server, please try again later")
							})
						})

						Context("and token generation succeeds", func() {
							const token = "token"

							BeforeEach(func() {
								gomock.InOrder(
									uamDAO.EXPECT().
										GetUser(user.Username).
										Return(user, nil),

									jwtCreator.EXPECT().
										GenerateToken(user.ID).
										Return(token, nil),
								)
							})

							It("returns successful response", func() {
								router.ServeHTTP(recorder, req)

								Expect(recorder.Code).To(Equal(http.StatusCreated))
								body := rest.LoginResponse{}
								json.Unmarshal([]byte(recorder.Body.String()), &body)
								Expect(body.Status).To(Equal(http.StatusCreated))
								Expect(body.Token).To(Equal(token))
							})
						})
					})
				})
			})
		})
	})

	Context("DeleteUser", func() {
		When("login request is sent and authentication passes", func() {

			BeforeEach(func() {
				req, _ = http.NewRequest("DELETE", "/protected/deletion", nil)
				req.Header.Set("Authorization", "Bearer sometoken")
			})

			Context("operation of deleting user from db fails", func() {
				BeforeEach(func() {
					uamDAO.EXPECT().
						DeleteUser(uint(userID)).
						Return(myerr.NewServerError("test-error"))
				})

				It("returns internal server error response", func() {
					router.ServeHTTP(recorder, req)
					assertErrorResponse(recorder, http.StatusInternalServerError, "Problem with the server, please try again later")
				})
			})

			Context("operation of deleting user from db succeeds", func() {
				Context("and user doesnt exist", func() {
					BeforeEach(func() {
						uamDAO.EXPECT().
							DeleteUser(uint(userID)).
							Return(myerr.NewItemNotFoundError("test-error"))
					})

					It("returns internal server error response", func() {
						router.ServeHTTP(recorder, req)
						assertErrorResponse(recorder, http.StatusInternalServerError, "Problem with the server, please try again later")
					})
				})

				Context("and user exists", func() {
					BeforeEach(func() {
						uamDAO.EXPECT().
							DeleteUser(uint(userID)).
							Return(nil)
					})

					It("returns success response", func() {
						router.ServeHTTP(recorder, req)

						Expect(recorder.Code).To(Equal(http.StatusOK))
						body := response.BasicResponse{}
						json.Unmarshal([]byte(recorder.Body.String()), &body)
						Expect(body.Status).To(Equal(http.StatusOK))
					})
				})

			})
		})
	})
})
