package payment

import (
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	def "github.com/xgmsx/rsf/order/internal/client"
	"github.com/xgmsx/rsf/order/internal/model"
	genPaymentV1 "github.com/xgmsx/rsf/shared/pkg/proto/payment/v1"
)

var _ def.PaymentClient = (*client)(nil)

type client struct {
	generatedClient genPaymentV1.PaymentServiceClient
}

func NewClient(addr string) *client {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to %s: %v\n", addr, err)
	}

	return &client{
		generatedClient: genPaymentV1.NewPaymentServiceClient(conn),
	}
}

func (c *client) PayOrder(ctx context.Context, userUUID, orderUUID uuid.UUID, paymentMethod model.PaymentMethod) (*uuid.UUID, error) {
	paymentMethodsMap := map[model.PaymentMethod]genPaymentV1.PaymentMethod{
		model.PaymentMethodCARD:          genPaymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		model.PaymentMethodSBP:           genPaymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
		model.PaymentMethodCREDITCARD:    genPaymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		model.PaymentMethodINVESTORMONEY: genPaymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
	}
	genPaymentMethod, ok := paymentMethodsMap[paymentMethod]
	if !ok {
		return nil, model.ErrPaymentMethodIsNotSupported
	}

	res, err := c.generatedClient.PayOrder(ctx, &genPaymentV1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		UserUuid:      userUUID.String(),
		PaymentMethod: genPaymentMethod,
	})
	if err != nil {
		return nil, err
	}

	txUUID, err := uuid.Parse(res.TransactionUuid)
	if err != nil {
		return nil, err
	}

	return &txUUID, nil
}
