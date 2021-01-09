package middleware_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"example.com/user/web-server/api/common/response"
	"example.com/user/web-server/internal/auth"
	authMock "example.com/user/web-server/internal/auth/auth_mocks"
	mw "example.com/user/web-server/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func setupRouter(filter mw.AuthzFilter) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/protected").Use(filter.Authz())
	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
	})
	return r
}

var _ = Describe("AuthzFilter", func() {
	var (
		router     *gin.Engine
		recorder   *httptest.ResponseRecorder
		jwtCreator *authMock.MockJwtCreator
	)

	BeforeEach(func() {
		controller := gomock.NewController(GinkgoT())

		jwtCreator = authMock.NewMockJwtCreator(controller)
		filter := mw.NewAuthzFilterImpl(jwtCreator)
		router = setupRouter(filter)
		recorder = httptest.NewRecorder()
	})

	Context("Authz()", func() {
		When("request to protected resource is sent", func() {
			var req *http.Request

			BeforeEach(func() {
				req, _ = http.NewRequest("GET", "/protected/ping", nil)
			})

			Context("and there isnt an Authorization header", func() {
				It("returns error", func() {
					router.ServeHTTP(recorder, req)
					Expect(recorder.Code).To(Equal(http.StatusForbidden))
					body := response.ErrorResponse{}
					json.Unmarshal([]byte(recorder.Body.String()), &body)
					Expect(body.ErrorCode).To(Equal(http.StatusForbidden))
					Expect(body.ErrorMsg).To(Equal("No Authorization header provided"))
				})
			})

			Context("and Authorization header is set", func() {
				Context("with invalid formatted token", func() {
					BeforeEach(func() {
						req.Header.Set("Authorization", "Some token")
					})
					It("returns error response", func() {
						router.ServeHTTP(recorder, req)
						Expect(recorder.Code).To(Equal(http.StatusBadRequest))
						body := response.ErrorResponse{}
						json.Unmarshal([]byte(recorder.Body.String()), &body)
						Expect(body.ErrorCode).To(Equal(http.StatusBadRequest))
						Expect(body.ErrorMsg).To(Equal("Incorrect Format of Authorization Token"))
					})
				})

				Context("with jwt token", func() {

					BeforeEach(func() {
						req.Header.Set("Authorization", "Bearer sometoken")
					})
					Context("and has wrong credentials", func() {
						BeforeEach(func() {
							jwtCreator.EXPECT().
								ValidateToken(gomock.Any()).
								Return(nil, errors.New("test error"))
						})

						It("returns error response", func() {
							router.ServeHTTP(recorder, req)
							Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
							body := response.ErrorResponse{}
							json.Unmarshal([]byte(recorder.Body.String()), &body)
							Expect(body.ErrorCode).To(Equal(http.StatusUnauthorized))
							Expect(body.ErrorMsg).To(Equal("Invalid Authorization token"))
						})
					})
					Context("and the credentials are legit", func() {
						BeforeEach(func() {
							jwtCreator.EXPECT().
								ValidateToken(gomock.Any()).
								Return(&auth.JwtClaim{
									UserID: 1,
								}, nil)
						})

						It("returns error response", func() {
							router.ServeHTTP(recorder, req)
							Expect(recorder.Code).To(Equal(http.StatusOK))
						})

					})
				})
			})

		})
	})

})
