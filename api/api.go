package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/steinfletcher/payments"
	"github.com/xeipuuv/gojsonschema"
)

type Server struct {
	Router  *gin.Engine
	service acme.PaymentService
	server  *http.Server
}

var jsonSchemaValidator = gojsonschema.NewStringLoader(acme.AttributesSchema)

// NewServer creates a new server with all application routes defined
// The caller must call `Start` to bind to the network and start serving requests
func NewServer(service acme.PaymentService) *Server {
	r := gin.Default()

	srv := &Server{Router: r, service: service}
	r.GET("/health", srv.healthCheck)

	v1 := r.Group("/v1")
	v1.Use(errorHandler)

	v1.GET("/payment", srv.getAllPayments)
	v1.GET("/payment/:id", srv.getPayment)
	v1.POST("/payment", srv.createPayment)
	v1.PUT("/payment/:id", srv.updatePayment)
	v1.DELETE("/payment/:id", srv.deletePayment)

	return srv
}

func (r *Server) createPayment(ctx *gin.Context) {
	payment := acme.Payment{}
	err := ctx.Bind(&payment)
	if err != nil {
		ctx.Error(acme.InvalidRequestBody)
		return
	}

	err = validatePayment(payment)
	if err != nil {
		ctx.Error(err)
		return
	}

	id, err := r.service.Create(payment)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Header("Location", id.String())
	ctx.AbortWithStatus(http.StatusCreated)
}

func validatePayment(p acme.Payment) error {
	if p.OrganisationID == uuid.Nil {
		err := acme.InvalidField
		err.Detail = "organisation Id must be provided"
		return err
	}
	result, err := gojsonschema.Validate(jsonSchemaValidator, gojsonschema.NewGoLoader(p.Attributes))
	if err != nil {
		return acme.ServerError
	}
	if !result.Valid() {
		err := acme.InvalidField
		err.Detail = fmt.Sprintf("invalid attributes: %s", result.Errors())
		return err
	}
	return nil
}

func (r *Server) getPayment(ctx *gin.Context) {
	id := ctx.Param("id")
	externalID, err := uuid.Parse(id)
	if err != nil {
		ctx.Error(acme.InvalidID)
		return
	}

	payment, err := r.service.Get(externalID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, payment)
}

func (r *Server) getAllPayments(ctx *gin.Context) {
	payments, err := r.service.GetAll()
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, payments)
}

func (r *Server) updatePayment(ctx *gin.Context) {
	id := ctx.Param("id")
	externalID, err := uuid.Parse(id)
	if err != nil {
		ctx.Error(acme.InvalidID)
		return
	}

	payment := acme.Payment{}
	err = ctx.Bind(&payment)
	if err != nil {
		ctx.Error(acme.InvalidField)
		return
	}

	err = r.service.Update(externalID, payment)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
}

func (r *Server) deletePayment(ctx *gin.Context) {
	id := ctx.Param("id")
	externalID, paramErr := uuid.Parse(id)
	if paramErr != nil {
		ctx.Error(acme.InvalidID)
		return
	}

	err := r.service.Delete(externalID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
}

func (r *Server) healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func (r *Server) Start(port string) {
	r.server = &http.Server{
		Addr:    ":" + port,
		Handler: r.Router,
	}
	defer r.Close()

	go func() {
		if err := r.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}

func (r *Server) Close() {
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := r.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown failed: %s.", err)
	}

	log.Println("Server shutdown complete.")
}
