package annotations

import (
	"encoding/json"
	"fmt"
)

const headerAnnotationPrefix = BfeAnnotationPrefix + "rewrite-header."

const (
	HeaderActionAnnotation      = headerAnnotationPrefix + "action"
)

/* examples: 
bfe.ingress.kubernetes.io/rewrite-header.actions: |-
 [
   {"cmd": "REQ_HEADER_ADD", "params": ["key", "value"]},
   {"cmd": "REQ_HEADER_SET", "params": ["key", "value"]},
   {"cmd": "REQ_HEADER_DEL", "params": ["key"]}
 ]
*/

type HeaderActions []HeaderAction
  
type HeaderAction struct {
	Cmd     string   // command of action
	Params  []string // params of action
}

/*
var headAnnotations = map[string]string{
	HeaderActionReqHeaderAdd:                 "REQ_HEADER_ADD",
	HeaderActionReqHeaderSet:                 "REQ_HEADER_SET",
	HeaderActionReqHeaderDel:                 "REQ_HEADER_DEL",
}
*/

// The annotation related to rewrite header action.
const (
	HeaderActionReqHeaderAdd = "REQ_HEADER_ADD"
	HeaderActionReqHeaderSet = "REQ_HEADER_SET"
	HeaderActionReqHeaderDel = "REQ_HEADER_DEL"
	HeaderActionRspHeaderAdd = "RSP_HEADER_ADD"
	HeaderActionRspHeaderSet = "RSP_HEADER_SET"
	HeaderActionRspHeaderDel = "RSP_HEADER_DEL"
)

// GetHeaderAction parse the cmd and the param of the header action from the annotations
func GetHeaderAction(annots map[string]string) (*HeaderActions, error) {
	if _, ok := annots[HeaderActionAnnotation]; !ok {
		return nil, nil
	}
	headerActions := make(HeaderActions, 0)
	if err := json.Unmarshal([]byte(annots[HeaderActionAnnotation]), &headerActions); err!= nil {
		return nil, fmt.Errorf("annotation %s's param is illegal, error: %s", annots[HeaderActionAnnotation], err)
	}

	return &headerActions, nil
}

