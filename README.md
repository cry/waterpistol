## Terraform

Install terraform and run ./c2_up to bring up a c2 instance.

It'll generate a new SSH privatekey in `${HOME}/.ssh/id_c2`.

SSH in with `ssh -i ${HOME}/.ssh/id_c2 ec2-user@...`

Put your secret keys and ID in `~/.aws/credentials`

```
[default]
aws_access_key_id=
aws_secret_access_key=
```


## Protogen
If updating protocol, run this command
protoc -I . --go_out=plugins=grpc:. common/messages/messages.proto
