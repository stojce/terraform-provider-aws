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
)

func ResourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

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
			"portal_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
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

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IoTSiteWiseConn

	name := d.Get("name").(string)
	input := &iotsitewise.CreateProjectInput{
		ProjectName: aws.String(name),
		PortalId:    aws.String(d.Get("portal_id").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		input.ProjectDescription = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Creating IoTSiteWise Project: %s", input)
	output, err := conn.CreateProject(input)

	if err != nil {
		return fmt.Errorf("error creating IoT Project (%s): %w", name, err)
	}

	d.SetId(aws.StringValue(output.ProjectId))

	return resourceProjectRead(d, meta)
}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IoTSiteWiseConn

	input := &iotsitewise.DescribeProjectInput{
		ProjectId: aws.String(d.Id()),
	}

	output, err := conn.DescribeProject(input)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] IoTSiteWise Project (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading IoTSiteWise Project (%s): %w", d.Id(), err)
	}

	d.Set("arn", output.ProjectArn)
	d.Set("name", output.ProjectName)
	d.Set("portal_id", output.PortalId)
	d.Set("description", output.ProjectDescription)
	d.Set("creation_date", output.ProjectCreationDate.Format(time.RFC3339))
	d.Set("last_update_date", output.ProjectLastUpdateDate.Format(time.RFC3339))

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IoTSiteWiseConn

	input := &iotsitewise.UpdateProjectInput{
		ProjectId:   aws.String(d.Id()),
		ProjectName: aws.String(d.Get("name").(string)),
	}

	if d.HasChange("description") {
		input.ProjectDescription = aws.String(d.Get("description").(string))
	}

	log.Printf("[DEBUG] Updating IoTSiteWise Project: %s", input)
	_, err := conn.UpdateProject(input)

	if err != nil {
		return fmt.Errorf("error updating IoTSiteWise Project (%s): %w", d.Id(), err)
	}

	return resourceProjectRead(d, meta)
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IoTSiteWiseConn

	log.Printf("[DEBUG] Deleting IoTSiteWise Project: %s", d.Id())
	_, err := conn.DeleteProject(&iotsitewise.DeleteProjectInput{
		ProjectId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, iotsitewise.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting IoT Project (%s): %w", d.Id(), err)
	}

	return nil
}
