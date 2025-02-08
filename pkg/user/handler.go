package user

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"myproject/pkg/common/utility"
	"myproject/pkg/handlers"
	"myproject/pkg/middleware"
	"strconv"

	// db "myproject/pkg/database"
	services "myproject/pkg/client"
	"myproject/pkg/config"
	"net/http"

	// "time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service   Service
	services  services.Services
	adminjw   middleware.MiddlewareJWT
	validator *validator.Validate
	cnf       config.Config
}

func NewHandler(service Service, srv services.Services, adTK middleware.MiddlewareJWT, cnf config.Config, validate *validator.Validate) *Handler {

	return &Handler{
		service:   service,
		services:  srv,
		adminjw:   adTK,
		cnf:       cnf,
		validator: validate,
	}
}
func (h *Handler) MountRoutes(app *fiber.App) {

	applicantApi := app.Group("/user")

	applicantApi.Use(h.adminjw.AdminAuthMiddleware())
	{

		///list orders
		applicantApi.Post("/addBlog", h.AddBlog)
		applicantApi.Patch("/addBlog/:id", h.UpdateBlog)
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
func (h *Handler) UpdateBlog(c *fiber.Ctx) error {

	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid blog ID",
		})
	}
	var updateReq handlers.UpdateBlogRequest
	if err := c.BodyParser(&updateReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	ctx := c.Context()
	// Call the repository to update the blog
	if err := h.service.UpdateBlog(ctx, updateReq, c.Locals("username").(string), id); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
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

func (h *Handler) AddBlog(c *fiber.Ctx) error {

	var request handlers.AddBlog

	if err := utility.ParseAndValidate(c, &request, h.validator); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}
	ctx := c.Context()
	if err := h.service.AddBlog(ctx, request, c.Locals("username").(string)); err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	fmt.Println("this is in the handler Register")

	return h.respondWithData(c, http.StatusOK, "success", nil)
}
