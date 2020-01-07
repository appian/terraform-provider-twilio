package twilio

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	log "github.com/sirupsen/logrus"
)

func resourceTwilioWorkflow() *schema.Resource {
	return &schema.Resource{
		Create: resourceTwilioWorkflowCreate,
		Read:   resourceTwilioWorkflowRead,
		Update: resourceTwilioWorkflowUpdate,
		Delete: resourceTwilioWorkflowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"sid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"friendly_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace_sid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"configuration": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func flattenWorkflowForCreate(d *schema.ResourceData) url.Values {
	v := make(url.Values)

	v.Add("FriendlyName", d.Get("friendly_name").(string))
	v.Add("Configuration", d.Get("configuration").(string))
	return v
}

func resourceTwilioWorkflowCreate(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioWorkflowCreate")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	workspaceSid := d.Get("workspace_sid").(string)
	createParams := flattenWorkflowForCreate(d)

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
		},
	).Debug("START client.TaskRouter.Workspace.Workflows.Create")

	workflow, err := client.TaskRouter.Workspace(workspaceSid).Workflows.Create(context, createParams)
	if err != nil {
		log.WithFields(
			log.Fields{
				"account_sid": config.AccountSID,
			},
		).WithError(err).Error("client.TaskRouter.Workspace.Workflows.Create failed")

		return err
	}
	d.SetId(workflow.Sid)
	d.Set("friendly_name", workflow.FriendlyName)
	d.Set("date_created", workflow.DateCreated)
	d.Set("date_updated", workflow.DateUpdated)
	d.Set("task_reservation_timeout", workflow.TaskReservationTimeout)
	return nil
}

func resourceTwilioWorkflowRead(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioWorkflowRead")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	sid := d.Id()
	workspaceSid := d.Get("workspace_sid").(string)

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
		},
	).Debug("START client.TaskRouter.Workspace.Workflows.Get")

	workflow, err := client.TaskRouter.Workspace(workspaceSid).Workflows.Get(context, sid)
	if err != nil {
		log.WithFields(
			log.Fields{
				"parent_account_sid": config.AccountSID,
			},
		).WithError(err).Error("client.TaskRouter.Workspace.Workflows.Get failed")

		return err
	}
	d.Set("friendly_name", workflow.FriendlyName)
	d.Set("date_created", workflow.DateCreated)
	d.Set("date_updated", workflow.DateUpdated)
	d.Set("task_reservation_timeout", workflow.TaskReservationTimeout)
	return nil
}

func resourceTwilioWorkflowUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioWorkflowUpdate")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	sid := d.Id()
	workspaceSid := d.Get("workspace_sid").(string)
	createParams := flattenWorkflowForCreate(d)

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
		},
	).Debug("START client.TaskRouter.Workspace.Workflows.Update")

	workflow, err := client.TaskRouter.Workspace(workspaceSid).Workflows.Update(context, sid, createParams)
	if err != nil {
		log.WithFields(
			log.Fields{
				"account_sid": config.AccountSID,
			},
		).WithError(err).Error("client.TaskRouter.Workspace.Workflows.Update failed")

		return err
	}
	d.SetId(workflow.Sid)
	d.Set("friendly_name", workflow.FriendlyName)
	d.Set("date_created", workflow.DateCreated)
	d.Set("date_updated", workflow.DateUpdated)
	d.Set("task_reservation_timeout", workflow.TaskReservationTimeout)
	return nil
}

func resourceTwilioWorkflowDelete(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioWorkflowDelete")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	sid := d.Id()
	workspaceSid := d.Get("workspace_sid").(string)

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
			"queue_sid":   sid,
		},
	).Debug("START client.TaskRouter.Workspace.Workflows.Delete")

	err := client.TaskRouter.Workspace(workspaceSid).Workflows.Delete(context, sid)

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
			"queue_sid":   sid,
		},
	).Debug("END client.TaskRouter.Workspace.Workflows.Delete")
	if err != nil {
		return fmt.Errorf("Failed to delete workflow: %s", err.Error())
	}
	return nil
}
