package service

import (
	orderservice "apidootoday/orderservice"
	subscriptionservice "apidootoday/subscription"
	userservice "apidootoday/user"
)

// OrderHandler :
type OrderHandler struct {
	UserService        *userservice.UserService
	SubscriptionSrvice *subscriptionservice.SubscriptionService
	OrderService       *orderservice.OrderService
}
