package provider

import (
	"github.com/klever-io/klever-go-sdk/models"
	"github.com/klever-io/klever-go-sdk/models/proto"
)

func (kc *kleverChain) Proposal(base *models.BaseTX, description string, parameters map[int32]string, duration uint32) (*proto.Transaction, error) {
	contracts := []interface{}{models.ProposalTXRequest{
		Parameters:     parameters,
		EpochsDuration: duration,
		Description:    description,
	}}

	data, err := kc.buildRequest(proto.TXContract_ProposalContractType, base, contracts)
	if err != nil {
		return nil, err
	}

	return kc.PrepareTransaction(data)
}

func (kc *kleverChain) Vote(base *models.BaseTX, proposalID uint64, amount float64, voteType uint64) (*proto.Transaction, error) {
	contracts := []interface{}{models.VoteTXRequest{
		Type:       uint32(voteType),
		ProposalID: proposalID,
		Amount:     int64(amount * 1000000),
	}}

	data, err := kc.buildRequest(proto.TXContract_VoteContractType, base, contracts)
	if err != nil {
		return nil, err
	}

	return kc.PrepareTransaction(data)
}
