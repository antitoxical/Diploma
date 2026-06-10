package handler

import (
	"telegramshop_backend/internal/service/admins"
	"telegramshop_backend/internal/service/basket"
	"telegramshop_backend/internal/service/categories"
	"telegramshop_backend/internal/service/delivery"
	"telegramshop_backend/internal/service/favorites"
	"telegramshop_backend/internal/service/firms"
	"telegramshop_backend/internal/service/notifications"
	"telegramshop_backend/internal/service/orders"
	"telegramshop_backend/internal/service/payments"
	"telegramshop_backend/internal/service/prices"
	"telegramshop_backend/internal/service/products"
	"telegramshop_backend/internal/service/promocodes"
	"telegramshop_backend/internal/service/reviews"
	"telegramshop_backend/internal/service/tags"
	"telegramshop_backend/internal/service/users"
	"telegramshop_backend/pkg/auth"
	"telegramshop_backend/pkg/telegram"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

var (
	apiPath     = "/api"
	versionPath = "v1"
	basePath    = apiPath + "/" + versionPath
)

type Handler struct {
	db                  *sqlx.DB
	userService         users.Service
	favoriteService     favorites.Service
	basketService       basket.Service
	orderService        orders.Service
	firmsService        firms.Service
	priceService        prices.Service
	categoryService     categories.Service
	productService      products.Service
	reviewService       reviews.Service
	notificationService notifications.Service
	adminService        admins.Service
	tagService          tags.Service
	deliveryService     delivery.Service
	paymentService      payments.Service
	promocodeService    promocodes.Service
	jwtSecret           string
}

func NewHandler(
	db *sqlx.DB,
	userService users.Service,
	favoriteService favorites.Service,
	basketService basket.Service,
	orderService orders.Service,
	firmsService firms.Service,
	priceService prices.Service,
	categoryService categories.Service,
	productService products.Service,
	reviewService reviews.Service,
	notificationService notifications.Service,
	adminService admins.Service,
	tagService tags.Service,
	deliveryService delivery.Service,
	paymentService payments.Service,
	promocodeService promocodes.Service,
	jwtSecret string,
) *Handler {
	return &Handler{
		db:                  db,
		userService:         userService,
		favoriteService:     favoriteService,
		basketService:       basketService,
		orderService:        orderService,
		firmsService:        firmsService,
		priceService:        priceService,
		categoryService:     categoryService,
		productService:      productService,
		reviewService:       reviewService,
		notificationService: notificationService,
		adminService:        adminService,
		tagService:          tagService,
		deliveryService:     deliveryService,
		paymentService:      paymentService,
		promocodeService:    promocodeService,
		jwtSecret:           jwtSecret,
	}
}

func (h *Handler) InitRouter(app *fiber.App, botToken string) {
	api := app.Group(basePath)

	api.Get("/products/search", h.SearchProducts)
	api.Get("/products/:id", h.GetProductByID)
	api.Get("/products", h.GetAllProducts)
	api.Get("/categories/roots", h.GetRootCategories)
	api.Get("/categories/:id/subcategories", h.GetSubcategories)
	api.Get("/categories/:id", h.GetCategoryByID)
	api.Get("/categories", h.GetAllCategories)
	api.Get("/firms/:id", h.GetFirmByID)
	api.Get("/firms", h.GetAllFirms)
	api.Get("/prices/:id", h.GetPriceByID)
	api.Get("/prices/product/:product_id", h.GetPricesByProductID)
	api.Get("/reviews/product/:product_id", h.GetReviewsByProduct)
	api.Get("/reviews/product/:product_id/stats", h.GetProductRatingStat)
	api.Get("/tags", h.GetAllTags)
	api.Get("/delivery-methods", h.GetAllDeliveryMethods)
	api.Get("/payment-methods", h.GetAllPaymentMethods)
	api.Get("/promocodes/validate", h.ValidatePromocode)

	api.Get("/uploads/images/:filename", h.ServeUploadedImage)

	api.Post("/admin/login", h.AdminLogin)

	admin := api.Group("/admin", auth.JWTMiddleware(h.jwtSecret), h.AdminAuditMiddleware())

	admin.Get("/admins", h.GetAllAdmins)
	admin.Get("/admins/:id", h.GetAdminByID)
	admin.Post("/admins", h.CreateAdmin)
	admin.Put("/admins/:id", h.UpdateAdmin)

	admin.Get("/reviews", h.GetAllReviews)

	admin.Get("/delivery-methods", h.GetAllDeliveryMethods)
	admin.Get("/payment-methods", h.GetAllPaymentMethods)

	admin.Get("/products", h.GetAllProducts)
	admin.Get("/products/:id", h.GetProductByID)
	admin.Get("/categories", h.GetAllCategories)
	admin.Get("/categories/:id", h.GetCategoryByID)
	admin.Get("/firms", h.GetAllFirms)
	admin.Get("/firms/:id", h.GetFirmByID)
	admin.Get("/tags", h.GetAllTags)
	admin.Get("/users", h.GetAllUsers)
	admin.Get("/dashboard/stats", h.GetDashboardStats)
	admin.Get("/orders/all", h.GetAllOrders)
	admin.Get("/orders/:id", h.GetOrder)
	admin.Patch("/orders/:id/status", h.UpdateOrderStatus)

	admin.Post("/uploads/image", h.UploadAdminImage)

	admin.Post("/products", h.CreateProduct)
	admin.Put("/products/:id", h.UpdateProduct)
	admin.Delete("/products/:id", h.DeleteProduct)
	admin.Put("/products/:id/image", h.AddProductImage)
	admin.Delete("/products/images/:image_id", h.RemoveProductImage)
	admin.Put("/products/:id/images", h.SetProductImages)
	admin.Patch("/products/:id/stock", h.UpdateStock)
	admin.Post("/products/:id/tags/:tag_id", h.AddProductTag)
	admin.Delete("/products/:id/tags/:tag_id", h.RemoveProductTag)

	admin.Post("/categories", h.CreateCategory)
	admin.Put("/categories/:id", h.UpdateCategory)
	admin.Delete("/categories/:id", h.DeleteCategory)
	admin.Put("/categories/:id/image", h.SetCategoryImage)
	admin.Delete("/categories/:id/image", h.RemoveCategoryImage)

	admin.Post("/firms", h.CreateFirm)
	admin.Put("/firms/:id", h.UpdateFirm)
	admin.Delete("/firms/:id", h.DeleteFirm)

	admin.Post("/prices", h.CreatePrice)
	admin.Put("/prices/:id", h.UpdatePrice)
	admin.Delete("/prices/:id", h.DeletePrice)
	admin.Delete("/prices/product/:product_id", h.DeletePricesByProductID)
	admin.Patch("/prices/:id/count", h.UpdatePriceCount)

	admin.Post("/tags", h.CreateTag)
	admin.Put("/tags/:id", h.UpdateTag)
	admin.Delete("/tags/:id", h.DeleteTag)

	admin.Post("/delivery-methods", h.CreateDeliveryMethod)
	admin.Put("/delivery-methods/:id", h.UpdateDeliveryMethod)
	admin.Delete("/delivery-methods/:id", h.DeleteDeliveryMethod)

	admin.Post("/payment-methods", h.CreatePaymentMethod)
	admin.Put("/payment-methods/:id", h.UpdatePaymentMethod)
	admin.Delete("/payment-methods/:id", h.DeletePaymentMethod)

	admin.Post("/promocodes", h.CreatePromocode)
	admin.Get("/promocodes", h.GetAllPromocodes)
	admin.Delete("/promocodes/:id", h.DeletePromocode)

	admin.Patch("/reviews/:id/moderate", h.ModerateReview)

	admin.Get("/audit-logs", h.GetAuditLogs)

	admin.Post("/notifications", h.CreateNotification)
	admin.Delete("/notifications/:id", h.DeleteNotification)

	user := api.Group("", telegram.InitDataMiddleware(botToken))

	user.Post("/users", h.CreateUser)
	user.Get("/users/:id", h.GetUser)
	user.Put("/users/:id", h.UpdateUserProfile)
	user.Delete("/users/:id", h.DeleteUser)

	user.Post("/addresses", h.CreateUserAddress)
	user.Get("/addresses/:user_id", h.GetUserAddresses)
	user.Delete("/addresses/:id", h.DeleteUserAddress)

	user.Post("/favorites", h.AddToFavorites)
	user.Get("/favorites/:user_id", h.GetUserFavorites)
	user.Delete("/favorites/:user_id/:product_id", h.RemoveFromFavorites)

	user.Post("/basket", h.AddToBasket)
	user.Get("/basket/:user_id", h.GetUserBasket)
	user.Put("/basket", h.UpdateBasketItem)
	user.Delete("/basket/:user_id/:product_id", h.RemoveFromBasket)

	user.Post("/orders", h.CreateOrder)
	user.Get("/orders/:id", h.GetOrder)
	user.Get("/orders/user/:user_id", h.GetUserOrders)
	user.Get("/orders/:id/history", h.GetOrderStatusHistory)

	user.Post("/reviews", h.UpsertReview)
	user.Get("/reviews/user/:user_id/product/:product_id", h.GetUserReview)
	user.Delete("/reviews/user/:user_id/product/:product_id", h.DeleteReview)

	user.Get("/notifications/:user_id", h.GetUserNotifications)
	user.Patch("/notifications/:id/read", h.MarkNotificationRead)
	user.Patch("/notifications/:user_id/read-all", h.MarkAllNotificationsRead)
}
