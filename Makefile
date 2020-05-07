all: test.svd.go rpi3_qemu.svd.go rpi3.svd.go

# on macs gsed is the name of the gnu sed installed by brew, you might need it in
# somecases because gnu sed has more features than the berkeley sed shipped with OSX
SED=sed

## this assumes that go install puts things somewhere that your PATH will find it
svd2go: structure.go svd.go template.go unmarshal_help.go useropts.go cmd/svd2go/*.go
	rm -f rpi3.svd.go rpi3_qemu.svd.go svd_out.go test.svd.go
	go install ./cmd/svd2go

rpi3_qemu.svd.go: svd2go rpi3_qemu.svd
	svd2go -o ./svd_out.go -p machine -b 'rpi3_qemu' rpi3_qemu.svd
	cat ./svd_out.go |  $(SED) '/^[[:blank:]]*$$/d' | $(SED) 's/^xxxblankxxx//g' | gofmt > rpi3_qemu.svd.go
	rm ./svd_out.go

rpi3.svd.go: svd2go rpi3.svd
	svd2go -o ./svd_out.go -p machine -b 'rpi3' rpi3.svd
	cat ./svd_out.go |  $(SED) '/^[[:blank:]]*$$/d' | $(SED) 's/^xxxblankxxx//g' | gofmt > rpi3.svd.go
	rm ./svd_out.go

test.svd.go: svd2go test.svd
	svd2go -o svd_out.go -i github.com/iansmith/svd/runtime/volatile -p svd test.svd
	cat svd_out.go |  $(SED) '/^[[:blank:]]*$$/d' | $(SED) 's/^xxxblankxxx//g' | gofmt > test.svd.go
	rm svd_out.go
