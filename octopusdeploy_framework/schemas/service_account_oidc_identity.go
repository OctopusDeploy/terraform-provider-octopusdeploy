package schemas

import (
	"regexp"

	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceAccountOIDCIdentityResourceName = "service_account_oidc_identity"
const ServiceAccountOIDCIdentityDatasourceName = "service_account_oidc_identity"

// These limits match the validation in Octopus Server.
const (
	claimConditionMaxValues      = 20
	claimConditionMaxValueLength = 256
)

// Octopus already checks these claims when it validates the token, so it rejects them in a claim condition.
var claimConditionReservedClaims = []string{"iss", "sub", "aud", "exp", "nbf", "iat", "jti"}

const claimConditionsDescription = "Extra conditions the token's claims must pass before Octopus exchanges the token for an access token. " +
	"The token must pass every condition. " +
	"When you leave this out, Terraform keeps the conditions that are already on the identity. " +
	"Set it to `[]` to remove every condition."

const claimConditionValuesDescription = "The values the claim can match. The claim passes when it matches any one value. " +
	"Values use the same matching rules as the subject: `*` matches any sequence of characters and `?` matches one character. " +
	"Matching is case sensitive. A value made only of `*` characters is rejected."

type ServiceAccountOIDCIdentitySchema struct{}

var _ EntitySchema = ServiceAccountOIDCIdentitySchema{}

func (d ServiceAccountOIDCIdentitySchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"id":   GetIdResourceSchema(),
			"name": GetNameResourceSchema(true),
			"service_account_id": util.ResourceString().
				Description("ID of the user to associate this identity to").
				Required().
				PlanModifiers(stringplanmodifier.RequiresReplace()).
				Build(),
			"issuer": util.ResourceString().
				Description("OIDC issuer url").
				Required().
				Build(),
			"subject": util.ResourceString().
				Description("OIDC subject claims").
				Required().
				Build(),
			"audience": util.ResourceString().
				Description("OIDC audience").
				Optional().
				Build(),
			"claim_conditions": resourceSchema.ListNestedAttribute{
				Description: claimConditionsDescription,
				Optional:    true,
				// Computed keeps conditions set outside Terraform, such as in the Octopus UI, when the config leaves this out.
				Computed:      true,
				PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
				NestedObject: resourceSchema.NestedAttributeObject{
					Attributes: map[string]resourceSchema.Attribute{
						"claim": resourceSchema.StringAttribute{
							Description: "The name of the token claim to check.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								stringvalidator.NoneOf(claimConditionReservedClaims...),
							},
						},
						"values": resourceSchema.ListAttribute{
							Description: claimConditionValuesDescription,
							Required:    true,
							ElementType: types.StringType,
							Validators: []validator.List{
								listvalidator.SizeBetween(1, claimConditionMaxValues),
								listvalidator.ValueStringsAre(
									stringvalidator.LengthBetween(1, claimConditionMaxValueLength),
									stringvalidator.RegexMatches(regexp.MustCompile(`[^*]`), "can't be made only of * characters"),
								),
							},
						},
					},
				},
			},
		},
		Description: "This resource manages manages OIDC service account for the associated user",
	}
}

func (d ServiceAccountOIDCIdentitySchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Attributes: map[string]datasourceSchema.Attribute{
			"id": GetIdDatasourceSchema(false),
			"service_account_id": util.DataSourceString().
				Description("ID of the user associated to this identity").
				Required().
				Build(),
			// Response
			"name": util.DataSourceString().
				Description("Name of the user associated to this identity").
				Computed().
				Build(),
			"issuer": util.DataSourceString().
				Description("OIDC issuer url").
				Computed().
				Build(),
			"subject": util.DataSourceString().
				Description("OIDC subject claims").
				Computed().
				Build(),
			"audience": util.DataSourceString().
				Description("OIDC audience").
				Computed().
				Build(),
			"claim_conditions": datasourceSchema.ListNestedAttribute{
				Description: "Extra conditions the token's claims must pass before Octopus exchanges the token for an access token.",
				Computed:    true,
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: map[string]datasourceSchema.Attribute{
						"claim": datasourceSchema.StringAttribute{
							Description: "The name of the token claim to check.",
							Computed:    true,
						},
						"values": datasourceSchema.ListAttribute{
							Description: "The values the claim can match.",
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

type OIDCServiceAccountSchemaModel struct {
	ServiceAccountID types.String `tfsdk:"service_account_id"`
	Name             types.String `tfsdk:"name"`
	Issuer           types.String `tfsdk:"issuer"`
	Subject          types.String `tfsdk:"subject"`
	Audience         types.String `tfsdk:"audience"`
	ClaimConditions  types.List   `tfsdk:"claim_conditions"`

	ResourceModel
}

type OIDCClaimConditionModel struct {
	Claim  types.String `tfsdk:"claim"`
	Values types.List   `tfsdk:"values"`
}

func OIDCClaimConditionObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"claim":  types.StringType,
			"values": types.ListType{ElemType: types.StringType},
		},
	}
}

type OIDCServiceAccountDatasourceSchemaModel struct {
	ID               types.String `tfsdk:"id"`
	ServiceAccountID types.String `tfsdk:"service_account_id"`
	Name             types.String `tfsdk:"name"`
	Issuer           types.String `tfsdk:"issuer"`
	Subject          types.String `tfsdk:"subject"`
	Audience         types.String `tfsdk:"audience"`
	ClaimConditions  types.List   `tfsdk:"claim_conditions"`
}
