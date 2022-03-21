package controllers

import (
	"strconv"

	"net/http"

	"github.com/labstack/echo/v4"
	httperors "github.com/myrachanto/custom-http-error"
	"github.com/myrachanto/location/src/cmd/model"
	service "github.com/myrachanto/location/src/cmd/services"
)

//locationController ..
var (
	LocationController LocationcontrollerInterface = &locationController{}
)

type locationController struct {
	service service.LocationServiceInterface
}
type LocationcontrollerInterface interface {
	OrderLocation(c echo.Context) error
	GetMAxOrders(c echo.Context) error
	DeleteOrder(c echo.Context) error
}

func NewlocationController(ser service.LocationServiceInterface) LocationcontrollerInterface {
	return &locationController{
		ser,
	}
}

/////////controllers/////////////////
func (controller locationController) OrderLocation(c echo.Context) error {
	// log.Println("----------------------------------- mwassssss")
	location := &model.Location{}
	if err := c.Bind(location); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code(), httperror)
	}
	order_id := c.Param("order_id")
	// log.Println("-----------------------------------", order_id)
	s, err1 := service.LocationService.OrderLocation(order_id, location)
	if err1 != nil {
		return c.JSON(err1.Code(), err1)
	}
	return c.JSON(http.StatusCreated, s)
}
func (controller locationController) GetMAxOrders(c echo.Context) error {
	order_id := c.Param("order_id")
	max, err := strconv.Atoi(c.QueryParam("max"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid Max number")
		return c.JSON(httperror.Code(), httperror)
	}
	// log.Println("----------------------------------- mwassssss", order_id, max)
	results, err3 := service.LocationService.GetMAxOrders(order_id, max)
	if err3 != nil {
		return c.JSON(err3.Code(), err3)
	}
	return c.JSON(http.StatusOK, results)
}
func (controller locationController) DeleteOrder(c echo.Context) error {
	order_id := c.Param("order_id")
	success, failure := service.LocationService.DeleteOrder(order_id)
	if failure != nil {
		return c.JSON(failure.Code(), failure)
	}
	return c.JSON(http.StatusOK, success)

}
