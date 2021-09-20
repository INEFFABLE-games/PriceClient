package menu

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/INEFFABLE-games/PositionService/protocol"
	"github.com/INEFFABLE-games/PriceService/models"
	"github.com/dixonwille/wmenu/v5"
	"github.com/inancgumus/screen"
	log "github.com/sirupsen/logrus"
	"priceClient/internal/client"
	"time"
)

func Start(ctx context.Context, currentPrices *[]models.Price, client client.PositionClient) {
	ticker := time.NewTicker(100 * time.Millisecond)
	client.ClientPositions = make(map[string]models.Price)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			mainMenu := wmenu.NewMenu("Select the option:")
			mainMenu.Option("Buy position", nil, false, func(opt wmenu.Opt) error {
				selectPrices := wmenu.NewMenu("Select price:")

				for _, v := range *currentPrices {
					selectPrices.Option(fmt.Sprintf("name: %v,bid: %v,ask: %v", v.Name, v.Bid, v.Ask), nil, false, func(opt wmenu.Opt) error {
						screen.Clear()

						userToken := "test"
						price := (*currentPrices)[opt.ID]

						marshalPrice, err := json.Marshal(price)
						if err != nil {
							log.WithFields(log.Fields{
								"handler ": "menu",
								"action ":  "marshal price",
							}).Errorf("unable to marshal price %v", err.Error())
						}

						res, err := client.Buy(ctx, &protocol.BuyRequest{
							UserToken: &userToken,
							Price:     marshalPrice,
						})
						if err != nil {
							log.WithFields(log.Fields{
								"handler ": "priceMenu",
								"action ":  "buy position",
							}).Errorf("unable to buy position %v", err.Error())
						}

						// add new price to local client map
						price.Ask = 0
						client.ClientPositions[price.Name] = price

						log.Infof("sending %v", price)

						log.Info(res.GetMessage())

						return nil
					})
				}

				err := selectPrices.Run()
				if err != nil {
					log.WithFields(log.Fields{
						"handler ": "menu",
						"action ":  "run selectPrices menu from Buy position",
					}).Errorf("Unnable to run menu %v", err.Error())
				}

				return nil
			})
			mainMenu.Option("Sell position", nil, false, func(opt wmenu.Opt) error {
				selectPrices := wmenu.NewMenu("Select price:")

				for _, v := range *currentPrices {
					selectPrices.Option(fmt.Sprintf("name: %v,bid: %v,ask: %v", v.Name, v.Bid, v.Ask), nil, false, func(opt wmenu.Opt) error {
						screen.Clear()

						userToken := "test"
						price := (*currentPrices)[opt.ID]

						marshalPrice, err := json.Marshal(price)
						if err != nil {
							log.WithFields(log.Fields{
								"handler ": "menu",
								"action ":  "marshal price",
							}).Errorf("unable to marshal price %v", err.Error())
						}

						res, err := client.Sell(ctx, &protocol.SellRequest{
							UserToken: &userToken,
							Price:     marshalPrice,
						})
						if err != nil {
							log.WithFields(log.Fields{
								"handler ": "priceMenu",
								"action ":  "sell position",
							}).Errorf("unable to sell position %v", err.Error())
						}

						delete(client.ClientPositions, price.Name)

						log.Info(res.GetMessage())

						return nil
					})
				}

				err := selectPrices.Run()
				if err != nil {
					log.WithFields(log.Fields{
						"handler ": "menu",
						"action ":  "run selectPrices menu from Buy position",
					}).Errorf("Unnable to run menu %v", err.Error())
				}

				return nil
			})
			mainMenu.Option("Show your positions", nil, false, func(opt wmenu.Opt) error {
				for _, v := range client.ClientPositions {

					if v.Ask != 0 {
						log.Infof("%v,CLOSED", v)
						continue
					}
					log.Info(v)
				}
				return nil
			})
			mainMenu.Option("Show current prices", nil, false, func(opt wmenu.Opt) error {
				screen.Clear()

				for _, v := range *currentPrices {
					log.Info(v)
				}

				back := wmenu.NewMenu("")
				back.Option("Back", nil, false, func(opt wmenu.Opt) error {
					screen.Clear()
					return nil
				})

				err := back.Run()
				if err != nil {
					log.WithFields(log.Fields{
						"handler ": "menu",
						"action ":  "run back menu from Show current prices",
					}).Errorf("Unnable to run menu %v", err.Error())
				}

				return nil
			})

			err := mainMenu.Run()
			if err != nil {
				log.WithFields(log.Fields{
					"handler ": "menu",
					"action ":  "run menu",
				}).Errorf("unable to run menu %v", err.Error())
			}
		}
	}
}
