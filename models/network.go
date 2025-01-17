package models

type NetworkConfig struct {
	NumMetachainNodes  uint64 `json:"klv_num_metachain_nodes"`
	ConsensusGroupSize uint64 `json:"klv_consensus_group_size"`
	ChainID            string `json:"klv_chain_id"`
	SlotInterval       uint64 `json:"klv_slot_time"`
	SlotsPerEpoch      uint64 `json:"klv_slots_per_epoch"`
	StartTime          uint64 `json:"klv_start_time"`
}
