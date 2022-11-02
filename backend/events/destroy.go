package events

func (e *DestroyEvent) PreHook(p *Processor) error {
	if e.PreHookFunc == nil {
		return nil
	}
	return e.PreHookFunc()
}

func (e *DestroyEvent) PostHook(p *Processor) error {
	if e.PostHookFunc == nil {
		return nil
	}
	return e.PostHookFunc()
}

func (e *DestroyEvent) Process(p *Processor) error {
	return p.deleteNode(string(e.NodeType), e.Prop, e.Value)
}
