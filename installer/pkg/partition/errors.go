package partition

type NewPartitionSizeError struct {
	Err string
}

func (e NewPartitionSizeError) Error() string {
	return e.Err
}

type NewPartitionError struct {
	Err string
}

func (e NewPartitionError) Error() string {
	return e.Err
}

type CreatePartitioningFileError struct {
	Message string
	Err     error
}

func (e *CreatePartitioningFileError) Error() string {
	return e.Message
}

func (e *CreatePartitioningFileError) Unwrap() error {
	return e.Err
}
