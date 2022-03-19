package simpletrace

type SpanOption func(s *Span)

// UseKind - declare the kind of the Span
func UseKind(kind Kind) SpanOption {
	return func(s *Span) {
		s.Kind = kind
	}
}

// FromParent - set the parent ID of the Span
func FromParent(parentId string) SpanOption {
	return func(s *Span) {
		s.ParentSpanId = parentId
	}
}

// TraceID - set the TraceId of the Span
func TraceID(id string) SpanOption {
	return func(s *Span) {
		s.TraceId = id
	}
}

// Shared - set shared option to Span
func Shared() SpanOption {
	return func(s *Span) {
		s.Shared = true
	}
}

// LocalEndpoint - define the local endpoint of the Span; parses the address to IPv6/IPv4 with port if set
func LocalEndpoint(name, address string) SpanOption {
	return func(s *Span) {
		s.LocalEndpoint = Service{ServiceName: name}
		s.LocalEndpoint.separateAddresses(address)
	}
}

// RemoteEndpoint - define the remote endpoint of the Span; parses the address to IPv6/IPv4 with port if set
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
