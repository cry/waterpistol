## Terraform

Install terraform and run ./c2_up to bring up a c2 instance.

It'll generate a new SSH privatekey in `${HOME}/.ssh/id_c2`.

SSH in with `ssh -i ${HOME}/.ssh/id_c2 ec2-user@...`
