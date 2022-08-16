// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

package cmd

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

const (
	descUse   = "descriptor"
	descShort = "Print credentials descriptor for Trustero internal use"
)

type desc struct {
	cmd *cobra.Command
}

func (d *desc) getCommand() *cobra.Command {
	return d.cmd
}

func (d *desc) setup() {
	d.cmd = &cobra.Command{
		Use:   descUse,
		Short: descShort,
		Args:  cobra.MinimumNArgs(0),
		RunE:  descriptor,
	}
}

// Cobra executes this function on descriptor command.
func descriptor(_ *cobra.Command, args []string) (err error) {
	var desc string
	if desc, err = toDescriptor(receptorImpl.GetCredentialObj()); err == nil {
		fmt.Println(desc)
	}
	return
}

type credential struct {
	Display     string `json:"display"`
	Placeholder string `json:"placeholder"`
	Field       string `json:"field"`
}

type credentials struct {
	Credentials []*credential `json:"credentials"`
}

func toDescriptor(credentialObj interface{}) (descriptor string, err error) {
	vt := reflect.Indirect(reflect.ValueOf(credentialObj)).Type()

	creds := &credentials{}
	for i := 0; i < vt.NumField(); i++ {
		tags := expandFieldTag(vt.Field(i))
		fname := vt.Field(i).Name
		display := getTagField(tags, displayField, fname)
		creds.Credentials = append(creds.Credentials, &credential{
			Display:     display,
			Field:       fname,
			Placeholder: getTagField(tags, placeholderField, strings.ToLower(fname)),
		})
	}

	var bytes []byte
	if bytes, err = json.MarshalIndent(creds, "", "  "); err == nil {
		descriptor = string(bytes)
	}
	return
}

func addCredentialFlags(credentialObj interface{}) (err error) {
	v := reflect.Indirect(reflect.ValueOf(credentialObj))
	vt := v.Type()
	for i := 0; i < vt.NumField(); i++ {
		tags := expandFieldTag(vt.Field(i))
		fname := vt.Field(i).Name
		display := getTagField(tags, displayField, fname)
		sptr := (*string)(reflect.Indirect(v.Field(i)).Addr().UnsafePointer())
		addStrFlagP("verify", sptr, strings.ToLower(fname), "", "", display)
		addStrFlagP("scan", sptr, strings.ToLower(fname), "", "", display)
	}
	return
}
