package main

import (
	"context"
	"fmt"
	"github.com/INEFFABLE-games/PriceService/models"
	"github.com/INEFFABLE-games/PriceService/protocol"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	clients "priceClient/internal/client"
	"priceClient/internal/config"
	"priceClient/internal/service"
)

func main() {

	cfg := config.NewConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// dial server
	conn, err := grpc.Dial(fmt.Sprintf(":%s", cfg.GrpcPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	// create stream
	client := protocol.NewPriceServiceClient(conn)

	stream, err := client.Send(ctx)
	if err != nil {
		log.Fatalf("unable to open stream %v", err)
	}

	priceChannel := make(chan []models.Price)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	priceClient := clients.NewPricesClient(ctx, priceChannel)

	//start grpc stream and write data into price channel
	go func() {
		err := priceClient.Send(stream)
		if err != nil {
			log.Errorf("unable to call send stream %v", err.Error())
		}

	}()

	//start listening prices channel
	go func() {
		priceService := service.NewPriceService(priceChannel)
		priceService.StartListening(ctx)
	}()

	<-c
	cancel()
	os.Exit(1)
}
