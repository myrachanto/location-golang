package routes

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/myrachanto/location/src/cmd/controllers"
	"github.com/myrachanto/location/src/cmd/repository"
	service "github.com/myrachanto/location/src/cmd/services"
	"github.com/spf13/viper"
)

func init() {
	log.SetPrefix("Location Api: ")
}

//StoreAPI =>entry point to routes
type Open struct {
	Port string `mapstructure:"HISTORY_SERVER_LISTEN_ADDR"`
}

// func LoadConfig(path string) (open Open, err error) {
func LoadConfig() (open Open, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&open)
	return
}
func LocationApi() {
	//check db connection//////////////////////
	fmt.Println("initialization----------------")
	controllers.NewlocationController(service.NewlocationService(repository.NewlocationRepo()))
	e := echo.New()
	repository.Locationrepo.StartClearHistory()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.POST("/location/:order_id/now", controllers.LocationController.OrderLocation)
	e.GET("/location/:order_id", controllers.LocationController.GetMAxOrders)
	e.DELETE("/location/:order_id", controllers.LocationController.DeleteOrder)

	open, err := LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// Start server
	if open.Port != "" {
		e.Logger.Fatal(e.Start(open.Port))
	} else {
		e.Logger.Fatal(e.Start(":8081"))
	}
}
