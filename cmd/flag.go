package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// Flags is a flag with a protobuf type
type Flags struct {
	EnumName  map[int32]string
	EnumValue map[string]int32
	Usage     string
	Value     int32
}

// GetUsage returns the usage string for the flag
func (f *Flags) GetUsage() string {
	return fmt.Sprintf("%s (%s)", f.Usage, f.listString())
}

func (f Flags) String() string {
	return f.EnumName[f.Value]
}

func (f *Flags) listString() string {
	var keys []string
	for key := range f.EnumValue {
		keys = append(keys, key)
	}
	return strings.Join(keys, "|")
}

// Set the flags with the given value
func (f *Flags) Set(value string) error {
	val, ok := f.EnumValue[value]
	if !ok {
		return fmt.Errorf("could not parse %q, expected one of %s", value, f.listString())
	}
	f.Value = val
	return nil
}

func (f *Flags) flagCompletion(toComplete string) ([]string, cobra.ShellCompDirective) {
	var complete []string
	for key := range f.EnumValue {
		if strings.HasPrefix(key, toComplete) {
			complete = append(complete, key)
		}
	}
	if len(complete) > 0 {
		return complete, cobra.ShellCompDirectiveNoFileComp
	}
	return nil, cobra.ShellCompDirectiveError | cobra.ShellCompDirectiveNoFileComp
}

// RegisterFlagCompletionFunc return lambda to complete the end of the flag
func (f *Flags) RegisterFlagCompletionFunc() func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return f.flagCompletion(toComplete)
	}
}

// CobraRegister flags and flag completion func
func (f *Flags) CobraRegister(cmd *cobra.Command, name string) {
	cmd.Flags().Var(f, name, f.GetUsage())
	cmd.RegisterFlagCompletionFunc(name, f.RegisterFlagCompletionFunc())
}

// Type don't known maybe underlying type
func (f *Flags) Type() string {
	return "string"
}
