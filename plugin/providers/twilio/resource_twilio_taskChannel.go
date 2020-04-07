package twilio

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	log "github.com/sirupsen/logrus"
)

func resourceTwilioTaskChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceTwilioTaskChannelCreate,
		Read:   resourceTwilioTaskChannelRead,
		Update: resourceTwilioTaskChannelUpdate,
		Delete: resourceTwilioTaskChannelDelete,
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
				Required: true,
			},
			"workspace_sid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"unique_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"channel_optimized_routing": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func flattenTaskChannelForCreate(d *schema.ResourceData) url.Values {
	v := make(url.Values)
	v.Add("FriendlyName", d.Get("friendly_name").(string))
	v.Add("UniqueName", d.Get("unique_name").(string))
	if val := d.Get("channel_optimized_routing").(bool); val {
		v.Add("ChannelOptimizedRouting", "true")
	}
	return v
}

func flattenTaskChannelForUpdate(d *schema.ResourceData) url.Values {
	v := make(url.Values)
	v.Add("Sid", d.Id())
	v.Add("FriendlyName", d.Get("friendly_name").(string))
	v.Add("WorkspaceSid", d.Get("workspace_sid").(string))
	return v
}

func resourceTwilioTaskChannelCreate(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioTaskChannelCreate")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	workspaceSid := d.Get("workspace_sid").(string)
	createParams := flattenTaskChannelForCreate(d)

	log.WithFields(
		log.Fields{
			"account_sid":   config.AccountSID,
			"workspace_sid": workspaceSid,
		},
	).Debug("START client.TaskRouter.Workspace.TaskChannels.Create")

	taskChannel, err := client.TaskRouter.Workspace(workspaceSid).TaskChannels.Create(context, createParams)
	if err != nil {
		log.WithFields(
			log.Fields{
				"account_sid": config.AccountSID,
			},
		).WithError(err).Error("client.TaskRouter.Workspace.TaskChannels.Create failed")

		return err
	}
	d.SetId(taskChannel.Sid)
	d.Set("friendly_name", taskChannel.FriendlyName)
	d.Set("date_created", taskChannel.DateCreated)
	d.Set("date_updated", taskChannel.DateUpdated)
	return nil
}

func resourceTwilioTaskChannelRead(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioTaskChannelRead")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	sid := d.Id()
	workspaceSid := d.Get("workspace_sid").(string)

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
		},
	).Debug("START client.TaskRouter.Workspace.TaskChannels.Get")

	taskChannel, err := client.TaskRouter.Workspace(workspaceSid).TaskChannels.Get(context, sid)
	if err != nil {
		log.WithFields(
			log.Fields{
				"parent_account_sid": config.AccountSID,
				"workspace_sid":      workspaceSid,
				"taskChannel_sid":    sid,
			},
		).WithError(err).Error("client.TaskRouter.Workspace.TaskChannels.Get failed")

		return err
	}
	d.Set("friendly_name", taskChannel.FriendlyName)
	d.Set("unique_name", taskChannel.DateUpdated)
	d.Set("date_created", taskChannel.DateCreated)
	d.Set("date_updated", taskChannel.DateUpdated)
	d.Set("channel_optimized_routing", taskChannel.ChannelOptimizedRouting)
	return nil
}

func resourceTwilioTaskChannelUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioTaskChannelUpdate")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	sid := d.Id()
	workspaceSid := d.Get("workspace_sid").(string)
	updateParams := flattenTaskChannelForUpdate(d)

	log.WithFields(
		log.Fields{
			"account_sid":     config.AccountSID,
			"workspace_sid":   workspaceSid,
			"taskChannel_sid": sid,
		},
	).Debug("client.TaskRouter.Workspace.TaskChannels.Update")

	taskChannel, err := client.TaskRouter.Workspace(workspaceSid).TaskChannels.Update(context, sid, updateParams)

	if err != nil {
		log.WithFields(
			log.Fields{
				"account_sid":     config.AccountSID,
				"workspace_sid":   workspaceSid,
				"taskChannel_sid": sid,
				"update_params":   updateParams,
			},
		).WithError(err).Error("client.TaskRouter.Workspace.TaskChannels.Update failed")

		return err
	}

	d.SetId(taskChannel.Sid)
	d.Set("sid", taskChannel.Sid)
	d.Set("friendly_name", taskChannel.FriendlyName)
	d.Set("date_created", taskChannel.DateCreated)
	d.Set("date_updated", taskChannel.DateUpdated)
	return nil
}

func resourceTwilioTaskChannelDelete(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioTaskChannelDelete")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	sid := d.Id()
	workspaceSid := d.Get("workspace_sid").(string)

	log.WithFields(
		log.Fields{
			"account_sid":     config.AccountSID,
			"workspace_sid":   workspaceSid,
			"taskChannel_sid": sid,
		},
	).Debug("START client.TaskRouter.Workspace.TaskChannels.Delete")

	err := client.TaskRouter.Workspace(workspaceSid).TaskChannels.Delete(context, sid)

	log.WithFields(
		log.Fields{
			"account_sid":     config.AccountSID,
			"workspace_sid":   workspaceSid,
			"taskChannel_sid": sid,
		},
	).Debug("END client.TaskRouter.Workspace.TaskChannels.Delete")
	if err != nil {
		return fmt.Errorf("Failed to delete taskChannel: %s", err.Error())
	}
	return nil
}
