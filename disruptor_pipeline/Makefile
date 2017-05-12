%.fa.gz:
	wget ftp://ftp.ensembl.org/pub/release-88/fasta/homo_sapiens/dna/$@

Homo%: Homo%.gz
	zcat $< > $@

chry.fa: Homo_sapiens.GRCh38.dna.chromosome.Y.fa
	ln -s $< $@

basecompl_chan: basecompl_chan.go
	go build $<

basecompl_disruptor: basecompl_disruptor.go
	go build $<

%.time.log: %.fa basecompl_chan basecompl_disruptor
	for i in $$(seq 1 4); do \
	echo "--------------------------------------------------------------------------------"; \
	export GOMAXPROCS=$$i; \
	echo "Setting GOMAXPROCS to $$i"; \
	echo "--------------------------------------------------------------------------------"; \
	sleep 1; \
	echo "Timing chan implementation ..."; \
	( /usr/bin/time --format 'Wall time: %e sec' ./basecompl_chan --in $< > output_chan.fa; ) 2>&1; \
	echo; \
	sleep 1; \
	echo "Timing disruptor implementation ..."; \
	( /usr/bin/time --format 'Wall time: %e sec' ./basecompl_disruptor --in $< > output_disruptor.fa; ) 2>&1; \
	echo; \
	done | tee $@;

times: chry.time.log
