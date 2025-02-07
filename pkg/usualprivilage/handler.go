package admin

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"myproject/pkg/common/utility"
	"myproject/pkg/handlers"
	"myproject/pkg/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
	services "myproject/pkg/client"
	db "myproject/pkg/database"
	"net/http"
)

type Handler struct {
	service   Service
	services  services.Services
	adminjw   middleware.MiddlewareJWT
	validator *validator.Validate
}

func NewHandler(service Service, srv services.Services, adTK middleware.MiddlewareJWT, validate *validator.Validate) *Handler {

	return &Handler{
		service:   service,
		services:  srv,
		adminjw:   adTK,
		validator: validate,
	}
}
func (h *Handler) MountRoutes(app *fiber.App) {
	//applicantApi := engine.Group(basePath)
	applicantApi := app.Group("/common")
	applicantApi.Post("/register", h.Register)
	applicantApi.Post("/login", h.Login)
	applicantApi.Post("/otpLogin", h.OtpLogin)

	//applicantApi.Use(h.adminjw.AdminAuthMiddleware())
	//{
	//applicantApi.Post("/deleteBlog", h.DeleteBlog)

	applicantApi.Get("/listBlog", h.Listing)

	//}
}

func (h *Handler) respondWithError(c *fiber.Ctx, code int, msg interface{}) error {
	return c.Status(code).JSON(fiber.Map{
		"msg": msg,
	})
}

func (h *Handler) respondWithData(c *fiber.Ctx, code int, message interface{}, data interface{}) error {
	if data == nil {
		data = "Successfully done"
	}
	return c.Status(code).JSON(fiber.Map{
		"msg":  message,
		"data": data,
	})
}
func (h *Handler) Listing(c *fiber.Ctx) error {
	fmt.Println("inside the app")
	ctx := c.Context()
	var conditions []db.WhereCondition
	pageCount := c.QueryInt("page")
	if pageCount < 1 {
		pageCount = 1
	}
	if c.Query("id") != "" {
		conditions = append(conditions, db.WhereCondition{
			Key:       "id",
			Value:     c.Query("id"),
			Condition: "=",
			Table:     "blogs",
			Joins:     "",
		})
	}
	if c.Query("user_id") != "" {
		conditions = append(conditions, db.WhereCondition{
			Key:       "id",
			Value:     c.Query("user_id"),
			Condition: "=",
			Table:     "users",
			Joins:     "",
		})
	}
	if c.Query("category") != "" {
		conditions = append(conditions, db.WhereCondition{
			Key:       "name",
			Value:     "'" + c.Query("category") + "'",
			Condition: "=",
			Table:     "categories",
			Joins:     "JOIN user_categories ON users.id = user_categories.user_id JOIN categories ON user_categories.category_id = categories.id ",
		})
	}

	if c.Query("name") != "" {
		conditions = append(conditions, db.WhereCondition{
			Key:       "title",
			Value:     "'%" + c.Query("name") + "%'",
			Condition: "ILIKE",
			Table:     "blogs",
		})
	}

	blogs, err := h.service.Listing(ctx, conditions, pageCount)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}

	return h.respondWithData(c, http.StatusOK, "success", blogs)

}

func (h *Handler) Register(c *fiber.Ctx) error {

	var request handlers.Register
	if err := utility.ParseAndValidate(c, &request, h.validator); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}

	if utility.IsValidPassword(request.Password) {
		return h.respondWithError(c, http.StatusBadRequest, map[string]interface{}{"invalid-request": "is not a valid password"})
	}
	log.Println("this is sucesssss")
	ctx := c.Context()
	if err := h.service.Register(ctx, request); err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	otp, err := h.services.SendEmailWithOTP(request.Email)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "error in sending otp"})

	}
	err = db.SetRedis(request.Email, otp, time.Minute*5)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "error in saving otp"})

	}
	_, err = db.GetRedis(request.Email)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "error in retriving redis data"})

	}

	return h.respondWithData(c, http.StatusOK, "success", nil)
}
func (h *Handler) Login(c *fiber.Ctx) error {
	// Parse request body into VendorRegisterRequest
	fmt.Println("this is in the handler Register")
	var request handlers.Login
	if err := utility.ParseAndValidate(c, &request, h.validator); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}

	ctx := c.Context()
	if err := h.service.Login(ctx, request); err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	fmt.Println("this is in the handler Register")
	fmt.Println("this is in the handler Register")
	token, err := h.adminjw.GenerateAdminToken(request.Email)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"token-generation": err.Error()})
	}

	fmt.Println("User logged in successfully")
	return h.respondWithData(c, http.StatusOK, "success", map[string]string{"token": token})
}
func (h *Handler) OtpLogin(c *fiber.Ctx) error {
	// Parse request body into VendorRegisterRequest
	fmt.Println("this is in the handler OtpLogin")
	var request handlers.Otp

	if err := utility.ParseAndValidate(c, &request, h.validator); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}
	fmt.Println("this is request", request)

	// Respond with success
	storedData, err := db.GetRedis(request.Email)
	fmt.Println("this is the keyy!!!!!", storedData, err)
	if storedData != request.Otp {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "wrong otp"})

	}
	ctx := c.Context()
	if err := h.service.OtpLogin(ctx, request); err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return h.respondWithData(c, http.StatusOK, "success", nil)
}
