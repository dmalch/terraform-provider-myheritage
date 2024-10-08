package internal

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-myheritage/internal/myheritage"
)

type ProfileResource struct {
	resource.ResourceWithConfigure
	apiKey types.String
}

func NewProfileResource() resource.Resource {
	return &ProfileResource{}
}

// Metadata provides the resource type name
func (r *ProfileResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "myheritage_profile"
}

func (r *ProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Always perform a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	provider, ok := req.ProviderData.(*MyHeritageProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *MyHeritageProvider, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.apiKey = provider.apiKey
}

// Schema defines the schema for the resource
func (r *ProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"first_name": schema.StringAttribute{
				Required: true,
			},
			"last_name": schema.StringAttribute{
				Required: true,
			},
			"birth_date": schema.StringAttribute{
				Optional: true,
			},
			"birth_place": schema.StringAttribute{
				Optional: true,
			},
			"death_date": schema.StringAttribute{
				Optional: true,
			},
			"death_place": schema.StringAttribute{
				Optional: true,
			},
			"cause_of_death": schema.StringAttribute{
				Optional: true,
			},
			"gender": schema.StringAttribute{
				Optional: true,
			},
			"father_id": schema.StringAttribute{
				Optional: true,
			},
			"mother_id": schema.StringAttribute{
				Optional: true,
			},
			"individual_id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"events": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:      true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						},
						"type": schema.StringAttribute{
							Required: true,
						},
						"date": schema.StringAttribute{
							Optional: true,
						},
						"additional_content": schema.StringAttribute{
							Optional: true,
						},
						"formatted_place": schema.StringAttribute{
							Optional: true,
						},
						"title": schema.StringAttribute{
							Computed:      true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						},
					},
				},
			},
		},
	}
}

type ProfileResourceModel struct {
	ID           types.String `tfsdk:"id"`
	FirstName    types.String `tfsdk:"first_name"`
	LastName     types.String `tfsdk:"last_name"`
	IndividualID types.String `tfsdk:"individual_id"`
	BirthDate    types.String `tfsdk:"birth_date"`
	BirthPlace   types.String `tfsdk:"birth_place"`
	DeathDate    types.String `tfsdk:"death_date"`
	DeathPlace   types.String `tfsdk:"death_place"`
	CauseOfDeath types.String `tfsdk:"cause_of_death"`
	Gender       types.String `tfsdk:"gender"`
	FatherID     types.String `tfsdk:"father_id"`
	MotherID     types.String `tfsdk:"mother_id"`
	Events       types.List   `tfsdk:"events"`
}

type EventModel struct {
	ID                types.String `tfsdk:"id"`
	Type              types.String `tfsdk:"type"`
	Date              types.String `tfsdk:"date"`
	AdditionalContent types.String `tfsdk:"additional_content"`
	FormattedPlace    types.String `tfsdk:"formatted_place"`
	Title             types.String `tfsdk:"title"`
}

func eventModelObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                 types.StringType,
			"type":               types.StringType,
			"date":               types.StringType,
			"additional_content": types.StringType,
			"formatted_place":    types.StringType,
			"title":              types.StringType,
		},
	}
}

// Create creates the resource
func (r *ProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Access the API key from the provider's configuration
	familyTreeID, err := myheritage.CreateProfile(r.apiKey.ValueString(), plan.FirstName.ValueString(), plan.LastName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating profile", err.Error())
		return
	}

	plan.IndividualID = types.StringValue(familyTreeID)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read reads the resource
func (r *ProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile, err := myheritage.GetProfileHeader(r.apiKey.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading profile", err.Error())
		return
	}

	if profile.FirstName != "" {
		state.FirstName = types.StringValue(profile.FirstName)
	}
	if profile.LastName != "" {
		state.LastName = types.StringValue(profile.LastName)
	}
	if profile.Individual.Id != "" {
		state.IndividualID = types.StringValue(profile.Individual.Id)
	}
	if profile.Individual.BirthDate.Text != "" {
		state.BirthDate = types.StringValue(profile.Individual.BirthDate.Text)
	}
	if profile.Individual.BirthPlace != "" {
		state.BirthPlace = types.StringValue(profile.Individual.BirthPlace)
	}
	if profile.Individual.DeathDate.Text != "" {
		state.DeathDate = types.StringValue(profile.Individual.DeathDate.Text)
	}
	if profile.Individual.DeathPlace != "" {
		state.DeathPlace = types.StringValue(profile.Individual.DeathPlace)
	}
	if profile.Individual.Gender != "" {
		state.Gender = types.StringValue(profile.Individual.Gender)
	}
	if profile.Individual.CauseOfDeath != "" {
		state.CauseOfDeath = types.StringValue(profile.Individual.CauseOfDeath)
	}

	individualDetails, err := myheritage.GetProfileDetails(r.apiKey.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading profile details", err.Error())
		return
	}

	if fatherId := individualDetails.GetFatherId(); fatherId != "" {
		state.FatherID = types.StringValue(fatherId)
	}

	if motherId := individualDetails.GetMotherId(); motherId != "" {
		state.MotherID = types.StringValue(motherId)
	}

	// Prepare a list to store updated events
	var events []EventModel

	for _, eventFact := range individualDetails.EventFacts.Data {
		if eventFact.IsFactOfRelative {
			continue
		}

		var event EventModel

		event.ID = types.StringValue(eventFact.Id)
		event.Type = types.StringValue(eventFact.Type)
		event.Date = types.StringValue(eventFact.Date.Text)
		event.AdditionalContent = types.StringValue(eventFact.AdditionalContent)
		event.Title = types.StringValue(eventFact.Title)
		event.FormattedPlace = types.StringValue(eventFact.FormattedPlace)

		events = append(events, event)
	}

	// Convert the slice of EventModel to a types.List
	eventList, diags := types.ListValueFrom(ctx, eventModelObjectType(), events)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the event list in the state
	state.Events = eventList

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Update updates the resource
func (r *ProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := myheritage.UpdateProfile(r.apiKey.ValueString(), plan.IndividualID.ValueString(), plan.FirstName.ValueString(), plan.LastName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error updating profile", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes the resource
func (r *ProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := myheritage.DeleteProfile(r.apiKey.ValueString(), state.IndividualID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting profile", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}
