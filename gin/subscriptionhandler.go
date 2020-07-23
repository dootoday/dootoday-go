package service

import (
	"apidootoday/config"
	orderservice "apidootoday/orderservice"
	subscriptionservice "apidootoday/subscription"
	userservice "apidootoday/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

// SubscriptionHandler :
type SubscriptionHandler struct {
	SubscriptionService *subscriptionservice.SubscriptionService
	OrderService        *orderservice.OrderService
	UserService         *userservice.UserService
}

// NewSubscriptionHandler :
func NewSubscriptionHandler(
	ss *subscriptionservice.SubscriptionService,
	os *orderservice.OrderService,
	us *userservice.UserService,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		SubscriptionService: ss,
		OrderService:        os,
		UserService:         us,
	}
}

// GetPlans :
func (sh *SubscriptionHandler) GetPlans(c *gin.Context) {
	type RequestBody struct {
		PromoCode string `form:"code"` // optional
	}
	var request RequestBody
	err := c.Bind(&request)
	if err != nil {
		glog.Error("Bad input")
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "code is missing"},
		)
		return
	}

	// this was set in context from auth middleware
	userID, ok := c.Get("user_id")

	if !ok {
		glog.Error("Could not get the user id from context")
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "could not get the user id from context"},
		)
		return
	}

	plans, err := sh.SubscriptionService.GetPlansToDisplay(
		userID.(uint), request.PromoCode,
	)
	// if err != nil {
	// 	glog.Error(err)
	// 	c.JSON(
	// 		http.StatusBadRequest,
	// 		gin.H{"error": err.Error()},
	// 	)
	// 	return
	// }
	type ResponseBody struct {
		ID                 uint   `json:"plan_id"`
		Name               string `json:"name"`
		Desc               string `json:"description"`
		AmountInCents      int    `json:"amount"`
		OfferAmountInCents int    `json:"offer_amount"`
	}
	resp := []ResponseBody{}
	for _, plan := range plans {
		resp = append(
			resp,
			ResponseBody{
				ID:                 plan.ID,
				Name:               plan.Name,
				Desc:               plan.Description,
				AmountInCents:      plan.AmountInCents,
				OfferAmountInCents: plan.OfferAmountInCents,
			},
		)
	}
	c.JSON(http.StatusOK, resp)
	return
}

// Subscribe :
func (sh *SubscriptionHandler) Subscribe(c *gin.Context) {
	uID, ok := c.Get("user_id")
	userID := uID.(uint)
	if !ok {
		glog.Error("Could not get the user id from context")
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "could not get the user id from context"},
		)
		return
	}
	pID := c.Param("plan_id")
	planID, err := strconv.ParseUint(pID, 10, 32)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid plan ID"},
		)
		return
	}
	plan, err := sh.SubscriptionService.GetPlanByID(uint(planID))
	if err != nil {
		glog.Error("Plan not found ", err)
		c.JSON(
			http.StatusNotFound,
			gin.H{"error": "Plan not found"},
		)
		return
	}
	// This is to verify if the user can avail the plan
	// third param true is very important
	err = sh.SubscriptionService.CreateSubscripton(userID, plan.ID, true)
	if err != nil {
		glog.Error("User can not use the plan ", err)
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": "User can not use the plan"},
		)
		return
	}
	type TaskResponse struct {
		KeyID       string `json:"key_id"`
		OrderID     string `json:"order_id"`
		Name        string `json:"name"`        // Company name
		Description string `json:"description"` // Company description
		Image       string `json:"image"`
		UserName    string `json:"user_full_name"`
		UserEmail   string `json:"user_email"`
		UserPhone   string `json:"user_phone"`
		CallBackURL string `json:"callback_url"`
		CancelURL   string `json:"cancel_url"`
		Amount      int    `json:"amount"`
	}
	resp := TaskResponse{}
	// else create a new order
	orderID, err := sh.OrderService.CreateNewOrder(userID, plan.ID, plan.OfferAmountInCents)
	if err != nil {
		glog.Error("Could not create a new order", err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Could not create a new order"},
		)
		return
	}
	user, err := sh.UserService.GetUserByID(userID)
	if err != nil {
		glog.Error("Could not fetch the user", err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Could not fetch the user"},
		)
		return
	}
	resp.KeyID = config.RPApiKey
	resp.OrderID = orderID
	resp.Name = config.DooTodayName
	resp.Description = config.DooTodayDesc
	resp.Image = config.DooTodayLogo
	resp.UserName = user.FirstName + " " + user.LastName
	resp.UserEmail = user.Email
	resp.UserPhone = "9066258469"
	resp.CallBackURL = config.BackendBase + "/v1/payment-success"
	resp.CancelURL = config.FrontendBase + "/me/subscription?cs=false"
	resp.Amount = plan.OfferAmountInCents
	c.JSON(http.StatusOK, resp)
	return
}

// PaymentSuccess :
func (sh *SubscriptionHandler) PaymentSuccess(c *gin.Context) {
	// c.Request.ParseMultipartForm(1000)
	// for key, value := range c.Request.PostForm {
	// 	fmt.Println(key, value)
	// }
	type RequestBody struct {
		RPOrderID   string `form:"razorpay_order_id" binding:"required"`
		RPPaymentID string `form:"razorpay_payment_id" binding:"required"`
		RPSignature string `form:"razorpay_signature" binding:"required"`
	}
	var req RequestBody
	err := c.Bind(&req)
	if err != nil {
		glog.Error("Bad input")
		c.Data(
			http.StatusOK,
			"text/html; charset=utf-8",
			[]byte(
				"<script> window.location.replace('"+
					config.FrontendBase+
					"/me/subscription?cs=false'); </script>",
			),
		)
		return
	}
	order, err := sh.OrderService.GetOrderByRPOrderID(req.RPOrderID)
	if err != nil {
		glog.Error("Can't find the order")
		c.Data(
			http.StatusOK,
			"text/html; charset=utf-8",
			[]byte(
				"<script> window.location.replace('"+
					config.FrontendBase+
					"/me/subscription?cs=false'); </script>",
			),
		)
		return
	}
	order.RPPaymentID = req.RPPaymentID
	order.RPSignature = req.RPSignature
	err = sh.OrderService.UpdateOrder(order)
	if err != nil {
		glog.Error("Can't find the order")
		c.Data(
			http.StatusOK,
			"text/html; charset=utf-8",
			[]byte(
				"<script> window.location.replace('"+
					config.FrontendBase+
					"/me/subscription?cs=false'); </script>",
			),
		)
		return
	}
	err = sh.SubscriptionService.CreateSubscripton(order.UserID, order.PlanID, false)
	if err != nil {
		glog.Error("Could not subscribe to the plan ", err)
		c.Data(
			http.StatusOK,
			"text/html; charset=utf-8",
			[]byte(
				"<script> window.location.replace('"+
					config.FrontendBase+
					"/me/subscription?cs=false'); </script>",
			),
		)
		return
	}

	c.Data(
		http.StatusOK,
		"text/html; charset=utf-8",
		[]byte(
			"<script> window.location.replace('"+
				config.FrontendBase+
				"/me/subscription?cs=true'); </script>",
		),
	)
	return
}
