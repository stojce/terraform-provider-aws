package iotsitewise

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourcePortal() *schema.Resource {
	return &schema.Resource{
		Create: resourcePortalCreate,
		Read:   resourcePortalRead,
		Update: resourcePortalUpdate,
		Delete: resourcePortalDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"portal_contact_email": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"role_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidARN,
			},

			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auth_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_update_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePortalCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IoTSiteWiseConn

	name := d.Get("name").(string)
	input := &iotsitewise.CreatePortalInput{
		PortalName: aws.String(name),
	}

	if v, ok := d.GetOk("portal_contact_email"); ok {
		input.PortalContactEmail = aws.String(v.(string))
	}

	if v, ok := d.GetOk("role_arn"); ok {
		input.RoleArn = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Creating IoTSiteWise Portal: %s", input)
	output, err := conn.CreatePortal(input)

	if err != nil {
		return fmt.Errorf("error creating IoT Portal (%s): %w", name, err)
	}

	d.SetId(aws.StringValue(output.PortalId))

	err = conn.WaitUntilPortalActive(&iotsitewise.DescribePortalInput{
		PortalId: output.PortalId,
	})

	if err != nil {
		return fmt.Errorf("error waiting for Portal (%s) to become ACTIVE: %w", aws.StringValue(output.PortalId), err)
	}

	return resourcePortalRead(d, meta)
}

func resourcePortalRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IoTSiteWiseConn

	input := &iotsitewise.DescribePortalInput{
		PortalId: aws.String(d.Id()),
	}

	output, err := conn.DescribePortal(input)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] IoTSiteWise Portal (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading IoTSiteWise Portal (%s): %w", d.Id(), err)
	}

	d.Set("arn", output.PortalArn)
	d.Set("role_arn", output.RoleArn)
	d.Set("name", output.PortalName)
	d.Set("portal_contact_email", output.PortalContactEmail)
	d.Set("auth_mode", output.PortalAuthMode)
	d.Set("description", output.PortalDescription)

	d.Set("client_id", output.PortalClientId)
	d.Set("start_url", output.PortalStartUrl)
	d.Set("creation_date", output.PortalCreationDate.Format(time.RFC3339))
	d.Set("last_update_date", output.PortalLastUpdateDate.Format(time.RFC3339))

	// Alarms *Alarms `locationName:"alarms" type:"structure"`
	// NotificationSenderEmail *string `locationName:"notificationSenderEmail" min:"1" type:"string"`
	// PortalLogoImageLocation *ImageLocation `locationName:"portalLogoImageLocation" type:"structure"`
	// PortalStatus *PortalStatus `locationName:"portalStatus" type:"structure" required:"true"`

	return nil
}

func resourcePortalUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IoTSiteWiseConn

	input := &iotsitewise.UpdatePortalInput{
		PortalId:           aws.String(d.Id()),
		PortalName:         aws.String(d.Get("name").(string)),
		PortalContactEmail: aws.String(d.Get("portal_contact_email").(string)),
		RoleArn:            aws.String(d.Get("role_arn").(string)),
	}

	log.Printf("[DEBUG] Updating IoTSiteWise Portal: %s", input)
	_, err := conn.UpdatePortal(input)

	if err != nil {
		return fmt.Errorf("error updating IoTSiteWise Portal (%s): %w", d.Id(), err)
	}

	return resourcePortalRead(d, meta)
}

func resourcePortalDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IoTSiteWiseConn

	log.Printf("[DEBUG] Deleting IoTSiteWise Portal: %s", d.Id())
	_, err := conn.DeletePortal(&iotsitewise.DeletePortalInput{
		PortalId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, iotsitewise.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting IoT Portal (%s): %w", d.Id(), err)
	}

	return nil
}
