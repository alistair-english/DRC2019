package arch

import (
	"fmt"
	"reflect"
)

// Router is responsible for routing action requests to the correct service
type Router struct {
	actionToServiceMap   map[reflect.Type]Service
	actionRequestChannel chan ActionRequest
}

// NewRouter creates a router
func NewRouter() *Router {
	return &Router{make(map[reflect.Type]Service), make(chan ActionRequest, 1000)}
}

// Register a service to the router so it can route action requests
func (r *Router) Register(service Service) error {
	// add the service to the action map
	actionType := service.GetActionRequestType()

	if actionType != nil {
		if _, exists := r.actionToServiceMap[actionType]; exists {
			return fmt.Errorf("Service already fulfilling ActionRequest of type: %v", actionType.String())
		}
		r.actionToServiceMap[actionType] = service
	}

	service.SetActionRequestChannel(r.actionRequestChannel)
	return nil
}

func (r *Router) Start() {
	for request := range r.actionRequestChannel {
		t := reflect.TypeOf(request)
		service := r.actionToServiceMap[t]
		if service != nil {
			go service.FulfullActionRequest(request)
		} else {
			fmt.Println("No service for request type: ", t.Name())
		}
	}
}
