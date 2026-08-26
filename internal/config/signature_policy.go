package config

type SignaturePolicy struct{ Allowed map[string]bool }

func LoadSignaturePolicy(names []string) *SignaturePolicy {
	if names == nil {
		return nil
	}
	p := &SignaturePolicy{Allowed: make(map[string]bool)}
	for _, name := range names {
		p.Allowed[name] = true
	}
	return p
}

func (p *SignaturePolicy) Allow(name string) bool { return p.Allowed[name] }
