package service

import (
	httperors "github.com/myrachanto/custom-http-error"
	"github.com/myrachanto/location/src/cmd/model"
	r "github.com/myrachanto/location/src/cmd/repository"
)

//locationService ...
var (
	LocationService LocationServiceInterface = &locationService{}
)

type LocationServiceInterface interface {
	OrderLocation(order_id string, location *model.Location) (string, httperors.HttpErr)
	GetMAxOrders(order_id string, max int) (*model.Order, httperors.HttpErr)
	DeleteOrder(order_id string) (string, httperors.HttpErr)
}

type locationService struct {
	repository r.LocationRepoInterface
}

func NewlocationService(repo r.LocationRepoInterface) LocationServiceInterface {
	return &locationService{
		repo,
	}
}

func (service locationService) OrderLocation(order_id string, location *model.Location) (string, httperors.HttpErr) {
	s, err1 := r.Locationrepo.OrderLocation(order_id, location)
	if err1 != nil {
		return "", err1
	}
	return s, nil

}
func (service locationService) GetMAxOrders(order_id string, max int) (*model.Order, httperors.HttpErr) {
	results, err := r.Locationrepo.GetMAxOrders(order_id, max)
	return results, err
}
func (service locationService) DeleteOrder(order_id string) (string, httperors.HttpErr) {
	success, failure := r.Locationrepo.DeleteOrder(order_id)
	return success, failure
}
