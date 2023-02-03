package header

import (
	"testing"
	"fmt"

//	netv1 "k8s.io/api/networking/v1"
	"github.com/bfenetworks/ingress-bfe/internal/bfeConfig/annotations"
//	"github.com/bfenetworks/ingress-bfe/internal/bfeConfig/configs/cache"
//	"github.com/bfenetworks/ingress-bfe/internal/bfeConfig/util"
//	"github.com/bfenetworks/bfe/bfe_modules/mod_header"
)

func TestAnnotationExtraction(t *testing.T) {
	m := make(map[string]string)
	m[annotations.HeaderActionAnnotation] = `[
		{"cmd": "REQ_HEADER_ADD", "params": ["key1", "value1"]},
		{"cmd": "REQ_HEADER_DEL", "params": ["key11", "value11", "value22"]},
		{"cmd": "REQ_HEADER_SET", "params": ["key2", "value2"]},
		{"cmd": "REQ_HEADER_DEL", "params": ["key3"]},
	]`
//		{"cmd": "REQ_HEADER_ADD", "params": ["key11", "value11", "value22"]},

	fmt.Printf("Data = %s\n\n", m[annotations.HeaderActionAnnotation])

	headerActions, err := annotations.GetHeaderAction(m)
	if err != nil {
		fmt.Printf("Error occured: %s", err)
		return
	}

	if headerActions == nil {
		fmt.Printf("Parse failed")
		return
	}

	if err := checkAction(headerActions); err != nil {
		fmt.Printf("Basic check Failed: %s\n\n", err.Error())
	} else {
		fmt.Printf("Basic check Passed\n\n")
	}

	for _, headerAction := range *headerActions {
		fmt.Printf("Cmd = %s\n", headerAction.Cmd)
//		fmt.Printf("Params = %s, Params[0] = %s, Params[1] = %s\n", headerAction.Params, headerAction.Params[0], headerAction.Params[1])
		fmt.Printf("Params = %s, Params[0] = %s\n", headerAction.Params, headerAction.Params[0])
	}

	return
}
