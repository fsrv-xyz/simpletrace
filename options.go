package simpletrace

type SpanOption func(s *Span)

// OptionName - declare the name of the Span
func OptionName(name string) SpanOption {
	return func(s *Span) {
		s.Name = name
	}
}

// OptionUseKind - declare the kind of the Span
func OptionUseKind(kind Kind) SpanOption {
	return func(s *Span) {
		s.Kind = kind
	}
}

// OptionFromParent - set the parent ID of the Span
func OptionFromParent(parentId string) SpanOption {
	return func(s *Span) {
		s.ParentSpanId = parentId
	}
}

// OptionTraceID - set the TraceId of the Span
func OptionTraceID(id string) SpanOption {
	return func(s *Span) {
		s.TraceId = id
	}
}

// OptionSpanID - set the SpanId of the Span
func OptionSpanID(id string) SpanOption {
	return func(s *Span) {
		s.SpanId = id
	}
}

// OptionShared - set shared option to Span
func OptionShared() SpanOption {
	return func(s *Span) {
		s.Shared = true
	}
}

// OptionLocalEndpoint - define the local endpoint of the Span; parses the address to IPv6/IPv4 with port if set
func OptionLocalEndpoint(name, address string) SpanOption {
	return func(s *Span) {
		s.LocalEndpoint = Service{ServiceName: name}
		s.LocalEndpoint.separateAddresses(address)
	}
}

// OptionRemoteEndpoint - define the remote endpoint of the Span; parses the address to IPv6/IPv4 with port if set
func OptionRemoteEndpoint(name, address string) SpanOption {
	return func(s *Span) {
		s.RemoteEndpoint = Service{ServiceName: name}
		s.RemoteEndpoint.separateAddresses(address)
	}
}

// OptionTags - assign tags to the span
func OptionTags(tags map[string]string) SpanOption {
	return func(s *Span) {
		for key, value := range tags {
			s.Tag(key, value)
		}
	}
}
