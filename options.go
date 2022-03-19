package simpletrace

type SpanOption func(s *Span)

// UseKind - declare the kind of the span
func UseKind(kind Kind) SpanOption {
	return func(s *Span) {
		s.Kind = kind
	}
}

// Shared - set shared option to span
func Shared() SpanOption {
	return func(s *Span) {
		s.Shared = true
	}
}

// LocalEndpoint - define the local endpoint of the span; parses the address to IPv6/IPv4 with port if set
func LocalEndpoint(name, address string) SpanOption {
	return func(s *Span) {
		s.LocalEndpoint = Service{ServiceName: name}
		s.LocalEndpoint.separateAddresses(address)
	}
}

// RemoteEndpoint - define the remote endpoint of the span; parses the address to IPv6/IPv4 with port if set
func RemoteEndpoint(name, address string) SpanOption {
	return func(s *Span) {
		s.RemoteEndpoint = Service{ServiceName: name}
		s.RemoteEndpoint.separateAddresses(address)
	}
}

// Tags - assign tags to the span
func Tags(tags map[string]string) SpanOption {
	return func(s *Span) {
		for key, value := range tags {
			s.Tag(key, value)
		}
	}
}
