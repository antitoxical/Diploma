package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"telegramshop_backend/internal/handler"
	adminsRepo "telegramshop_backend/internal/repository/admins"
	"telegramshop_backend/internal/repository/basket"
	"telegramshop_backend/internal/repository/categories"
	"telegramshop_backend/internal/repository/delivery"
	"telegramshop_backend/internal/repository/favorites"
	"telegramshop_backend/internal/repository/firms"
	"telegramshop_backend/internal/repository/notifications"
	"telegramshop_backend/internal/repository/orders"
	"telegramshop_backend/internal/repository/payments"
	"telegramshop_backend/internal/repository/prices"
	"telegramshop_backend/internal/repository/products"
	"telegramshop_backend/internal/repository/promocodes"
	"telegramshop_backend/internal/repository/reviews"
	"telegramshop_backend/internal/repository/tags"
	"telegramshop_backend/internal/repository/users"

	adminsService "telegramshop_backend/internal/service/admins"
	basketService "telegramshop_backend/internal/service/basket"
	categoriesService "telegramshop_backend/internal/service/categories"
	deliveryService "telegramshop_backend/internal/service/delivery"
	favoritesService "telegramshop_backend/internal/service/favorites"
	firmsService "telegramshop_backend/internal/service/firms"
	notificationsService "telegramshop_backend/internal/service/notifications"
	ordersService "telegramshop_backend/internal/service/orders"
	paymentsService "telegramshop_backend/internal/service/payments"
	pricesService "telegramshop_backend/internal/service/prices"
	productsService "telegramshop_backend/internal/service/products"
	promocodesService "telegramshop_backend/internal/service/promocodes"
	reviewsService "telegramshop_backend/internal/service/reviews"
	tagsService "telegramshop_backend/internal/service/tags"
	usersService "telegramshop_backend/internal/service/users"

	"telegramshop_backend/pkg/postgres"
	"telegramshop_backend/pkg/telegram"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	_ "telegramshop_backend/docs"

	_ "github.com/lib/pq"
)

func main() {
	db, err := postgres.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := users.NewRepository(db)
	basketRepo := basket.NewRepository(db)
	favoritesRepo := favorites.NewRepository(db)
	ordersRepo := orders.NewRepository(db)
	productsRepo := products.NewRepository(db)
	pricesRepo := prices.NewRepository(db)
	categoriesRepo := categories.NewRepository(db)
	firmsRepo := firms.NewRepository(db)
	reviewsRepo := reviews.NewRepository(db)
	notificationsRepo := notifications.NewRepository(db)
	adminRepo := adminsRepo.NewRepository(db)
	tagsRepo := tags.NewRepository(db)
	deliveryRepo := delivery.NewRepository(db)
	paymentsRepo := payments.NewRepository(db)
	promocodesRepo := promocodes.NewRepository(db)

	userService := usersService.NewService(userRepo)
	basketService := basketService.NewService(basketRepo)
	favoritesService := favoritesService.NewService(favoritesRepo)
	pricesService := pricesService.NewService(pricesRepo)
	categoriesService := categoriesService.NewService(categoriesRepo)
	firmsService := firmsService.NewService(firmsRepo)
	reviewsService := reviewsService.NewService(reviewsRepo, userRepo)
	notificationsService := notificationsService.NewService(notificationsRepo, userRepo)
	adminService := adminsService.NewService(adminRepo)
	tagsService := tagsService.NewService(tagsRepo)
	deliveryService := deliveryService.NewService(deliveryRepo)
	paymentsService := paymentsService.NewService(paymentsRepo)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "change-me-in-production"
		log.Println("WARNING: JWT_SECRET not set, using insecure default")
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Println("WARNING: BOT_TOKEN not set, Telegram initData validation will reject all requests")
	}

	webAppURL := os.Getenv("WEBAPP_URL")
	adminPanelURL := os.Getenv("ADMIN_PANEL_URL")
	if adminPanelURL == "" {
		adminPanelURL = "http://localhost:5173"
	}
	supportContact := os.Getenv("SUPPORT_CONTACT")

	var adminChatID int64
	if v := os.Getenv("ADMIN_TELEGRAM_CHAT_ID"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			adminChatID = id
		} else {
			log.Printf("WARNING: invalid ADMIN_TELEGRAM_CHAT_ID: %v", err)
		}
	}

	var bot *telegram.BotClient
	if botToken != "" {
		bot = telegram.NewBotClient(botToken, webAppURL, adminPanelURL, supportContact, adminChatID)
	}

	productsService := productsService.NewService(productsRepo, favoritesRepo, bot)
	promocodesService := promocodesService.NewService(promocodesRepo, userRepo, bot)
	ordersService := ordersService.NewService(ordersRepo, userRepo, productsRepo, promocodesRepo, notificationsService, bot)

	h := handler.NewHandler(
		db,
		userService,
		favoritesService,
		basketService,
		ordersService,
		firmsService,
		pricesService,
		categoriesService,
		productsService,
		reviewsService,
		notificationsService,
		adminService,
		tagsService,
		deliveryService,
		paymentsService,
		promocodesService,
		jwtSecret,
	)

	app := fiber.New(fiber.Config{
		BodyLimit: 12 * 1024 * 1024,
	})
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowMethods:  "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:  "Origin,Content-Type,Accept,Authorization",
		ExposeHeaders: "Content-Length",
	}))

	prometheus := fiberprometheus.New("telegramshop")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	h.InitRouter(app, botToken)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if bot != nil {
		go bot.StartPolling(ctx)
	}

	go func() {
		port := os.Getenv("SERVER_PORT")
		if port == "" {
			port = "8080"
		}
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Fiber Listen error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down gracefully...")
	_ = app.Shutdown()
}
