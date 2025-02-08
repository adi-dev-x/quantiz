package admin

import (
	"github.com/go-playground/validator/v10"
	"myproject/pkg/common/utility"
	"myproject/pkg/handlers"
	"myproject/pkg/middleware"
	"strconv"
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
	applicantApi := app.Group("/admin")
	applicantApi.Post("/register", h.Register)

	applicantApi.Use(h.adminjw.AdminAuthMiddleware())
	{
		applicantApi.Get("/deleteBlog/:id", h.DeleteBlog)

	}
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

func (h *Handler) Register(c *fiber.Ctx) error {

	var request handlers.Register
	if err := utility.ParseAndValidate(c, &request, h.validator); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}

	if utility.IsValidPassword(request.Password) {
		return h.respondWithError(c, http.StatusBadRequest, map[string]interface{}{"invalid-request": "is not a valid password"})
	}

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

func (h *Handler) DeleteBlog(c *fiber.Ctx) error {

	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid blog ID",
		})
	}

	ctx := c.Context()

	if err := h.service.DeleteBlog(ctx, c.Locals("username").(string), id); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}

	return h.respondWithData(c, http.StatusOK, "success", nil)
}
