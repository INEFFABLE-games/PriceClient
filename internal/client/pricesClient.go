package client

import (
	"context"
	"encoding/json"
	"github.com/INEFFABLE-games/PriceService/models"
	"github.com/INEFFABLE-games/PriceService/protocol"
	log "github.com/sirupsen/logrus"
	"io"
	"time"
)

type PricesClient struct {
	priceChannel chan []models.Price
	ctx          context.Context
}

func (p *PricesClient) Send(stream protocol.PriceService_SendClient) error {

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-p.ctx.Done():
			return nil
		case <-ticker.C:

			data, err := stream.Recv()
			if err == io.EOF {
				log.WithFields(log.Fields{
					"handler": "client",
					"action":  "get data from stream",
				}).Errorf("End of file %v", err.Error())
			}

			butchOfPrices := []models.Price{}

			err = json.Unmarshal(data.ButchOfPrices, &butchOfPrices)
			if err != nil {
				log.WithFields(log.Fields{
					"handler": "pricesClient",
					"action":  "unmarshal butch of prices",
				}).Errorf("unable to unmarshal butch of prices %v", err.Error())
			}

			go func() {
				p.priceChannel <- butchOfPrices
			}()

		}
	}
}

func NewPricesClient(ctx context.Context, c chan []models.Price) *PricesClient {
	return &PricesClient{
		ctx:          ctx,
		priceChannel: c}
}
