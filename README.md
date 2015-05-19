# generic-types-go

This repository is intended to house generally usable types that are missing
from the standard library (Domains) or a bit more specific for the
containerized world we live in (DockerImage, DockerPort).  All types should
support JSON serialization and a validation logic.  Errors are wrapped with
github.com/juju/errgo.
