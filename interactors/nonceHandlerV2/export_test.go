package nonceHandlerV2

import (
	"github.com/klever-io/klever-go-sdk/core/address"
	"github.com/klever-io/klever-go-sdk/interactors"
	"github.com/klever-io/klever-go/data/transaction"
	"github.com/klever-io/klever-go/tools/check"
)

// NewAddressNonceHandlerWithPrivateAccess -
func NewAddressNonceHandlerWithPrivateAccess(proxy interactors.Proxy, address address.Address) (*addressNonceHandler, error) {
	if check.IfNil(proxy) {
		return nil, interactors.ErrNilProxy
	}
	if check.IfNil(address) {
		return nil, interactors.ErrNilAddress
	}
	return &addressNonceHandler{
		address:      address,
		proxy:        proxy,
		transactions: make(map[uint64]*transaction.Transaction),
	}, nil
}
