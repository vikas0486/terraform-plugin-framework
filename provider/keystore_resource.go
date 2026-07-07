package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type KeystoreResource struct {
	client *APIClient
}

type KeystoreResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func NewKeystoreResource() resource.Resource {
	return &KeystoreResource{}
}

func (r *KeystoreResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {

	resp.TypeName = "thales_keystore"

}

func (r *KeystoreResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*APIClient)

	if !ok {

		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			"Expected *APIClient.",
		)

		return
	}

	r.client = client
}

func (r *KeystoreResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {

	resp.Schema = schema.Schema{

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				Computed: true,
			},

			"name": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *KeystoreResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {

	var data KeystoreResourceModel

	resp.Diagnostics.Append(
		req.Plan.Get(ctx, &data)...,
	)

	if resp.Diagnostics.HasError() {
		return
	}

	keystore, err := r.client.CreateKeystore(
		data.Name.ValueString(),
	)

	if err != nil {

		resp.Diagnostics.AddError(
			"Unable to create keystore",
			err.Error(),
		)

		return
	}

	data.ID = types.StringValue(
		keystore.ID,
	)

	resp.Diagnostics.Append(
		resp.State.Set(ctx, &data)...,
	)
}

func (r *KeystoreResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {

	var data KeystoreResourceModel

	resp.Diagnostics.Append(
		req.State.Get(ctx, &data)...,
	)

	if resp.Diagnostics.HasError() {
		return
	}

	keystore, err := r.client.GetKeystore(
		data.ID.ValueString(),
	)

	if errors.Is(err, ErrNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {

		resp.Diagnostics.AddError(
			"Unable to read keystore",
			err.Error(),
		)

		return
	}

	data.Name = types.StringValue(
		keystore.Name,
	)

	resp.Diagnostics.Append(
		resp.State.Set(ctx, &data)...,
	)
}

func (r *KeystoreResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Day-3: Call PUT API
}

func (r *KeystoreResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Day-3: Call DELETE API
}
