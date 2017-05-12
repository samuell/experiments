%.fa.gz:
	wget ftp://ftp.ensembl.org/pub/release-88/fasta/homo_sapiens/dna/$@

%: %.gz
	zcat $< > $@

chr1.fa: Homo_sapiens.GRCh38.dna.chromosome.1.fa
	ln -s %< %@

all.fa: Homo_sapiens.GRCh38.dna.toplevel.fa.gz
	ln -s %< %@

basecompl_chan: basecompl_chan.go
	go build $<

basecompl_disruptor: basecompl_disruptor.go
	go build $<

times_all.log: all.fa basecompl_chan basecompl_disruptor
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

times_chr1.log: all.fa basecompl_chan basecompl_disruptor
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

times: times_chr1.log times_all.log
