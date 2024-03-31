## Use Ubuntu 18.04 with ssh key

```
curl -o root-drive-with-ssh.img https://s3.amazonaws.com/spec.ccfc.min/ci-artifacts/disks/x86_64/ubuntu-18.04.ext4

curl -o vmlinux https://s3.amazonaws.com/spec.ccfc.min/img/quickstart_guide/x86_64/kernels/vmlinux.bin

curl -o root-drive-ssh-key https://s3.amazonaws.com/spec.ccfc.min/ci-artifacts/disks/x86_64/ubuntu-18.04.id_rsa
```

## Download firecracker binary
```
curl -L https://github.com/firecracker-microvm/firecracker/releases/download/v1.7.0/firecracker-v1.7.0-x86_64.tgz | tar -xz

mv release-v1.7.0-x86_64/firecracker-v1.7.0-x86_64 firecracker

rm -rf release-v1.7.0-x86_64
```