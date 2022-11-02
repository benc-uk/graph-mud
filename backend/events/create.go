package events

func (e *CreateEvent) PreHook(p *Processor) error {
	if e.PreHookFunc == nil {
		return nil
	}
	return e.PreHookFunc()
}

func (e *CreateEvent) PostHook(p *Processor) error {
	if e.PostHookFunc == nil {
		return nil
	}
	return e.PostHookFunc()
}

func (e *CreateEvent) Process(p *Processor) error {
	return p.createNode(string(e.Type), e.Props)
}
