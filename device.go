package awair

import (
	"context"
	"fmt"
	"net/url"
)

type Device struct {
	Name         string  `json:"name"`
	Id           int64   `json:"deviceId"`
	UUID         string  `json:"deviceUUID"`
	Type         string  `json:"deviceType"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Preference   string  `json:"preference"`
	Timezone     string  `json:"timezone"`
	RoomType     string  `json:"roomType"`
	SpaceType    string  `json:"spaceType"`
	LocationName string  `json:"locationName"`

	c *Client `json:"-"`
}

func (d *Device) baseURL(userPrefix bool) string {
	var user string
	if userPrefix {
		user = "users/self/"
	}
	return fmt.Sprintf("https://developer-apis.awair.is/v1/%sdevices/%s/%d",
		user, d.Type, d.Id)
}

func (d *Device) Latest(ctx context.Context) (*Observation, error) {
	var wrapper struct {
		Data []*Observation `json:"data"`
	}
	err := d.c.get(ctx, d.baseURL(true)+"/air-data/latest?"+
		(url.Values{
			"fahrenheit": []string{fmt.Sprint(d.c.Options.PreferFahrenheit)},
		}).Encode(), &wrapper)
	if err != nil {
		return nil, err
	}
	if len(wrapper.Data) != 1 {
		return nil, fmt.Errorf("unexpected situation")
	}
	return wrapper.Data[0], nil
}
