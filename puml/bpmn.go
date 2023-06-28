package puml

// BpmnSymbol is a unicode code point from the bpmn font,
// see https://cdn.staticaly.com/gh/bpmn-io/bpmn-font/master/dist/demo.html
type BpmnSymbol string

const (
	bpmn_icon_start_event_none BpmnSymbol = "e845"
	bpmn_icon_task             BpmnSymbol = "e821"
)
