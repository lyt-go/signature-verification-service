package logger

type AuditContext struct {
	Subject  string
	Verifier string
}

type AuditEvent struct {
	Subject  string
	Verifier string
}

func (c *AuditContext) Event() AuditEvent {
	return AuditEvent{Subject: c.Subject, Verifier: c.Verifier}
}

type AuditSink interface{ Emit(*AuditContext) }
