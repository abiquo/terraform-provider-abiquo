package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var deviceSchema = map[string]*schema.Schema{
	"devicetype":  Required().Renew().Link(),
	"endpoint":    Required().Renew().String(),
	"description": Optional().String(),
	"name":        Required().String(),
	"password":    Required().String(),
	"username":    Required().String(),
	"enterprise":  Optional().Renew().Link(),
	"datacenter":  Required().Link(),
}

var deviceResource = &schema.Resource{
	Schema: deviceSchema,
	Delete: resourceDelete,
	Exists: resourceExists("device"),
	Create: resourceCreate(deviceDTO, nil, deviceRead, deviceEndpoint),
	Update: resourceUpdate(deviceDTO, "device"),
	Read:   resourceRead(deviceDTO, deviceRead, "device"),
}

func deviceDTO(d *resourceData) core.Resource {
	return &abiquo.Device{
		Description: d.string("description"),
		Endpoint:    d.string("endpoint"),
		Name:        d.string("name"),
		Username:    d.string("username"),
		Password:    d.string("password"),
		Default:     true,
		DTO: core.NewDTO(
			d.link("enterprise"),
			d.link("devicetype"),
		),
	}
}

func deviceEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("datacenter")+"/devices", "device")
}

func deviceRead(d *resourceData, resource core.Resource) (err error) {
	device := resource.(*abiquo.Device)
	d.Set("endpoint", device.Endpoint)
	d.Set("name", device.Name)
	d.Set("password", device.Password)
	d.Set("username", device.Username)

	if _, ok := d.GetOk("description"); ok {
		d.Set("description", device.Description)
	}

	if _, ok := d.GetOk("enterprise"); ok {
		d.Set("enterprise", device.Rel("enterprise").URL())
	}

	return
}