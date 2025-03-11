package e

import "log/slog"

func SlogErr(err error) slog.Attr {
	if err != nil {
		return slog.String("error", err.Error())
	}

	return slog.Attr{}
}
