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
	maskAny        = errgo.MaskFunc(errgo.Any)
	updatedTLDs    bool
	updateTLDMutex sync.Mutex
)

// updateTLDsIfNeeded performs a webrequest to update the
// list of toplevel domain names.
// This is done only once per process.
func updateTLDsIfNeeded() error {
	if updatedTLDs {
		return nil
	}

	updateTLDMutex.Lock()
	defer updateTLDMutex.Unlock()

	if !updatedTLDs {
		// Fetch a new list of TLDs from the internet on startup.
		if err := web.UpdateTLDs(web.IANA); err != nil {
			return maskAny(err)
		}
		updatedTLDs = true
	}
	return nil
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
	if err := updateTLDsIfNeeded(); err != nil {
		// We don't fail Validate here, because we still have a backup
		// with our builtin TLD list.
		log.Printf("[ERROR] Failed to update TLDs: %v\n", err)
	}
	v := validate.NewValidator()

	if err := v.Validate(web.NewDomain(d.String())); err != nil {
		return maskAny(errgo.Newf("Invalid domain: %s", d.String()))
	}

	return nil
}
