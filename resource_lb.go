package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var algorithms = []string{"ROUND_ROBIN", "LEAST_CONNECTIONS", "SOURCE_IP"}
var protocols = []string{"TCP", "HTTP", "HTTPS"}

var lbSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"algorithm": &schema.Schema{
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateString(algorithms),
	},
	"routingrules": &schema.Schema{
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"protocolin": &schema.Schema{
					Required:     true,
					Type:         schema.TypeString,
					ValidateFunc: validateString(protocols),
				},
				"protocolout": &schema.Schema{
					Required:     true,
					Type:         schema.TypeString,
					ValidateFunc: validateString(protocols),
				},
				"portout": &schema.Schema{
					Required:     true,
					Type:         schema.TypeInt,
					ValidateFunc: validatePort,
				},
				"portin": &schema.Schema{
					Required:     true,
					Type:         schema.TypeInt,
					ValidateFunc: validatePort,
				},
			},
		},
		MinItems: 1,
		Required: true,
		Type:     schema.TypeList,
	},
	"privatenetwork": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
	"virtualdatacenter": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
}

var lbResource = &schema.Resource{
	Schema: lbSchema,
	Delete: resourceDelete,
	Exists: resourceExists("loadbalancer"),
	Create: resourceCreate(lbNew, nil, lbRead, lbEndpoint),
	Update: resourceUpdate(lbNew, nil, "loadbalancer"),
	Read:   resourceRead(lbNew, lbRead, "loadbalancer"),
}

func lbAddresses(d *resourceData) abiquo.LoadBalancerAddresses {
	return abiquo.LoadBalancerAddresses{
		Collection: []abiquo.LoadBalancerAddress{
			abiquo.LoadBalancerAddress{Internal: false},
		},
	}
}

func lbRules(d *resourceData) (rules []abiquo.RoutingRule) {
	for _, r := range d.slice("routingrules") {
		rule := abiquo.RoutingRule{}
		mapDecoder(r, &rule)
		rules = append(rules, rule)
	}
	return
}

func lbNew(d *resourceData) core.Resource {
	return &abiquo.LoadBalancer{
		Name:                  d.string("name"),
		Algorithm:             d.string("algorithm"),
		LoadBalancerAddresses: lbAddresses(d),
		RoutingRules: abiquo.RoutingRules{
			Collection: lbRules(d),
		},
		DTO: core.NewDTO(
			d.link("virtualdatacenter"),
			d.link("privatenetwork"),
		),
	}
}

func lbRead(d *resourceData, resource core.Resource) (err error) {
	lb := resource.(*abiquo.LoadBalancer)
	d.Set("name", lb.Name)
	d.Set("algorithm", lb.Algorithm)
	return
}

func lbUpdate(rd *schema.ResourceData, m interface{}) (err error) {
	d := newResourceData(rd, "loadbalancer")
	lb := lbNew(d).(*abiquo.LoadBalancer)
	if err = core.Update(d, lb); err == nil {
		err = lb.SetRules(lbRules(d))
	}
	return
}

func lbEndpoint(d *resourceData) (link *core.Link) {
	vdc := new(abiquo.VirtualDatacenter)
	if core.Read(d.link("virtualdatacenter"), vdc) == nil {
		endpoint := vdc.Rel("device")
		if endpoint == nil {
			return nil
		}

		device := new(abiquo.Device)
		if core.Read(endpoint, device) == nil {
			link = device.Rel("loadbalancers").SetType("loadbalancer")
		}
	}
	return
}
