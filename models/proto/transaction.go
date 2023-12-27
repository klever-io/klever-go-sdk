package proto

import "context"

func (x *Transaction) AddSignature(s []byte) {
	if x.Signature == nil {
		x.Signature = make([][]byte, 0)
	}

	x.Signature = append(x.Signature, s)
}

type Signer interface {
	Sign([]byte) ([]byte, error)
}

type Broadcaster interface {
	BroadcastTransaction(context.Context, *Transaction) (string, error)
}

func (x *Transaction) Sign(signer Signer) error {
	signature, err := signer.Sign(x.Hash)
	if err != nil {
		return err
	}

	x.AddSignature(signature)

	return nil
}

func (x *Transaction) Broadcast(provider Broadcaster) (string, error) {
	return provider.BroadcastTransaction(context.Background(), x)
}
