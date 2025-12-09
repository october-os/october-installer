package partition

type NewPartitionSizeError struct {
	Err string
}

func (e NewPartitionSizeError) Error() string {
	return e.Err
}
