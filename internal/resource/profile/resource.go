package profile

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/dmalch/terraform-provider-myheritage/internal/config"
	"github.com/dmalch/terraform-provider-myheritage/internal/myheritage"
)

type Resource struct {
	resource.ResourceWithConfigure
	apiKey types.String
}

func NewProfileResource() resource.Resource {
	return &Resource{}
}

// Metadata provides the resource type name
func (r *Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "myheritage_profile"
}

func (r *Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Always perform a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(*config.MyHeritageProviderConfig)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *config.MyHeritageProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.apiKey = cfg.ApiKey
}

type ResourceModel struct {
	ID           types.String `tfsdk:"id"`
	FirstName    types.String `tfsdk:"first_name"`
	LastName     types.String `tfsdk:"last_name"`
	IndividualID types.String `tfsdk:"individual_id"`
	Gender       types.String `tfsdk:"gender"`
	FatherID     types.String `tfsdk:"father_id"`
	MotherID     types.String `tfsdk:"mother_id"`
	Events       types.List   `tfsdk:"events"`
	Notes        types.List   `tfsdk:"notes"`
}

type EventModel struct {
	ID                types.String `tfsdk:"id"`
	Type              types.String `tfsdk:"type"`
	Date              types.String `tfsdk:"date"`
	Content           types.String `tfsdk:"content"`
	AdditionalContent types.String `tfsdk:"additional_content"`
	FormattedPlace    types.String `tfsdk:"formatted_place"`
	Title             types.String `tfsdk:"title"`
	CauseOfDeath      types.String `tfsdk:"cause_of_death"`
	SpouseId          types.String `tfsdk:"spouse_id"`
	Notes             types.List   `tfsdk:"notes"`
	Media             types.List   `tfsdk:"media"`
}

func eventModelObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                 types.StringType,
			"type":               types.StringType,
			"date":               types.StringType,
			"content":            types.StringType,
			"additional_content": types.StringType,
			"formatted_place":    types.StringType,
			"title":              types.StringType,
			"cause_of_death":     types.StringType,
			"spouse_id":          types.StringType,
			"notes": types.ListType{
				ElemType: noteModelObjectType(),
			},
			"media": types.ListType{
				ElemType: mediaModelObjectType(),
			},
		},
	}
}

type NoteModel struct {
	ID   types.String `tfsdk:"id"`
	Text types.String `tfsdk:"text"`
}

func noteModelObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":   types.StringType,
			"text": types.StringType,
		},
	}
}

type MediaModel struct {
	Name types.String `tfsdk:"name"`
	Link types.String `tfsdk:"link"`
}

func mediaModelObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name": types.StringType,
			"link": types.StringType,
		},
	}
}

// Create creates the resource
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ID, err := myheritage.CreateProfile(r.apiKey.ValueString(), plan.FirstName.ValueString(), plan.LastName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating profile", err.Error())
		return
	}

	plan.ID = types.StringValue(ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read reads the resource
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceModel
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
	if profile.Individual.Gender != "" {
		state.Gender = types.StringValue(profile.Individual.Gender)
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

		noteList, diags := notesToList(ctx, eventFact.Notes.Data)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		mediaList, diags := mediaToList(ctx, eventFact.Media.Data)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		event := EventModel{
			ID:                types.StringValue(eventFact.Id),
			Type:              types.StringValue(eventFact.Type),
			Date:              types.StringValue(eventFact.Date.Text),
			AdditionalContent: types.StringValue(eventFact.AdditionalContent),
			Title:             types.StringValue(eventFact.Title),
			FormattedPlace:    types.StringValue(eventFact.FormattedPlace),
			CauseOfDeath:      types.StringValue(eventFact.CauseOfDeath),
			SpouseId:          getSpouseId(eventFact),
			Content:           types.StringValue(eventFact.Content),
			Notes:             noteList,
			Media:             mediaList,
		}

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

	individualBiography, err := myheritage.GetIndividualBiography(r.apiKey.ValueString(), state.IndividualID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading biography details", err.Error())
		return
	}

	noteList, diags := notesToList(ctx, individualBiography.Notes.Data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Notes = noteList

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func getSpouseId(eventFact myheritage.EventFact) basetypes.StringValue {
	if eventFact.Spouse != nil && eventFact.Spouse.Id != "" {
		return types.StringValue(eventFact.Spouse.Id)
	}

	return types.StringNull()
}

func notesToList(ctx context.Context, noteRecords []myheritage.Note) (basetypes.ListValue, diag.Diagnostics) {
	var noteModels []NoteModel

	for _, noteRecord := range noteRecords {
		noteModels = append(noteModels, NoteModel{
			ID:   types.StringValue(noteRecord.Id),
			Text: types.StringValue(noteRecord.Text),
		})
	}

	// Convert the slice of NoteModel to a types.List
	return types.ListValueFrom(ctx, noteModelObjectType(), noteModels)
}

func mediaToList(ctx context.Context, mediaRecords []myheritage.Media) (basetypes.ListValue, diag.Diagnostics) {
	var mediaModels []MediaModel

	for _, mediaRecord := range mediaRecords {
		mediaModels = append(mediaModels, MediaModel{
			Name: types.StringValue(mediaRecord.Name),
			Link: types.StringValue(mediaRecord.Link),
		})
	}

	return types.ListValueFrom(ctx, mediaModelObjectType(), mediaModels)
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Update updates the resource
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ResourceModel
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
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ResourceModel
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
