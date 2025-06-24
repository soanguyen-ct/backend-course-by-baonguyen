// You can edit this code!
// Click here and start typing.
package main

import (
	"ct-backend-course-baonguyen/config"
	"ct-backend-course-baonguyen/internal/controller"
	postgres "ct-backend-course-baonguyen/internal/storage/postgres"
	"ct-backend-course-baonguyen/internal/usecase"
	auth "ct-backend-course-baonguyen/pkg/auth"
	"ct-backend-course-baonguyen/pkg/bucket"
	"ct-backend-course-baonguyen/pkg/validator"

	_ "github.com/labstack/echo-jwt/v4"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	conf := config.Config{
		Port:           "27017",
		MongoURI:       "mongodb://localhost:27017",
		MongoDB:        "ct_backend_course",
		MongoCollImage: "images",
		MongoCollUser:  "users",
		// PostgreSQL config
		PostgresHost:     "localhost",
		PostgresPort:     "5432",
		PostgresUser:     "postgres",
		PostgresPassword: "postgres",
		PostgresDBName:   "ct_backend_course",
		PostgresSSLMode:  "disable",
		GoogleCredFile:   "path/to/your/google/credentials.json",
		GoogleBucketName: "your-google-bucket-name",
	}

	// demoDB := mongostore.MustDatabase(conf.MongoURI, conf.MongoDB)
	// userStore := mongostore.NewUserCollection(demoDB, conf.MongoCollUser)
	// imageStore := mongostore.NewImageCollection(demoDB, conf.MongoCollImage)

	// For PostgreSQL
	demoDB := postgres.MustDatabase(conf.PostgresConnString())
	userStore := postgres.NewUserRepository(demoDB.DB)
	imageStore := postgres.NewImageRepository(demoDB.DB)

	//imgBucket := bucket.MustNewGoogleStorageClient(context.TODO(), conf.GoogleBucketName, conf.GoogleCredFile)
	imgBucket := bucket.NewFake()

	uc := usecase.NewUseCase(imageStore, userStore, imgBucket)
	hdl := controller.NewHandler(uc)

	srv := newServer(hdl)
	if err := srv.Start(":8090"); err != nil {
		log.Error(err)
	}
}

func newServer(hdl *controller.Handler) *echo.Echo {
	e := echo.New()
	e.Validator = validator.NewCustomValidator()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	public := e.Group("/api/public")
	private := e.Group("/api/private")
	private.Use(auth.AuthMiddleware(), auth.ExtractUserNameFn)

	public.POST("/register", hdl.Register)
	public.POST("/login", hdl.Login)

	private.GET("/self", hdl.Self)
	private.POST("/upload", hdl.UploadImage)
	private.PUT("/change-password", hdl.ChangePassword)

	return e
}
