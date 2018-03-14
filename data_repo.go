package main

import (
	"errors"
	"fmt"

	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var repoDataSchema = map[string]*schema.Schema{
	"datacenter": &schema.Schema{
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
}

var repoDataSource = &schema.Resource{
	Schema: repoDataSchema,
	Read:   dataRepoRead,
}

func dataRepoRead(d *schema.ResourceData, p interface{}) (err error) {
	enterprise := p.(*provider).Enterprise()
	if enterprise == nil {
		return errors.New("The user enterprise was not found")
	}

	finder := func(r core.Resource) bool {
		return r.Rel("datacenter").Title == d.Get("datacenter").(string)
	}
	repo := enterprise.Rel("datacenterrepositories").Collection(nil).Find(finder)
	if repo == nil {
		return fmt.Errorf("datacenter repository for datacenter %q was not found", d.Get("datacenter"))
	}

	d.SetId(repo.URL())
	return
}
