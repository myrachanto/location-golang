package repository

import (
	"fmt"
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
		fmt.Println("-------------------------------------", l.Orders)
		location.Thetime = time.Now()
		locs := []model.Location{}
		for j, g := range l.Orders {
			if g.OrderId == order_id {
				locs = append(locs, *location)
			}
			g.Location = append(g.Location, locs...)
			l.Orders[j] = g
		}
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

// this functions queries the order with the maximum number of location histories as per the query
func (l *locationrepo) GetMAxOrders(order_id string, max int) (order *model.Order, err httperors.HttpErr) {
	log.Println("testing --------------max", max)
	orderring := model.Order{}
	if max == 0 {
		for _, o := range l.Orders {
			log.Println("testing --------------", o)
			l.sorttime(o.Location)
			orderring = o
		}
		return &orderring, nil
	}
	for _, o := range l.Orders {
		if o.OrderId == order_id {
			if max > len(o.Location) {
				orderring.OrderId = order_id
				l.sorttime(o.Location)
				orderring.Location = o.Location
			} else {
				orderring.OrderId = order_id
				l.sorttime(o.Location)
				orderring.Location = o.Location[:max]
			}
		}
	}
	return &orderring, nil
}

// sort time
func (l *locationrepo) sorttime(dat []model.Location) {
	for i := 0; i < len(dat); i++ {
		for j := i + 1; j < len(dat); j++ {
			if dat[i].Thetime.Unix() < dat[j].Thetime.Unix() {
				dat[i], dat[j] = dat[j], dat[i]
			}
		}
	}
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
