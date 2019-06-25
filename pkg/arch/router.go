package arch

import (
	"fmt"

	"golang.org/x/tools/go/types/typeutil"
)

// Router is responsible for routing action requests to the correct service
type Router struct {
	actionToServiceMap   typeutil.Map
	actionRequestChannel chan ActionRequest
}

// NewRouter creates a router
func NewRouter() *Router {
	channel := make(chan ActionRequest, 1000)
	return &Router{typeutil.Map{}, channel}
}

// Register a service to the router so it can route action requests
func (r *Router) Register(service Service) error {
	// add the service to the action map
	actionType := service.GetActionRequestType()
	if actionType == nil {
		return nil
	}
	if r.actionToServiceMap.At(actionType) != nil {
		return fmt.Errorf("Service already fulfilling ActionRequest of type: %v", actionType.String())
	}

	r.actionToServiceMap.Set(actionType, service)
	return nil
}
