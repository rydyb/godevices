package main

type Singletone struct {
	Channel           uint8   `name:"channel" required:"" help:"The channel number to configure."`
	LogAmplitudeScale float64 `name:"log-amplitude" required:"" help:"The relative singletone amplitude in dBm."`
	Frequency         float64 `name:"frequency" required:"" help:"The frequency of the singletone in Hz."`
}

func (cmd *Singletone) Run(ctx *Context) error {
	return ctx.FlexDDSSlot.Singletone(cmd.Channel, cmd.LogAmplitudeScale, cmd.Frequency)
}
