package twilio

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	log "github.com/sirupsen/logrus"
)

func resourceTwilioPhoneNumber() *schema.Resource {
	return &schema.Resource{
		Create: resourceTwilioPhoneNumberCreate,
		Read:   resourceTwilioPhoneNumberRead,
		Update: resourceTwilioPhoneNumberUpdate,
		Delete: resourceTwilioPhoneNumberDelete,
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
			"country_code": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"search": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			// TODO: We should also be able to handle "capabilities" but skipping it
			// because it is challenging to parse lists and pass them along to the underlying
			// go library
		},
	}
}

func flattenPhoneNumberForSearch(d *schema.ResourceData) url.Values {
	v := make(url.Values)
	search := d.Get("search").(string)
	v.Add("Contains", search)
	// TODO: The below code is commented out because we are not able
	// to pin down the area code by doing this if statement. Need to figure it out
	// if strings.Contains(search, "*") {
	// 	v.Add("Contains", search)
	// } else {
	// 	v.Add("AreaCode", search)
	// }
	return v
}

func flattenPhoneNumberForBuying(d *schema.ResourceData, phoneNumber string) url.Values {
	v := make(url.Values)
	v.Add("FriendlyName", d.Get("friendly_name").(string))
	v.Add("PhoneNumber", phoneNumber)
	return v
}

func flattenPhoneNumberForUpdate(d *schema.ResourceData) url.Values {
	v := make(url.Values)
	v.Add("FriendlyName", d.Get("friendly_name").(string))
	return v
}

func resourceTwilioPhoneNumberCreate(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioPhoneNumberCreate")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	searchParams := flattenPhoneNumberForSearch(d)
	countryCode := d.Get("country_code").(string)

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
		},
	).Debug("START client.AvailableNumbers.Local.GetPage")

	numbers, err := client.AvailableNumbers.Local.GetPage(context, countryCode, searchParams)
	if err != nil {
		log.WithFields(
			log.Fields{
				"account_sid": config.AccountSID,
			},
		).WithError(err).Error("client.AvailableNumbers.Local.GetPage failed")

		return err
	}

	if !(len(numbers.Numbers) >= 1) {
		log.WithFields(
			log.Fields{
				"account_sid": config.AccountSID,
			},
		).WithError(err).Error("client.AvailableNumbers.Local.GetPage failed to find a valid number")
		return nil
	}

	phoneNumber := numbers.Numbers[0]
	createParams := flattenPhoneNumberForBuying(d, string(phoneNumber.PhoneNumber))

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
		},
	).Debug("START client.IncomingNumbers.Create")

	boughtNumber, err := client.IncomingNumbers.Create(context, createParams)
	if err != nil {
		log.WithFields(
			log.Fields{
				"account_sid": config.AccountSID,
			},
		).WithError(err).Error("client.IncomingNumbers.Create failed")

		return err
	}
	d.SetId(boughtNumber.Sid)
	d.Set("friendly_name", boughtNumber.FriendlyName)
	d.Set("phone_number", boughtNumber.PhoneNumber)
	d.Set("date_created", boughtNumber.DateCreated)
	d.Set("date_updated", boughtNumber.DateUpdated)
	d.Set("capabilities", boughtNumber.Capabilities)
	return nil
}

func resourceTwilioPhoneNumberRead(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioPhoneNumberRead")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	sid := d.Id()

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
		},
	).Debug("START client.IncomingNumbers.Get")

	phoneNumber, err := client.IncomingNumbers.Get(context, sid)
	if err != nil {
		log.WithFields(
			log.Fields{
				"account_sid": config.AccountSID,
			},
		).WithError(err).Error("client.IncomingNumbers.Get failed")

		return err
	}
	d.Set("friendly_name", phoneNumber.FriendlyName)
	d.Set("phone_number", phoneNumber.PhoneNumber)
	d.Set("date_created", phoneNumber.DateCreated)
	d.Set("date_updated", phoneNumber.DateUpdated)
	d.Set("capabilities", phoneNumber.Capabilities)
	return nil
}

func resourceTwilioPhoneNumberUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioPhoneNumberUpdate")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	sid := d.Id()
	createParams := flattenPhoneNumberForUpdate(d)

	log.WithFields(
		log.Fields{
			"account_sid": config.AccountSID,
		},
	).Debug("START client.IncomingNumbers.Update")

	updatedNumber, err := client.IncomingNumbers.Update(context, sid, createParams)
	if err != nil {
		log.WithFields(
			log.Fields{
				"account_sid": config.AccountSID,
			},
		).WithError(err).Error("client.IncomingNumbers.Update failed")

		return err
	}
	d.SetId(updatedNumber.Sid)
	d.Set("friendly_name", updatedNumber.FriendlyName)
	d.Set("phone_number", updatedNumber.PhoneNumber)
	d.Set("date_created", updatedNumber.DateCreated)
	d.Set("date_updated", updatedNumber.DateUpdated)
	d.Set("capabilities", updatedNumber.Capabilities)
	return nil
}

func resourceTwilioPhoneNumberDelete(d *schema.ResourceData, meta interface{}) error {
	log.Debug("ENTER resourceTwilioPhoneNumberDelete")

	client := meta.(*TerraformTwilioContext).client
	config := meta.(*TerraformTwilioContext).configuration
	context := context.TODO()

	sid := d.Id()

	log.WithFields(
		log.Fields{
			"account_sid":     config.AccountSID,
			"phoneNumber_sid": sid,
		},
	).Debug("START client.IncomingNumbers.Delete")

	err := client.IncomingNumbers.Release(context, sid)

	log.WithFields(
		log.Fields{
			"account_sid":     config.AccountSID,
			"phoneNumber_sid": sid,
		},
	).Debug("END client.IncomingNumbers.Delete")
	if err != nil {
		return fmt.Errorf("Failed to delete phoneNumber: %s", err.Error())
	}
	return nil
}
