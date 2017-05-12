# Comparing channel and disruptor pattern performance

This is an experiment to see how using the [LMAX disruptor pattern](https://lmax-exchange.github.io/disruptor/)
as implemented in a slightly adapted version of [Darren Elwood](https://github.com/textnode)'s [gringo](https://github.com/textnode/gringo)
lock-free ring buffer, as an alternative to plain Go channels.

The example program used, is a chain of 8 similar [base-complementer](https://en.wikipedia.org/wiki/Complementarity_(molecular_biology)#DNA_and_RNA_base_pair_complementarity)
on a [58 Mb DNA sequence text file](ftp://ftp.ensembl.org/pub/release-67/fasta/homo_sapiens/dna/Homo_sapiens.GRCh37.67.dna_rm.chromosome.Y.fa.gz),
in [FASTA format](https://en.wikipedia.org/wiki/FASTA_format).

I'm sure I have made some stupid errors leading to sub-par performance in these codes. Please contact
me on [@smllmp](http://twitter.com/smllmp) or in the [issue tracker](https://github.com/samuell/experiments/issues)
to suggest improvements!

## Prerequisites

- Linux
- The Go programming language
- The Afero file system abstraction layer<br>
  (Install with `go get github.com/spf13/afero`)

## How to run

```
make times
```

## Example output

This is output I get when I run on my Laptop, with an Intel Core i5-4210U CPU (1.7GHz base, 2.7GHz max):

```
--------------------------------------------------------------------------------
Setting GOMAXPROCS to 1
--------------------------------------------------------------------------------
Timing chan implementation ...
Wall time: 4.05 sec

Timing disruptor implementation ...
Wall time: 3.25 sec

--------------------------------------------------------------------------------
Setting GOMAXPROCS to 2
--------------------------------------------------------------------------------
Timing chan implementation ...
Wall time: 2.78 sec

Timing disruptor implementation ...
Wall time: 2.65 sec

--------------------------------------------------------------------------------
Setting GOMAXPROCS to 3
--------------------------------------------------------------------------------
Timing chan implementation ...
Wall time: 2.82 sec

Timing disruptor implementation ...
Wall time: 2.77 sec

--------------------------------------------------------------------------------
Setting GOMAXPROCS to 4
--------------------------------------------------------------------------------
Timing chan implementation ...
Wall time: 3.05 sec

Timing disruptor implementation ...
Wall time: 3.39 sec
```
