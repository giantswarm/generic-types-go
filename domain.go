package generictypes

import (
	"encoding/json"
	"log"

	"github.com/giantswarm/validate"
	"github.com/giantswarm/validate/web"
	"github.com/juju/errgo"
)

func init() {
	// Fetch a new list of TLDs from the internet on startup.
	if err := web.UpdateTLDs(web.IANA); err != nil {
		log.Printf("[ERROR] Failed to update TLDs: %v\n", err)
	}
}

type Domain string

func (d *Domain) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Domain) UnmarshalJSON(data []byte) error {
	var input string

	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}

	*d = Domain(input)

	if err := d.Validate(); err != nil {
		return err
	}

	return nil
}

func (d *Domain) String() string {
	return string(*d)
}

func (d *Domain) Validate() error {
	v := validate.NewValidator()

	if err := v.Validate(web.NewDomain(d.String())); err != nil {
		return errgo.Mask(errgo.Newf("Invalid domain: %s", d.String()), errgo.Any)
	}

	return nil
}
