// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

import fmt "fmt"

func (dst *Status) SetFields(src *Status, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "code":
			if len(subs) > 0 {
				return fmt.Errorf("'code' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Code = src.Code
			} else {
				var zero int32
				dst.Code = zero
			}
		case "message":
			if len(subs) > 0 {
				return fmt.Errorf("'message' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Message = src.Message
			} else {
				var zero string
				dst.Message = zero
			}
		case "details":
			if len(subs) > 0 {
				return fmt.Errorf("'details' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Details = src.Details
			} else {
				dst.Details = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}
