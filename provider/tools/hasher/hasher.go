package hasher

// NewHasher will return a new instance of hasher based on the value stored in config
func NewHasher() (Hasher, error) {
	return &Blake2b{}, nil
}
