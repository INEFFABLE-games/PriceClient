package service

import (
	"context"
	"github.com/INEFFABLE-games/PriceService/models"
	log "github.com/sirupsen/logrus"
	"time"
)

type PriceService struct {
	pricesChannel chan []models.Price
}

func (p *PriceService) StartListening(ctx context.Context) {

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:

			butchOfPrices := <-p.pricesChannel

			for _, v := range butchOfPrices {
				log.Info(v)
			}

		}
	}

}

func NewPriceService(pricesChannel chan []models.Price) *PriceService {
	return &PriceService{pricesChannel: pricesChannel}
}
