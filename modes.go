package awair

import (
	"context"
)

func (d *Device) GetDisplayMode(ctx context.Context) (string, error) {
	var wrapper struct {
		Mode string `json:"mode"`
	}

	return wrapper.Mode, d.c.get(ctx, d.baseURL(false)+"/display", &wrapper)
}

func (d *Device) SetDisplayMode(ctx context.Context, mode string) error {
	return d.c.put(ctx, d.baseURL(false)+"/display", map[string]string{"mode": mode})
}

func (d *Device) GetKnockingMode(ctx context.Context) (string, error) {
	var wrapper struct {
		Mode string `json:"mode"`
	}

	return wrapper.Mode, d.c.get(ctx, d.baseURL(false)+"/knocking", &wrapper)
}

func (d *Device) SetKnockingMode(ctx context.Context, mode string) error {
	return d.c.put(ctx, d.baseURL(false)+"/knocking", map[string]string{"mode": mode})
}

func (d *Device) GetLEDMode(ctx context.Context) (string, error) {
	var wrapper struct {
		Mode string `json:"mode"`
	}

	return wrapper.Mode, d.c.get(ctx, d.baseURL(false)+"/led", &wrapper)
}

func (d *Device) SetLEDMode(ctx context.Context, mode string) error {
	return d.c.put(ctx, d.baseURL(false)+"/led", map[string]string{"mode": mode})
}
