package lib

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Option struct {
	P         any
	Name      string
	Shorthand string
	Value     any
	Usage     string
}

var options = []Option{}

func AddOption(cmd *cobra.Command, o Option) {
	switch o.Value.(type) {
	case string:
		if o.Shorthand == "" {
			cmd.PersistentFlags().StringVar(o.P.(*string), o.Name, o.Value.(string), o.Usage)
		} else {
			cmd.PersistentFlags().StringVarP(o.P.(*string), o.Name, o.Shorthand, o.Value.(string), o.Usage)
		}
	case int:
		if o.Shorthand == "" {
			cmd.PersistentFlags().IntVar(o.P.(*int), o.Name, o.Value.(int), o.Usage)
		} else {
			cmd.PersistentFlags().IntVarP(o.P.(*int), o.Name, o.Shorthand, o.Value.(int), o.Usage)
		}
	default:
		return
	}
	viper.BindPFlag(o.Name, cmd.PersistentFlags().Lookup(o.Name))
	options = append(options, o)
}

func LoadOptions() {
	for _, o := range options {
		switch o.Value.(type) {
		case string:
			*o.P.(*string) = viper.GetString(o.Name)
		case int:
			*o.P.(*int) = viper.GetInt(o.Name)
		}
	}
}

// type StringVarP struct {
// 	P         *string
// 	Name      string
// 	Shorthand string
// 	Value     string
// 	Usage     string
// }

// func AddStringVarPs(cmd *cobra.Command, options []StringVarP) {
// 	for _, o := range options {
// 		cmd.PersistentFlags().StringVarP(o.P, o.Name, o.Shorthand, o.Value, o.Usage)
// 		viper.BindPFlag(o.Name, cmd.PersistentFlags().Lookup(o.Name))
// 	}
// }
// func LoadStringVarPs(options []StringVarP) {
// 	for _, o := range options {
// 		*o.P = viper.GetString(o.Name)
// 	}
// }

// type IntVarP struct {
// 	P         *int
// 	Name      string
// 	Shorthand string
// 	Value     int
// 	Usage     string
// }

// func AddIntVarPs(cmd *cobra.Command, options []IntVarP) {
// 	for _, o := range options {
// 		cmd.PersistentFlags().IntVarP(o.P, o.Name, o.Shorthand, o.Value, o.Usage)
// 		viper.BindPFlag(o.Name, cmd.PersistentFlags().Lookup(o.Name))
// 	}
// }
// func LoadIntVarPs(options []IntVarP) {
// 	for _, o := range options {
// 		*o.P = viper.GetInt(o.Name)
// 	}
// }
