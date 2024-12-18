package lib

import (
	"time"

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
	case bool:
		if o.Shorthand == "" {
			cmd.PersistentFlags().BoolVar(o.P.(*bool), o.Name, o.Value.(bool), o.Usage)
		} else {
			cmd.PersistentFlags().BoolVarP(o.P.(*bool), o.Name, o.Shorthand, o.Value.(bool), o.Usage)
		}
	case time.Duration:
		if o.Shorthand == "" {
			cmd.PersistentFlags().DurationVar(o.P.(*time.Duration), o.Name, o.Value.(time.Duration), o.Usage)
		} else {
			cmd.PersistentFlags().DurationVarP(o.P.(*time.Duration), o.Name, o.Shorthand, o.Value.(time.Duration), o.Usage)
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
		case bool:
			*o.P.(*bool) = viper.GetBool(o.Name)
		case time.Duration:
			*o.P.(*time.Duration) = viper.GetDuration(o.Name)
		}
	}
}
