package std

// A Statistic is a kind of TemporalBuffer that -only- tracks the number of times something is hit within a window of observance.
// For example - Every Epiphany will add an entry to the Statistic stored in std.Path{"Performance", "Epiphanies"} whenever it "materializes" into something.
//
// NOTE: You can still query the temporal information for the exact instant of each event, but
// there are no guarantees (or requirements) of what data the actor associates with it.
type Statistic struct {
	*TemporalBuffer[any]

	created bool
}

func NewStatistic() *Statistic {
	s := &Statistic{
		TemporalBuffer: NewTemporalBuffer[any](),
		created:        true,
	}
	s.TemporalBuffer.sanityPassthrough = s.sanityCheck
	return s
}

func (s *Statistic) sanityCheck() {
	if !s.created {
		panic("please create a std.Statistic through std.NewStatistic()")
	}
}
