package container

// RewindableGenerator is a struct that holds a generator function and a count.
type RewindableGenerator struct {
	generator func() []interface{}
	count     int
}

// NewRewindableGenerator creates a new RewindableGenerator with the provided generator function and count.
func NewRewindableGenerator(generator func() []interface{}, count int) *RewindableGenerator {
	return &RewindableGenerator{
		generator: generator,
		count:     count,
	}
}

// GetIterator calls the generator function and returns the result.
func (rg *RewindableGenerator) GetIterator() []interface{} {
	return rg.generator()
}

// Count returns the count of tagged services.
func (rg *RewindableGenerator) Count() int {
	return rg.count
}
