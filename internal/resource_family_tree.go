package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/dmalch/terraform-provider-myheritage/internal/myheritage"
)

func resourceFamilyTree() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFamilyTreeCreate,
		ReadContext:   resourceFamilyTreeRead,
		UpdateContext: resourceFamilyTreeUpdate,
		DeleteContext: resourceFamilyTreeDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceFamilyTreeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	apiKey := m.(string) // Retrieve the API key from the meta interface

	// Call MyHeritage API to create the family tree
	// Assume you have a function createFamilyTree that interacts with the API
	familyTreeID, err := myheritage.CreateFamilyTree(apiKey, name, description)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(familyTreeID)

	return diags
}

func resourceFamilyTreeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	familyTreeID := d.Id()
	apiKey := m.(string) // Retrieve the API key from the meta interface

	// Call MyHeritage API to read the family tree
	// Assume you have a function getFamilyTree that interacts with the API
	familyTree, err := myheritage.GetFamilyTree(apiKey, familyTreeID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", familyTree.Name)
	d.Set("description", familyTree.Description)

	return diags
}

func resourceFamilyTreeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	familyTreeID := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	apiKey := m.(string) // Retrieve the API key from the meta interface

	// Call MyHeritage API to update the family tree
	// Assume you have a function updateFamilyTree that interacts with the API
	err := myheritage.UpdateFamilyTree(apiKey, familyTreeID, name, description)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceFamilyTreeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	familyTreeID := d.Id()
	apiKey := m.(string) // Retrieve the API key from the meta interface

	// Call MyHeritage API to delete the family tree
	// Assume you have a function deleteFamilyTree that interacts with the API
	err := myheritage.DeleteFamilyTree(apiKey, familyTreeID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
