Homo_sapiens.GRCh37.67.dna_rm.chromosome.Y.fa.gz:
	wget ftp://ftp.ensembl.org/pub/release-67/fasta/homo_sapiens/dna/$@

Homo_sapiens.GRCh37.67.dna_rm.chromosome.Y.fa: Homo_sapiens.GRCh37.67.dna_rm.chromosome.Y.fa.gz
	zcat $< > $@

get_data: Homo_sapiens.GRCh37.67.dna_rm.chromosome.Y.fa

basecompl_chan: basecompl_chan.go
	go build $<

basecompl_disruptor: basecompl_disruptor.go
	go build $<

times.log: get_data basecompl_chan basecompl_disruptor
	for i in $$(seq 1 4); do \
	echo "--------------------------------------------------------------------------------"; \
	export GOMAXPROCS=$$i; \
	echo "Setting GOMAXPROCS to $$i"; \
	echo "--------------------------------------------------------------------------------"; \
	sleep 1; \
	echo "Timing chan implementation ..."; \
	( /usr/bin/time --format 'Wall time: %e sec' ./basecompl_chan > output_chan.fa; ) 2>&1; \
	echo; \
	sleep 1; \
	echo "Timing disruptor implementation ..."; \
	( /usr/bin/time --format 'Wall time: %e sec' ./basecompl_disruptor > output_disruptor.fa; ) 2>&1; \
	echo; \
	done | tee $@;

times: times.log
