all: svd2go rpi3_qemu.svd.go rpi3.svd.go

## this assumes that go install puts things somewhere that your PATH will find it
svd2go: *.go cmd/svd2go/*.go
	rm -f rpi3.svd.go rpi3_qemu.svd.go svd_out.go
	go install ./cmd/svd2go

rpi3_qemu.svd.go: rpi3_qemu.svd svd2go
	svd2go -o ./svd_out.go -p machine -b 'rpi3_qemu' rpi3_qemu.svd
	cat ./svd_out.go |  gsed '/^[[:blank:]]*$$/d' | gofmt > rpi3_qemu.svd.go
	rm ./svd_out.go

rpi3.svd.go: rpi3.svd svd2go
	svd2go -o ./svd_out.go -p machine -b 'rpi3' rpi3.svd
	cat ./svd_out.go |  gsed '/^[[:blank:]]*$$/d' | gofmt > rpi3.svd.go
	rm ./svd_out.go
