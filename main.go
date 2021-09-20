package main

import (
	"context"
	"encoding/json"
	"fmt"
	protocol2 "github.com/INEFFABLE-games/PositionService/protocol"
	"github.com/INEFFABLE-games/PriceService/models"
	"github.com/INEFFABLE-games/PriceService/protocol"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io"
	"os"
	"os/signal"
	"priceClient/internal/client"
	"priceClient/internal/config"
	"priceClient/internal/menu"
	"time"
)

var currentPrices []models.Price

func main() {

	cfg := config.NewConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// dial server
	conn, err := grpc.Dial(fmt.Sprintf(":%s", cfg.GrpcPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	conn2, err := grpc.Dial(fmt.Sprintf(":%s", cfg.PositionsPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	// create stream
	grpcClient := protocol.NewPriceServiceClient(conn)
	grpcClient2 := protocol2.NewPositionServiceClient(conn2)

	stream, err := grpcClient.Send(ctx)
	if err != nil {
		log.Fatalf("unable to open stream %v", err)
	}
	posClient := client.NewPositionClient(grpcClient2)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	//priceClient := clients.NewPricesClient(ctx, priceChannel)

	//start grpc stream and write data into price channel
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:

				data, err := stream.Recv()
				if err == io.EOF {
					log.WithFields(log.Fields{
						"handler": "main",
						"action":  "get data from stream",
					}).Errorf("End of file %v", err.Error())
				}

				butchOfPrices := []models.Price{}

				err = json.Unmarshal(data.ButchOfPrices, &butchOfPrices)
				if err != nil {
					log.WithFields(log.Fields{
						"handler": "main",
						"action":  "unmarshal butch of prices",
					}).Errorf("unable to unmarshal butch of prices %v", err.Error())
				}

				go func() {
					currentPrices = butchOfPrices
				}()
			}
		}
	}()

	//start listening prices channel
	go func() {
		/*priceService := service.NewPriceService(priceChannel)
		priceService.StartListening(ctx,&currentPrices)*/
	}()

	go func() {
		menu.Start(ctx, &currentPrices, *posClient)
	}()

	<-c
	cancel()
	os.Exit(1)
}
