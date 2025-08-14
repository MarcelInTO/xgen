// Copyright 2020 - 2024 The xgen Authors. All rights reserved. Use of this
// source code is governed by a BSD-style license that can be found in the
// LICENSE file.
//
// Package xgen written in pure Go providing a set of functions that allow you
// to parse XSD (XML schema files). This library needs Go version 1.10 or
// later.

package xgen

import "encoding/xml"

// OnAttributeGroup handles parsing event on the attributeGroup start
// elements. The attributeGroup element is used to group a set of attribute
// declarations so that they can be incorporated as a group into complex type
// definitions.
func (opt *Options) OnAttributeGroup(ele xml.StartElement, protoTree []interface{}) (err error) {
	attributeGroup := AttributeGroup{}
	for _, attr := range ele.Attr {
		if attr.Name.Local == "name" {
			attributeGroup.Name = attr.Value
		}
		if attr.Name.Local == "ref" {
			attributeGroup.Name = attr.Value
			attributeGroup.Ref, err = opt.GetValueType(attr.Value, protoTree)
			if err != nil {
				return
			}
		}
	}


	// Track nesting level
	opt.AttributeGroupNestLevel++

	// Check if we're inside another AttributeGroup (nested attributeGroup reference)
	if opt.AttributeGroup.Len() > 0 && opt.InAttributeGroup {
		// Add this attributeGroup reference to the parent attributeGroup
		parent := opt.AttributeGroup.Peek().(*AttributeGroup)
		parent.AttributeGroups = append(parent.AttributeGroups, attributeGroup)
		return
	}

	if opt.ComplexType.Len() == 0 {
		opt.InAttributeGroup = true
		opt.CurrentEle = opt.InElement
		opt.AttributeGroup.Push(&attributeGroup)
		return
	}

	if opt.ComplexType.Len() > 0 {
		opt.ComplexType.Peek().(*ComplexType).AttributeGroup = append(opt.ComplexType.Peek().(*ComplexType).AttributeGroup, attributeGroup)
		return
	}
	return
}

// EndAttributeGroup handles parsing event on the attributeGroup end elements.
func (opt *Options) EndAttributeGroup(ele xml.EndElement, protoTree []interface{}) (err error) {

	// Decrease nesting level
	opt.AttributeGroupNestLevel--

	// Only pop the stack when we return to the top level (nesting level 0) 
	// and we have an attribute group on the stack
	if opt.AttributeGroup.Len() > 0 && opt.AttributeGroupNestLevel == 0 && opt.InAttributeGroup {
		ag := opt.AttributeGroup.Pop().(*AttributeGroup)
		opt.ProtoTree = append(opt.ProtoTree, ag)
		opt.CurrentEle = ""
		opt.InAttributeGroup = false
	}
	return
}
