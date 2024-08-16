// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

package cmd

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
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
		Use:          descUse,
		Short:        descShort,
		Args:         cobra.MinimumNArgs(0),
		RunE:         descriptor,
		SilenceUsage: true,
	}
	d.cmd.FParseErrWhitelist.UnknownFlags = true
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

type descriptors struct {
	Credentials  []*credential `json:"credentials"`
	Config       interface{}   `json:"config"`
	ReceptorType string        `json:"receptorType"`
}

func toDescriptor(credentialObj interface{}) (descriptor string, err error) {
	vt := reflect.Indirect(reflect.ValueOf(credentialObj)).Type()

	creds := &descriptors{}
	for i := 0; i < vt.NumField(); i++ {
		tags := expandFieldTag(vt.Field(i))
		fname := vt.Field(i).Name
		display := getTagField(tags, displayField, fname)
		if strings.ToLower(fname) == "oauth" {
			fname = "oauth"
			display = ""
		}

		placeholder := getTagField(tags, placeholderField, strings.ToLower(fname))
		if display != "" || fname == "oauth" {
			creds.Credentials = append(creds.Credentials, &credential{
				Display:     display,
				Field:       fname,
				Placeholder: placeholder,
			})
		}
	}

	creds.ReceptorType = GetParsedReceptorType()
	creds.Config = GetReceptorConfig()
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

// GetParsedReceptorType returns a normalized receptor type string.  A receptor type string only includes letters,
// numbers, "-", and "_".  All other characters will be converted to "_".
func GetParsedReceptorType() (parsedName string) {
	receptorName := receptorImpl.GetReceptorType()
	regex, _ := regexp.Compile(`[^-a-z0-9A-Z_]`)
	res := regex.ReplaceAll([]byte(receptorName), []byte("_"))
	return string(res)
}

func GetReceptorConfig() interface{} {
	return receptorImpl.GetConfigObj()
}
