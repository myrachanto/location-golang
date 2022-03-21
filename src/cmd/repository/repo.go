package repository

import (
	"log"
	"strconv"
	"time"

	httperors "github.com/myrachanto/custom-http-error"
	"github.com/myrachanto/location/src/cmd/model"
	"github.com/spf13/viper"
)

//locationrepo ...
var (
	Locationrepo LocationRepoInterface = &locationrepo{}
)

type Key struct {
	EncryptionKey string `mapstructure:"EncryptionKey"`
}

type locationrepo struct {
	Orders []model.Order
}

type LocationRepoInterface interface {
	OrderLocation(order_id string, location *model.Location) (string, httperors.HttpErr)
	GetMAxOrders(order_id string, max int) (*model.Order, httperors.HttpErr)
	orderExistById(order_id string) *model.Order
	DeleteOrder(order_id string) (string, httperors.HttpErr)
	addtoLocationsToOrder(location *model.Location)
	StartClearHistory()
	// removeorder(key int)
}

func NewlocationRepo() *locationrepo {
	return &locationrepo{}
}

type Duraion struct {
	Duration string `mapstructure:"LOCATION_HISTORY_TTL_SECONDS"`
}

// func LoadConfig(path string) (open Open, err error) {
func LoadConfig() (open Duraion, err error) {
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

//this function creates an order
func (l *locationrepo) OrderLocation(order_id string, location *model.Location) (string, httperors.HttpErr) {
	if err1 := location.Validate(); err1 != nil {
		return "", err1
	}
	//check if the orderExist in memory
	orderExist := l.orderExistById(order_id)
	if orderExist != nil {
		//add to the locations
		location.Thetime = time.Now()
		l.addtoLocationsToOrder(location)
		return "location added successifully", nil
	}
	order := &model.Order{}
	order.OrderId = order_id
	location.Thetime = time.Now()
	order.Location = append(order.Location, *location)
	l.Orders = append(l.Orders, *order)
	return "location created successifully", nil
}

//check if the order exist
func (l *locationrepo) orderExistById(order_id string) *model.Order {
	for _, o := range l.Orders {
		if o.OrderId == order_id {
			return &o
			break
		}
	}
	return nil

}

//this is a private function adds to the in-memory database
func (l *locationrepo) addtoLocationsToOrder(location *model.Location) {
	for _, o := range l.Orders {
		if o.OrderId == o.OrderId {
			o.Location = append(o.Location, *location)
			break
		}
	}

}

// this functions queries the order with the maximum number of location histories as per the query
func (l *locationrepo) GetMAxOrders(order_id string, max int) (order *model.Order, err httperors.HttpErr) {
	log.Println("testing --------------", l.Orders)
	orderring := model.Order{}
	for _, o := range l.Orders {
		if o.OrderId == order_id {
			orderring.OrderId = order_id
			locations := []model.Location{}
			for k, v := range o.Location {
				if k <= max {
					locations = append(locations, v)
					break
				}
			}
			orderring.Location = locations
			break
		}
	}
	return &orderring, nil
}
func (l *locationrepo) DeleteOrder(order_id string) (string, httperors.HttpErr) {
	orders := []model.Order{}
	for _, o := range l.Orders {
		if o.OrderId != order_id {
			orders = append(orders, o)
		}
	}
	l.Orders = orders
	return "Successifully deleted", nil
}

// this functions calls itself as recursive function to run after the seconds recorded in the env file
func (l *locationrepo) StartClearHistory() {
	times, _ := LoadConfig()
	timer, _ := strconv.Atoi(times.Duration)
	l.ClearHistory()
	_ = time.AfterFunc(time.Second*time.Duration(timer), l.StartClearHistory)
}

//this is the self cleannup function
func (l *locationrepo) ClearHistory() {
	times, err := LoadConfig()
	if err != nil {
		return
	}
	timer, err := strconv.Atoi(times.Duration)
	if err != nil {
		return
	}
	locations := []model.Location{}
	ordr := &model.Order{}
	for k, order := range l.Orders {
		for _, location := range order.Location {
			// fmt.Println(">>>>>>>>>>>>>>>>>", time.Since(location.Thetime), "---------", time.Duration(timer)*time.Second)
			if time.Since(location.Thetime) < time.Duration(timer)*time.Second {
				// fmt.Println("---------------sfrdddd")
				locations = append(locations, location)
			}
		}
		// fmt.Println("aaaaaaaaaaaaaa", locations)
		ordr.Location = locations
		ordr.OrderId = order.OrderId
		l.Orders[k] = *ordr

		// fmt.Println("ffffffffffffffff", order)

	}
}
