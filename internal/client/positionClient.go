package client

import (
	"context"
	"github.com/INEFFABLE-games/PositionService/protocol"
	"github.com/INEFFABLE-games/PriceService/models"
	"google.golang.org/grpc"
)

type PositionClient struct {
	client          protocol.PositionServiceClient
	ClientPositions map[string]models.Price
}

func (p *PositionClient) Buy(ctx context.Context, in *protocol.BuyRequest, opts ...grpc.CallOption) (*protocol.BuyReply, error) {
	reply, err := p.client.Buy(ctx, &protocol.BuyRequest{
		UserToken: in.UserToken,
		Price:     in.GetPrice(),
	})
	return reply, err
}

func (p *PositionClient) Sell(ctx context.Context, in *protocol.SellRequest, opts ...grpc.CallOption) (*protocol.SellReply, error) {
	reply, err := p.client.Sell(ctx, &protocol.SellRequest{
		UserToken: in.UserToken,
		Price:     in.Price,
	})
	return reply, err
}

func (p *PositionClient) Get(ctx context.Context, in *protocol.GetRequest, opts ...grpc.CallOption) (*protocol.GetReply, error) {
	reply, err := p.client.Get(ctx, &protocol.GetRequest{
		UserToken: in.UserToken,
	})
	return reply, err
}

func NewPositionClient(client protocol.PositionServiceClient) *PositionClient {
	return &PositionClient{client: client}
}
