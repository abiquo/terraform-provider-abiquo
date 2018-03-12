package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var limitResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		// Soft limits
		"cpusoft":  Optional().Number(),
		"hdsoft":   Optional().Number(),
		"ipsoft":   Optional().Number(),
		"ramsoft":  Optional().Number(),
		"reposoft": Optional().Number(),
		"volsoft":  Optional().Number(),
		"vlansoft": Optional().Number(),
		// Hard limits
		"cpuhard":  Optional().Number(),
		"hdhard":   Optional().Number(),
		"iphard":   Optional().Number(),
		"ramhard":  Optional().Number(),
		"repohard": Optional().Number(),
		"vlanhard": Optional().Number(),
		"volhard":  Optional().Number(),
		// Links
		"location":   Required().Renew().Link(),
		"enterprise": Required().Renew().Link(),
		"hwprofiles": &schema.Schema{
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateURL,
			},
			Optional: true,
			Set:      schema.HashString,
			Type:     schema.TypeSet,
		},
		"backups": &schema.Schema{
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateURL,
			},
			Optional: true,
			Set:      schema.HashString,
			Type:     schema.TypeSet,
		},
	},
	Exists: resourceExists("limit"),
	Read:   resourceRead(limitNew, limitRead, "limit"),
	Update: resourceUpdate(limitNew, nil, "limit"),
	Create: resourceCreate(limitNew, nil, limitRead, limitEndpoint),
	Delete: resourceDelete,
}

func limitNew(d *resourceData) core.Resource {
	limit := &abiquo.Limit{
		// Soft limits
		CPUSoft:  d.int("cpusoft"),
		HDSoft:   d.int("hdsoft"),
		IPSoft:   d.int("ipsoft"),
		RAMSoft:  d.int("ramsoft"),
		RepoSoft: d.int("reposoft"),
		VolSoft:  d.int("VolSoft"),
		VLANSoft: d.int("vlansoft"),
		// Hard limits
		CPUHard:  d.int("cpuhard"),
		HDHard:   d.int("hdhard"),
		IPHard:   d.int("iphard"),
		RAMHard:  d.int("ramhard"),
		RepoHard: d.int("repohard"),
		VolHard:  d.int("volhard"),
		VLANHard: d.int("vlanhard"),
		// Links
		DTO: core.NewDTO(
			d.link("location"),
		),
	}

	// Backups
	backups := d.set("backups")
	if backups != nil && backups.Len() > 0 {
		for _, entry := range backups.List() {
			href := entry.(string)
			limit.Add(core.NewLinkType(href, "backuppolicy").SetRel("backuppolicy"))
		}
	}

	// HWprofiles
	hwprofiles := d.set("hwprofiles")
	if hwprofiles != nil && hwprofiles.Len() > 0 {
		for _, entry := range d.set("hwprofiles").List() {
			href := entry.(string)
			limit.Add(core.NewLinkType(href, "hardwareprofile").SetRel("hardwareprofile"))
		}
	}

	return limit
}

func limitEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("enterprise")+"/limits", "limit")
}

func limitRead(d *resourceData, resource core.Resource) (err error) {
	limit := resource.(*abiquo.Limit)

	backups := mapHrefs(limit.Links.Filter(func(l *core.Link) bool {
		return l.IsMedia("backuppolicy")
	}))

	hwprofiles := mapHrefs(limit.Links.Filter(func(l *core.Link) bool {
		return l.IsMedia("hwprofile")
	}))

	d.Set("backups", backups)
	d.Set("hwprofiles", hwprofiles)
	// Soft limits
	d.Set("cpusoft", limit.CPUSoft)
	d.Set("hdsoft", limit.HDSoft)
	d.Set("ipsoft", limit.IPSoft)
	d.Set("ramsoft", limit.RAMSoft)
	d.Set("reposoft", limit.RepoSoft)
	d.Set("volsoft", limit.VolSoft)
	d.Set("vlansoft", limit.VLANSoft)
	// Hard limits
	d.Set("cpuhard", limit.CPUHard)
	d.Set("hdhard", limit.HDHard)
	d.Set("iphard", limit.IPHard)
	d.Set("ramhard", limit.RAMHard)
	d.Set("repohard", limit.RepoHard)
	d.Set("volhard", limit.VolHard)
	d.Set("vlanhard", limit.VLANHard)
	return
}
