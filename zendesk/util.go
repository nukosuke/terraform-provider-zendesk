package zendesk

import "github.com/hashicorp/terraform/helper/schema"

func setSchemaFields(d *schema.ResourceData, m map[string]interface{}) error {
	for k, v := range m {
		err := d.Set(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}
