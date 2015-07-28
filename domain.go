package generictypes

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/giantswarm/validate"
	"github.com/giantswarm/validate/web"
	"github.com/juju/errgo"
)

var (
	updatedTLDs    bool
	updateTLDMutex sync.Mutex
)

// updateTLDsIfNeeded performs a webrequest to update the
// list of toplevel domain names.
// This is done only once per process.
func updateTLDsIfNeeded() {
	if updatedTLDs {
		return
	}

	updateTLDMutex.Lock()
	defer updateTLDMutex.Unlock()

	if !updatedTLDs {
		// Fetch a new list of TLDs from the internet on startup.
		if err := web.UpdateTLDs(web.IANA); err != nil {
			log.Printf("[ERROR] Failed to update TLDs: %v\n", err)
		}
		updatedTLDs = true
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
	updateTLDsIfNeeded()
	v := validate.NewValidator()

	if err := v.Validate(web.NewDomain(d.String())); err != nil {
		return errgo.Mask(errgo.Newf("Invalid domain: %s", d.String()), errgo.Any)
	}

	return nil
}
