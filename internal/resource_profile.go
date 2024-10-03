package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/dmalch/terraform-provider-myheritage/internal/myheritage"
)

func resourceProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProfileCreate,
		ReadContext:   resourceProfileRead,
		UpdateContext: resourceProfileUpdate,
		DeleteContext: resourceProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceProfileImport,
		},
		Schema: map[string]*schema.Schema{
			"first_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"birth_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"birth_place": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"death_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"death_place": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cause_of_death": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gender": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"father_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mother_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"individual_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	apiKey := m.(string) // Retrieve the API key from the meta interface

	// Call MyHeritage API to create the family tree
	// Assume you have a function createFamilyTree that interacts with the API
	familyTreeID, err := myheritage.CreateProfile(apiKey, name, description)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(familyTreeID)

	return diags
}

func resourceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	profileID := d.Id()
	apiKey := m.(string) // Retrieve the API key from the meta interface

	d, err := retrieveProfile(apiKey, profileID, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func retrieveProfile(apiKey string, profileID string, d *schema.ResourceData) (*schema.ResourceData, error) {
	profile, err := myheritage.GetProfileHeader(apiKey, profileID)
	if err != nil {
		return nil, err
	}

	if err := d.Set("first_name", profile.FirstName); err != nil {
		return nil, err
	}

	if err := d.Set("last_name", profile.LastName); err != nil {
		return nil, err
	}

	if err := d.Set("birth_date", profile.Individual.BirthDate.Text); err != nil {
		return nil, err
	}

	if err := d.Set("death_date", profile.Individual.DeathDate.Text); err != nil {
		return nil, err
	}

	if err := d.Set("birth_place", profile.Individual.BirthPlace); err != nil {
		return nil, err
	}

	if err := d.Set("death_place", profile.Individual.DeathPlace); err != nil {
		return nil, err
	}

	if err := d.Set("cause_of_death", profile.Individual.CauseOfDeath); err != nil {
		return nil, err
	}

	if err := d.Set("gender", profile.Individual.Gender); err != nil {
		return nil, err
	}

	if err := d.Set("individual_id", profile.Individual.ID); err != nil {
		return nil, err
	}

	individualDetails, err := myheritage.GetProfileDetails(apiKey, profileID)
	if err != nil {
		return nil, err
	}

	if fatherId := individualDetails.GetFatherId(); fatherId != "" {
		if err := d.Set("father_id", fatherId); err != nil {
			return nil, err
		}
	}

	if motherId := individualDetails.GetMotherId(); motherId != "" {
		if err := d.Set("mother_id", motherId); err != nil {
			return nil, err
		}
	}

	return d, nil
}

func resourceProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	familyTreeID := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	apiKey := m.(string) // Retrieve the API key from the meta interface

	// Call MyHeritage API to update the family tree
	// Assume you have a function updateFamilyTree that interacts with the API
	err := myheritage.UpdateProfile(apiKey, familyTreeID, name, description)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	familyTreeID := d.Id()
	apiKey := m.(string) // Retrieve the API key from the meta interface

	// Call MyHeritage API to delete the family tree
	// Assume you have a function deleteFamilyTree that interacts with the API
	err := myheritage.DeleteProfile(apiKey, familyTreeID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func resourceProfileImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	profileID := d.Id()
	apiKey := m.(string) // Retrieve the API key from the meta interface

	d, err := retrieveProfile(apiKey, profileID, d)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
