package puml

import "fmt"

// BpmnSymbol is a unicode code point from the bpmn font,
// see https://cdn.staticaly.com/gh/bpmn-io/bpmn-font/master/dist/demo.html
type BpmnSymbol string

const (
	bpmn_icon_start_event_none     BpmnSymbol = "e845"
	bpmn_icon_task                 BpmnSymbol = "e821"
	bpmn_icon_gateway_xor          BpmnSymbol = "e80f"
	bpmn_icon_end_event_terminate  BpmnSymbol = "e836"
	bpmn_icon_end_event_message    BpmnSymbol = "e83a"
	bpmn_icon_end_event_error      BpmnSymbol = "e822"
	bpmn_icon_end_event_cancel     BpmnSymbol = "e811"
	bpmn_icon_receive              BpmnSymbol = "e829"
	bpmn_icon_data_object          BpmnSymbol = "e84b"
	bpmn_icon_data_output                     = "e867"
	bpmn_icon_data_input                      = "e866"
	bpmn_icon_text_annotation                 = "e86b"
	bpmn_icon_manual_task                     = "e840"
	bpmn_icon_subprocess_collapsed BpmnSymbol = "e81f"
	bpmn_icon_loop_marker          BpmnSymbol = "e809"
)

const (
	ColorExternFunc                = "#f6d8ec"
	ColorInternFunc                = "#d8f6f3"
	ColorEvent                     = "#ff992a"
	ColorData                      = "#f5e253"
	ColorErrorEvent                = "#ec4d4e"
	ColorTask                      = "#3399fe"
	ColorTaskUndefined             = "#94c9ff"
	ColorIntermediateFnResultEvent = "#fccd9a"
)

func bpmSym(symbol BpmnSymbol) string {
	return fmt.Sprintf("<size:25><font:bpmn><U+%s></font></size>", symbol)
}
